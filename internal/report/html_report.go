package report

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// HTMLReporter generates a fully self-contained HTML report with embedded data.
type HTMLReporter struct{}

func NewHTMLReporter() *HTMLReporter { return &HTMLReporter{} }

// Generate produces the full HTML report string.
func (h *HTMLReporter) Generate(findings model.FindingsBundle, insights *model.InsightsBundle) (string, error) {
	findingsJSON, err := json.Marshal(findings)
	if err != nil {
		return "", fmt.Errorf("marshal findings: %w", err)
	}

	insightsJSON := []byte("null")
	if insights != nil {
		insightsJSON, err = json.Marshal(insights)
		if err != nil {
			return "", fmt.Errorf("marshal insights: %w", err)
		}
	}

	generatedAt := time.Now().Format("January 2, 2006 at 15:04 MST")

	// Build per-finding insight lookup map for Go-side rendering (used in template)
	insightMap := map[string]model.FindingInsight{}
	if insights != nil {
		for _, fi := range insights.PerFinding {
			insightMap[fi.FindingID] = fi
		}
	}

	data := templateData{
		GeneratedAt:  generatedAt,
		FindingsJSON: template.JS(findingsJSON),
		InsightsJSON: template.JS(insightsJSON),
		Findings:     findings,
		Insights:     insights,
		InsightMap:   insightMap,
	}

	var sb strings.Builder
	if err := reportTmpl.Execute(&sb, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return sb.String(), nil
}

type templateData struct {
	GeneratedAt  string
	FindingsJSON template.JS
	InsightsJSON template.JS
	Findings     model.FindingsBundle
	Insights     *model.InsightsBundle
	InsightMap   map[string]model.FindingInsight
}

// severityColor returns a CSS hex color for a severity string.
func severityColor(s string) string {
	switch strings.ToLower(s) {
	case "critical":
		return "#e53e3e"
	case "high":
		return "#dd6b20"
	case "medium":
		return "#d69e2e"
	default:
		return "#38a169"
	}
}

// severityEmoji returns an emoji for a severity string.
func severityEmoji(s string) string {
	switch strings.ToLower(s) {
	case "critical":
		return "🔴"
	case "high":
		return "🟠"
	case "medium":
		return "🟡"
	default:
		return "🟢"
	}
}

// profileEmoji returns a category emoji.
func profileEmoji(cat string) string {
	switch strings.ToLower(cat) {
	case "cpu":
		return "⚡"
	case "heap", "allocs":
		return "🧠"
	case "goroutine":
		return "🧵"
	case "mutex":
		return "🔒"
	case "block":
		return "⏸️"
	default:
		return "📊"
	}
}

// truncate truncates a string to n runes.
func truncate(n int, s string) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}

var funcMap = template.FuncMap{
	"severityColor": severityColor,
	"severityEmoji": severityEmoji,
	"profileEmoji":  profileEmoji,
	"truncate":      truncate,
	"lower":         strings.ToLower,
	"title": func(s string) string {
		if len(s) == 0 {
			return s
		}
		return strings.ToUpper(s[:1]) + s[1:]
	},
	"add": func(a, b int) int { return a + b },
	"hasInsight": func(m map[string]model.FindingInsight, id string) bool {
		_, ok := m[id]
		return ok
	},
	"getInsight": func(m map[string]model.FindingInsight, id string) model.FindingInsight {
		return m[id]
	},
	"insightsEnabled": func(ins *model.InsightsBundle) bool {
		return ins != nil && ins.DisabledReason == "" && ins.ExecutiveSummary.Overview != ""
	},
}

var reportTmpl = template.Must(template.New("report").Funcs(funcMap).Parse(rawHTMLTemplate))

const rawHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>TriageProf — Performance Report</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.3/dist/chart.umd.min.js"></script>
<style>
/* ===== RESET & TOKENS ===== */
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0f1117;--surface:#1a1d27;--surface2:#22263a;--border:#2d3148;
  --text:#e2e8f0;--muted:#8892a4;--accent:#7c6af7;--accent2:#56cfb2;
  --red:#fc8181;--orange:#f6ad55;--yellow:#f6e05e;--green:#68d391;
  --red-bg:rgba(252,129,129,.12);--orange-bg:rgba(246,173,85,.12);
  --yellow-bg:rgba(246,224,94,.12);--green-bg:rgba(104,211,145,.12);
  --radius:12px;--radius-sm:8px;--shadow:0 4px 24px rgba(0,0,0,.4);
}
body{font-family:'Inter',system-ui,-apple-system,sans-serif;background:var(--bg);
  color:var(--text);line-height:1.6;min-height:100vh}
a{color:var(--accent);text-decoration:none}

/* ===== LAYOUT ===== */
.wrap{max-width:1100px;margin:0 auto;padding:32px 20px}
.grid2{display:grid;grid-template-columns:1fr 1fr;gap:20px}
.grid3{display:grid;grid-template-columns:repeat(3,1fr);gap:16px}
.grid4{display:grid;grid-template-columns:repeat(4,1fr);gap:16px}
@media(max-width:700px){.grid2,.grid3,.grid4{grid-template-columns:1fr}}

/* ===== CARDS ===== */
.card{background:var(--surface);border:1px solid var(--border);border-radius:var(--radius);padding:24px;box-shadow:var(--shadow)}
.card+.card,.card+.finding{margin-top:20px}
.card h2{font-size:1.1rem;font-weight:600;margin-bottom:16px;color:var(--text);
  display:flex;align-items:center;gap:8px}
.card h2 .icon{opacity:.7}

/* ===== HEADER ===== */
.header{background:linear-gradient(135deg,#1e1b4b 0%,#312e81 50%,#1a1d27 100%);
  border-bottom:1px solid var(--border);padding:40px 20px 32px;text-align:center;margin-bottom:32px}
.header h1{font-size:2rem;font-weight:700;letter-spacing:-0.03em;
  background:linear-gradient(90deg,#a78bfa,#56cfb2);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.header .sub{color:var(--muted);margin-top:6px;font-size:.95rem}
.header .meta{display:flex;justify-content:center;gap:24px;margin-top:20px;flex-wrap:wrap}
.header .meta span{font-size:.85rem;color:var(--muted);display:flex;align-items:center;gap:6px}
.header .meta strong{color:var(--text)}

/* ===== SCORE GAUGE ===== */
.gauge-wrap{display:flex;flex-direction:column;align-items:center;gap:8px}
.gauge-ring{position:relative;width:120px;height:120px}
.gauge-ring svg{transform:rotate(-90deg)}
.gauge-num{position:absolute;inset:0;display:flex;align-items:center;justify-content:center;
  font-size:1.8rem;font-weight:700;color:var(--text)}
.gauge-label{font-size:.8rem;color:var(--muted);text-align:center}

/* ===== SEVERITY BADGE ===== */
.badge{display:inline-flex;align-items:center;gap:4px;padding:3px 10px;border-radius:20px;
  font-size:.75rem;font-weight:600;letter-spacing:.02em;text-transform:uppercase}
.badge.critical{background:var(--red-bg);color:var(--red);border:1px solid rgba(252,129,129,.3)}
.badge.high{background:var(--orange-bg);color:var(--orange);border:1px solid rgba(246,173,85,.3)}
.badge.medium{background:var(--yellow-bg);color:var(--yellow);border:1px solid rgba(246,224,94,.3)}
.badge.low,.badge.info{background:var(--green-bg);color:var(--green);border:1px solid rgba(104,211,145,.3)}

/* ===== STAT TILES ===== */
.stat{background:var(--surface2);border:1px solid var(--border);border-radius:var(--radius-sm);
  padding:18px;text-align:center}
.stat .val{font-size:2rem;font-weight:700;line-height:1}
.stat .lbl{font-size:.8rem;color:var(--muted);margin-top:4px}
.stat.critical .val{color:var(--red)}
.stat.high .val{color:var(--orange)}
.stat.medium .val{color:var(--yellow)}
.stat.low .val{color:var(--green)}
.stat.accent .val{color:var(--accent)}

/* ===== FINDINGS ===== */
.finding{background:var(--surface);border:1px solid var(--border);border-radius:var(--radius);
  overflow:hidden;box-shadow:var(--shadow);margin-top:20px;transition:border-color .2s}
.finding:hover{border-color:var(--accent)}
.finding-hdr{padding:18px 24px;display:flex;align-items:center;gap:14px;
  background:var(--surface2);border-bottom:1px solid var(--border);cursor:pointer;user-select:none}
.finding-hdr .emoji{font-size:1.4rem}
.finding-hdr .info{flex:1}
.finding-hdr .title{font-weight:600;font-size:1rem}
.finding-hdr .sub{font-size:.82rem;color:var(--muted);margin-top:2px}
.finding-hdr .right{display:flex;align-items:center;gap:10px}
.finding-hdr .score{font-size:.85rem;color:var(--muted);background:var(--border);
  padding:2px 10px;border-radius:20px}
.finding-body{padding:24px;display:none}
.finding-body.open{display:block}

/* ===== HOTSPOT BARS ===== */
.hotspot-list{display:flex;flex-direction:column;gap:8px;margin-top:12px}
.hotspot{display:flex;flex-direction:column;gap:4px}
.hotspot-meta{display:flex;justify-content:space-between;font-size:.8rem}
.hotspot-fn{color:var(--text);font-family:monospace;max-width:75%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.hotspot-val{color:var(--muted)}
.hotspot-bar{height:6px;border-radius:3px;background:var(--border)}
.hotspot-fill{height:100%;border-radius:3px;background:linear-gradient(90deg,var(--accent),var(--accent2));
  transition:width .6s cubic-bezier(.4,0,.2,1)}

/* ===== AI INSIGHTS ===== */
.ai-card{background:linear-gradient(135deg,rgba(124,106,247,.15),rgba(86,207,178,.1));
  border:1px solid rgba(124,106,247,.4);border-radius:var(--radius);padding:24px;margin-top:12px}
.ai-card .ai-header{display:flex;align-items:center;gap:10px;margin-bottom:16px}
.ai-card .ai-header h3{font-size:.95rem;font-weight:600;color:var(--accent)}
.ai-badge{display:inline-flex;align-items:center;gap:4px;padding:2px 10px;border-radius:12px;
  font-size:.75rem;background:rgba(124,106,247,.2);color:var(--accent);border:1px solid rgba(124,106,247,.3)}
.ai-section{margin-top:14px}
.ai-section h4{font-size:.82rem;text-transform:uppercase;letter-spacing:.06em;color:var(--muted);margin-bottom:8px}
.ai-narrative{font-size:.9rem;line-height:1.7;color:var(--text);background:rgba(0,0,0,.2);
  padding:14px;border-radius:var(--radius-sm);border-left:3px solid var(--accent)}
.ai-list{list-style:none;display:flex;flex-direction:column;gap:6px}
.ai-list li{font-size:.88rem;color:var(--text);padding:8px 12px;background:rgba(0,0,0,.2);
  border-radius:var(--radius-sm);display:flex;align-items:flex-start;gap:8px}
.ai-list li::before{content:"→";color:var(--accent2);flex-shrink:0;margin-top:1px}
.code-block{background:#0d1117;border:1px solid var(--border);border-radius:var(--radius-sm);
  padding:12px 16px;font-family:'Fira Code','Cascadia Code',monospace;font-size:.8rem;
  color:#a5d6ff;overflow-x:auto;margin-top:6px;white-space:pre-wrap;word-break:break-all}
.confidence-pill{display:inline-flex;align-items:center;gap:4px;padding:2px 10px;
  border-radius:12px;font-size:.75rem;font-weight:600}
.confidence-high{background:rgba(104,211,145,.15);color:var(--green);border:1px solid rgba(104,211,145,.3)}
.confidence-mid{background:rgba(246,224,94,.15);color:var(--yellow);border:1px solid rgba(246,224,94,.3)}
.confidence-low{background:rgba(252,129,129,.15);color:var(--red);border:1px solid rgba(252,129,129,.3)}

/* ===== RECOMMENDATIONS PANEL ===== */
.reco-grid{display:flex;flex-direction:column;gap:16px;margin-top:4px}
.reco{background:var(--surface2);border:1px solid var(--border);border-radius:var(--radius-sm);padding:18px;transition:border-color .2s}
.reco:hover{border-color:var(--accent)}
.reco-hdr{display:flex;align-items:flex-start;gap:12px;margin-bottom:10px}
.reco-num{width:28px;height:28px;border-radius:50%;display:flex;align-items:center;justify-content:center;
  font-size:.8rem;font-weight:700;flex-shrink:0;margin-top:2px}
.reco-num.high{background:var(--orange-bg);color:var(--orange)}
.reco-num.medium{background:var(--yellow-bg);color:var(--yellow)}
.reco-num.low{background:var(--green-bg);color:var(--green)}
.reco-num.critical{background:var(--red-bg);color:var(--red)}
.reco-title-block{flex:1}
.reco-title{font-size:.95rem;font-weight:600;line-height:1.4}
.reco-meta{display:flex;flex-wrap:wrap;gap:8px;margin-top:6px}
.reco-tag{display:inline-flex;align-items:center;gap:4px;padding:2px 9px;border-radius:10px;
  font-size:.72rem;font-weight:600;border:1px solid}
.reco-tag.effort{background:rgba(99,102,241,.12);color:#a5b4fc;border-color:rgba(99,102,241,.3)}
.reco-tag.complexity-Low{background:rgba(104,211,145,.12);color:var(--green);border-color:rgba(104,211,145,.3)}
.reco-tag.complexity-Medium{background:rgba(246,224,94,.12);color:var(--yellow);border-color:rgba(246,224,94,.3)}
.reco-tag.complexity-High{background:rgba(252,129,129,.12);color:var(--red);border-color:rgba(252,129,129,.3)}
.reco-tag.category{background:rgba(86,207,178,.1);color:var(--accent2);border-color:rgba(86,207,178,.25)}
.reco-impact-bar{display:flex;align-items:center;gap:10px;padding:8px 12px;background:rgba(104,211,145,.06);
  border:1px solid rgba(104,211,145,.2);border-radius:6px;margin-bottom:10px}
.reco-impact-label{font-size:.75rem;font-weight:700;text-transform:uppercase;letter-spacing:.05em;color:var(--green);white-space:nowrap}
.reco-impact-val{font-size:.85rem;color:var(--fg)}
.reco-validation{margin-top:8px}
.reco-validation-title{font-size:.72rem;font-weight:700;text-transform:uppercase;letter-spacing:.05em;color:var(--muted);margin-bottom:4px}
.reco-validation-list{display:flex;flex-wrap:wrap;gap:6px}
.reco-validation-item{font-size:.77rem;padding:2px 8px;background:rgba(255,255,255,.04);border:1px solid var(--border);border-radius:10px;color:var(--muted)}
.reco-code-block{background:#0d1117;border:1px solid rgba(255,255,255,.08);border-radius:8px;
  padding:12px 16px;font-family:'Fira Code','Cascadia Code',monospace;font-size:.8rem;
  color:#a5d6ff;overflow-x:auto;margin-top:10px;white-space:pre-wrap;word-break:break-all}
.reco code{font-family:monospace;font-size:.8rem;background:rgba(0,0,0,.3);padding:2px 6px;border-radius:4px;color:#a5d6ff}

/* ===== RISK LIST ===== */
.risk-list{display:flex;flex-direction:column;gap:10px;margin-top:4px}
.risk{padding:12px 16px;border-radius:var(--radius-sm);border-left:4px solid;display:flex;gap:12px;align-items:flex-start}
.risk.critical{background:var(--red-bg);border-color:var(--red)}
.risk.high{background:var(--orange-bg);border-color:var(--orange)}
.risk.medium{background:var(--yellow-bg);border-color:var(--yellow)}
.risk.low{background:var(--green-bg);border-color:var(--green)}
.risk-body .risk-title{font-weight:600;font-size:.9rem}
.risk-body .risk-detail{font-size:.82rem;color:var(--muted);margin-top:3px}

/* ===== CHART WRAPPER ===== */
.chart-wrap{position:relative;height:240px;margin-top:12px}

/* ===== KEY THEMES ===== */
.themes{display:flex;flex-wrap:wrap;gap:8px;margin-top:8px}
.theme-tag{background:rgba(124,106,247,.15);color:var(--accent);border:1px solid rgba(124,106,247,.3);
  padding:4px 12px;border-radius:20px;font-size:.8rem;font-weight:500}

/* ===== TOGGLE ===== */
.chevron{transition:transform .25s;font-size:.8rem;color:var(--muted);margin-left:auto}
.finding-hdr.open .chevron{transform:rotate(180deg)}

/* ===== TABLE ===== */
.tbl{width:100%;border-collapse:collapse;font-size:.82rem;margin-top:8px}
.tbl th{background:var(--surface2);color:var(--muted);font-weight:500;text-align:left;
  padding:8px 10px;border-bottom:1px solid var(--border)}
.tbl td{padding:8px 10px;border-bottom:1px solid rgba(45,49,72,.5);color:var(--text);
  font-family:monospace;white-space:nowrap;max-width:300px;overflow:hidden;text-overflow:ellipsis}
.tbl tr:hover td{background:var(--surface2)}

/* ===== FOOTER ===== */
.footer{text-align:center;padding:40px 0 20px;color:var(--muted);font-size:.82rem}
.footer strong{color:var(--accent)}

/* ===== SECTION TITLE ===== */
.section-title{font-size:1.4rem;font-weight:700;margin-bottom:20px;
  display:flex;align-items:center;gap:10px;letter-spacing:-.02em}
.section-title .dot{width:8px;height:8px;border-radius:50%;background:var(--accent);flex-shrink:0}

/* ===== FILTERS ===== */
.filters{display:flex;gap:8px;flex-wrap:wrap;margin-bottom:16px}
.filter-btn{background:var(--surface2);border:1px solid var(--border);color:var(--muted);
  padding:6px 16px;border-radius:20px;cursor:pointer;font-size:.82rem;font-weight:500;transition:.2s}
.filter-btn:hover,.filter-btn.active{background:var(--accent);color:#fff;border-color:var(--accent)}

/* ===== SCROLL ANIMATION ===== */
.fade-up{opacity:0;transform:translateY(16px);transition:opacity .4s,transform .4s}
.fade-up.visible{opacity:1;transform:none}
</style>
</head>
<body>

<!-- HEADER -->
<div class="header">
  <h1>⚡ TriageProf Performance Report</h1>
  <p class="sub">Deterministic profiling analysis · Powered by Mistral AI</p>
  <div class="meta">
    <span><strong>Generated:</strong> {{.GeneratedAt}}</span>
    <span><strong>Findings:</strong> {{len .Findings.Findings}}</span>
    {{if insightsEnabled .Insights}}<span>🤖 <strong>AI:</strong> {{.Insights.Model}}</span>{{end}}
  </div>
</div>

<div class="wrap">

<!-- ===== OVERVIEW ROW ===== -->
<div class="grid2 fade-up">
  <!-- Score card -->
  <div class="card">
    <h2><span class="icon">🎯</span> Health Score</h2>
    <div style="display:flex;align-items:center;gap:32px;flex-wrap:wrap">
      <div class="gauge-wrap">
        <div class="gauge-ring" id="gaugeRing">
          <svg width="120" height="120" viewBox="0 0 120 120">
            <circle cx="60" cy="60" r="52" fill="none" stroke="#22263a" stroke-width="10"/>
            <circle id="gaugeFill" cx="60" cy="60" r="52" fill="none"
              stroke="url(#gGrad)" stroke-width="10"
              stroke-linecap="round" stroke-dasharray="327" stroke-dashoffset="327"/>
            <defs>
              <linearGradient id="gGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" style="stop-color:#7c6af7"/>
                <stop offset="100%" style="stop-color:#56cfb2"/>
              </linearGradient>
            </defs>
          </svg>
          <div class="gauge-num" id="gaugeNum">—</div>
        </div>
        <div class="gauge-label">out of 100</div>
      </div>
      <div style="flex:1;min-width:180px">
        <div style="font-size:.9rem;color:var(--muted);margin-bottom:12px">Overall analysis quality</div>
        {{with .Findings.Summary}}
        {{if .TopIssueTags}}
        <div style="font-size:.82rem;color:var(--muted);margin-bottom:6px">Top Issues</div>
        <div class="themes">
          {{range .TopIssueTags}}<span class="theme-tag">{{.}}</span>{{end}}
        </div>
        {{end}}
        {{if .Notes}}
        <div style="margin-top:12px">
          {{range .Notes}}<div style="font-size:.82rem;color:var(--muted)">• {{.}}</div>{{end}}
        </div>
        {{end}}
        {{end}}
      </div>
    </div>
  </div>

  <!-- Stat tiles -->
  <div style="display:flex;flex-direction:column;gap:12px">
    <div class="grid2">
      <div class="stat critical"><div class="val" id="sCritical">0</div><div class="lbl">Critical</div></div>
      <div class="stat high"><div class="val" id="sHigh">0</div><div class="lbl">High</div></div>
    </div>
    <div class="grid2">
      <div class="stat medium"><div class="val" id="sMedium">0</div><div class="lbl">Medium</div></div>
      <div class="stat low"><div class="val" id="sLow">0</div><div class="lbl">Low</div></div>
    </div>
  </div>
</div>

<!-- ===== AI EXECUTIVE SUMMARY ===== -->
{{if insightsEnabled .Insights}}
{{with .Insights}}
<div class="card fade-up" style="margin-top:20px;background:linear-gradient(135deg,rgba(124,106,247,.08),rgba(86,207,178,.05));border-color:rgba(124,106,247,.35)">
  <h2>🤖 Mistral AI Executive Summary
    <span class="ai-badge" style="margin-left:auto">{{.Model}}</span>
  </h2>
  <p style="font-size:.95rem;line-height:1.75;color:var(--text)">{{.ExecutiveSummary.Overview}}</p>

  {{if .ExecutiveSummary.KeyThemes}}
  <div class="themes" style="margin-top:14px">
    {{range .ExecutiveSummary.KeyThemes}}<span class="theme-tag">{{.}}</span>{{end}}
  </div>
  {{end}}

  <div class="grid3" style="margin-top:20px">
    <div style="background:rgba(0,0,0,.2);border-radius:8px;padding:14px;text-align:center">
      <div style="font-size:1.4rem;font-weight:700;color:var(--accent)">{{.ExecutiveSummary.OverallSeverity | title}}</div>
      <div style="font-size:.78rem;color:var(--muted);margin-top:4px">Overall Severity</div>
    </div>
    <div style="background:rgba(0,0,0,.2);border-radius:8px;padding:14px;text-align:center">
      <div style="font-size:1.4rem;font-weight:700;color:var(--accent2)">{{.ExecutiveSummary.Confidence}}%</div>
      <div style="font-size:.78rem;color:var(--muted);margin-top:4px">AI Confidence</div>
    </div>
    <div style="background:rgba(0,0,0,.2);border-radius:8px;padding:14px;text-align:center">
      <div style="font-size:1.4rem;font-weight:700;color:#f6e05e">{{len .PerFinding}}</div>
      <div style="font-size:.78rem;color:var(--muted);margin-top:4px">Findings Analyzed</div>
    </div>
  </div>
</div>
{{end}}
{{end}}

<!-- ===== CHARTS ROW ===== -->
<div class="grid2 fade-up" style="margin-top:20px">
  <div class="card">
    <h2><span class="icon">📊</span> Severity Distribution</h2>
    <div class="chart-wrap"><canvas id="chartSeverity"></canvas></div>
  </div>
  <div class="card">
    <h2><span class="icon">📂</span> Category Breakdown</h2>
    <div class="chart-wrap"><canvas id="chartCategory"></canvas></div>
  </div>
</div>

<!-- ===== TOP RISKS ===== -->
{{if insightsEnabled .Insights}}
{{if .Insights.TopRisks}}
<div class="card fade-up" style="margin-top:20px">
  <h2>🚨 Top Risks</h2>
  <div class="risk-list">
    {{range .Insights.TopRisks}}
    <div class="risk {{.Severity | lower}}">
      <span style="font-size:1.2rem">{{.Severity | lower | severityEmoji}}</span>
      <div class="risk-body">
        <div class="risk-title">{{.Description}}</div>
        <div class="risk-detail">{{if .Impact}}Impact: {{.Impact}}{{end}}{{if .PotentialImpact}} · {{.PotentialImpact}}{{end}}</div>
      </div>
      <div style="margin-left:auto;flex-shrink:0"><span class="badge {{.Severity | lower}}">{{.Severity}}</span></div>
    </div>
    {{end}}
  </div>
</div>
{{end}}
{{end}}

<!-- ===== RECOMMENDATIONS ===== -->
{{if insightsEnabled .Insights}}
{{if .Insights.TopActions}}
<div class="card fade-up" style="margin-top:20px">
  <h2>💡 Recommendations</h2>
  <div class="reco-grid">
    {{range $i, $a := .Insights.TopActions}}
    <div class="reco">
      <div class="reco-hdr">
        <div class="reco-num {{$a.Priority | lower}}">{{add $i 1}}</div>
        <div class="reco-title-block">
          <div class="reco-title">{{$a.Description}}</div>
          <div class="reco-meta">
            <span class="badge {{$a.Priority | lower}}" style="font-size:.7rem">{{$a.Priority}} priority</span>
            {{if $a.EstimatedEffort}}<span class="reco-tag effort">⏱ {{$a.EstimatedEffort}}</span>{{end}}
            {{if $a.ImplementationComplexity}}<span class="reco-tag complexity-{{$a.ImplementationComplexity}}">⚙ {{$a.ImplementationComplexity}} complexity</span>{{end}}
            {{range $a.Categories}}<span class="reco-tag category">{{.}}</span>{{end}}
          </div>
        </div>
      </div>
      {{if $a.ExpectedImpact}}
      <div class="reco-impact-bar">
        <span class="reco-impact-label">📈 Expected impact</span>
        <span class="reco-impact-val">{{$a.ExpectedImpact}}</span>
      </div>
      {{end}}
      {{if $a.ValidationMetrics}}
      <div class="reco-validation">
        <div class="reco-validation-title">✅ How to validate</div>
        <div class="reco-validation-list">
          {{range $a.ValidationMetrics}}<span class="reco-validation-item">{{.}}</span>{{end}}
        </div>
      </div>
      {{end}}
      {{range $a.CodeExamples}}
      <div class="reco-code-block">{{.}}</div>
      {{end}}
    </div>
    {{end}}
  </div>
</div>
{{end}}
{{end}}

<!-- ===== DETAILED FINDINGS ===== -->
<div style="margin-top:32px" class="fade-up">
  <div class="section-title"><span class="dot"></span> Detailed Findings</div>

  <div class="filters" id="filterBar">
    <button class="filter-btn active" data-f="all">All</button>
    <button class="filter-btn" data-f="critical">🔴 Critical</button>
    <button class="filter-btn" data-f="high">🟠 High</button>
    <button class="filter-btn" data-f="medium">🟡 Medium</button>
    <button class="filter-btn" data-f="low">🟢 Low</button>
  </div>

  <div id="findingsList">
  {{range $i, $f := .Findings.Findings}}
  <div class="finding" data-sev="{{$f.Severity | lower}}">
    <div class="finding-hdr" onclick="toggleFinding(this)">
      <div class="emoji">{{$f.Category | profileEmoji}}</div>
      <div class="info">
        <div class="title">{{$f.Title}}</div>
        <div class="sub">{{$f.Category | title}} · {{$f.ID}}</div>
      </div>
      <div class="right">
        <span class="badge {{$f.Severity | lower}}">{{$f.Severity}}</span>
        <span class="finding-score">Score {{$f.Score}}</span>
        <span class="chevron">▼</span>
      </div>
    </div>
    <div class="finding-body">

      <!-- HOTSPOTS -->
      {{if $f.Top}}
      <div style="margin-bottom:24px">
        <div style="font-size:.82rem;text-transform:uppercase;letter-spacing:.06em;color:var(--muted);margin-bottom:10px">🔥 Top Hotspots</div>
        <div class="hotspot-list">
          {{$maxVal := 1.0}}
          {{range $f.Top}}{{if gt .Cum $maxVal}}{{end}}{{end}}
          {{range $j, $fr := $f.Top}}{{if lt $j 8}}
          <div class="hotspot">
            <div class="hotspot-meta">
              <span class="hotspot-fn" title="{{$fr.Function}}">{{$fr.Function | truncate 60}}</span>
              <span class="hotspot-val">{{printf "%.0f" $fr.Cum}}</span>
            </div>
            <div class="hotspot-bar">
              <div class="hotspot-fill" style="width:0%" data-target="{{printf "%.2f" $fr.Cum}}"></div>
            </div>
          </div>
          {{end}}{{end}}
        </div>
        <!-- full table toggle -->
        {{if gt (len $f.Top) 8}}
        <details style="margin-top:12px">
          <summary style="cursor:pointer;font-size:.82rem;color:var(--accent)">Show all {{len $f.Top}} frames</summary>
          <div style="overflow-x:auto;margin-top:8px">
          <table class="tbl">
            <thead><tr><th>Function</th><th>File</th><th>Line</th><th>Cum</th><th>Flat</th></tr></thead>
            <tbody>
              {{range $f.Top}}<tr>
                <td title="{{.Function}}">{{.Function | truncate 50}}</td>
                <td title="{{.File}}">{{.File | truncate 40}}</td>
                <td>{{.Line}}</td>
                <td>{{printf "%.0f" .Cum}}</td>
                <td>{{printf "%.0f" .Flat}}</td>
              </tr>{{end}}
            </tbody>
          </table>
          </div>
        </details>
        {{end}}
      </div>
      {{end}}

      <!-- ALLOCATION ANALYSIS -->
      {{if $f.AllocationAnalysis}}
      {{with $f.AllocationAnalysis}}
      <div style="margin-bottom:24px">
        <div style="font-size:.82rem;text-transform:uppercase;letter-spacing:.06em;color:var(--muted);margin-bottom:10px">🧠 Allocation Analysis</div>
        <div class="grid3">
          <div class="stat accent"><div class="val" style="font-size:1.4rem">{{printf "%.0f" .TotalAllocations}}</div><div class="lbl">Total Allocs</div></div>
          <div class="stat accent"><div class="val" style="font-size:1.4rem">{{printf "%.1f" .TopConcentration}}%</div><div class="lbl">Top Concentration</div></div>
          <div class="stat"><div class="val" style="font-size:1.4rem"><span class="badge {{.Severity | lower}}">{{.Severity}}</span></div><div class="lbl">Severity</div></div>
        </div>
        {{if .Hotspots}}
        <div style="overflow-x:auto;margin-top:12px">
        <table class="tbl">
          <thead><tr><th>Function</th><th>File</th><th>Count</th><th>%</th></tr></thead>
          <tbody>
            {{range .Hotspots}}<tr>
              <td title="{{.Function}}">{{.Function | truncate 45}}</td>
              <td title="{{.File}}">{{.File | truncate 30}}</td>
              <td>{{printf "%.0f" .Count}}</td>
              <td>{{printf "%.1f" .Percent}}%</td>
            </tr>{{end}}
          </tbody>
        </table>
        </div>
        {{end}}
      </div>
      {{end}}
      {{end}}

      <!-- AI INSIGHTS FOR THIS FINDING -->
      {{if hasInsight $.InsightMap $f.ID}}
      {{$ins := getInsight $.InsightMap $f.ID}}
      <div class="ai-card">
        <div class="ai-header">
          <span style="font-size:1.1rem">🤖</span>
          <h3>Mistral AI Analysis</h3>
          {{if ge $ins.Confidence 80}}<span class="confidence-pill confidence-high">{{$ins.Confidence}}% confidence</span>
          {{else if ge $ins.Confidence 50}}<span class="confidence-pill confidence-mid">{{$ins.Confidence}}% confidence</span>
          {{else}}<span class="confidence-pill confidence-low">{{$ins.Confidence}}% confidence</span>{{end}}
          <div style="margin-left:auto;display:flex;gap:8px;flex-wrap:wrap">
            {{if $ins.PerformanceImpact}}<span class="reco-tag effort">📈 {{$ins.PerformanceImpact}}</span>{{end}}
            {{if $ins.ImplementationComplexity}}<span class="reco-tag complexity-{{$ins.ImplementationComplexity}}">⚙ {{$ins.ImplementationComplexity}} complexity</span>{{end}}
          </div>
        </div>

        {{if $ins.Narrative}}
        <div class="ai-section">
          <h4>Analysis</h4>
          <div class="ai-narrative">{{$ins.Narrative}}</div>
        </div>
        {{end}}

        <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:14px">
          {{if $ins.LikelyRootCauses}}
          <div class="ai-section" style="margin-top:0">
            <h4>🔍 Likely Root Causes</h4>
            <ul class="ai-list">{{range $ins.LikelyRootCauses}}<li>{{.}}</li>{{end}}</ul>
          </div>
          {{end}}
          {{if $ins.Suggestions}}
          <div class="ai-section" style="margin-top:0">
            <h4>💡 Suggestions</h4>
            <ul class="ai-list">{{range $ins.Suggestions}}<li>{{.}}</li>{{end}}</ul>
          </div>
          {{end}}
        </div>

        {{if $ins.CodeExamples}}
        <div class="ai-section">
          <h4>📝 Code Examples</h4>
          {{range $ins.CodeExamples}}<div class="code-block">{{.}}</div>{{end}}
        </div>
        {{end}}

        <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:14px">
          {{if $ins.NextMeasurements}}
          <div class="ai-section" style="margin-top:0">
            <h4>📊 Next Measurements</h4>
            <ul class="ai-list">{{range $ins.NextMeasurements}}<li>{{.}}</li>{{end}}</ul>
          </div>
          {{end}}
          {{if $ins.Caveats}}
          <div class="ai-section" style="margin-top:0">
            <h4>⚠️ Caveats</h4>
            <ul class="ai-list">{{range $ins.Caveats}}<li>{{.}}</li>{{end}}</ul>
          </div>
          {{end}}
        </div>

        {{if $ins.BeforeAfterMetrics}}
        <div class="ai-section">
          <h4>� Before / After Metrics</h4>
          <div class="reco-validation-list">
            {{range $ins.BeforeAfterMetrics}}<span class="reco-validation-item">{{.}}</span>{{end}}
          </div>
        </div>
        {{end}}

      </div>
      {{end}}

    </div><!-- /finding-body -->
  </div>
  {{end}}
  </div><!-- /findingsList -->
</div>

</div><!-- /wrap -->

<div class="footer">
  Generated by <strong>TriageProf</strong> · Powered by <strong>Mistral AI</strong>
</div>

<script>
// ===== DATA =====
const FINDINGS = {{.FindingsJSON}};
const INSIGHTS = {{.InsightsJSON}};

// ===== GAUGE =====
(function(){
  const score = FINDINGS.summary && FINDINGS.summary.overall_score || 0;
  const circ = 2 * Math.PI * 52; // 326.7
  const offset = circ - (score / 100) * circ;
  const fill = document.getElementById('gaugeFill');
  const num = document.getElementById('gaugeNum');
  if(fill){ setTimeout(()=>{ fill.style.strokeDashoffset = offset; fill.style.transition='stroke-dashoffset 1s ease'; },100); }
  if(num){ num.textContent = score; }
})();

// ===== SEVERITY STATS =====
(function(){
  const counts = {critical:0,high:0,medium:0,low:0};
  (FINDINGS.findings||[]).forEach(f=>{ const s=(f.severity||'low').toLowerCase(); if(counts[s]!==undefined)counts[s]++; });
  document.getElementById('sCritical').textContent=counts.critical;
  document.getElementById('sHigh').textContent=counts.high;
  document.getElementById('sMedium').textContent=counts.medium;
  document.getElementById('sLow').textContent=counts.low;
})();

// ===== CHARTS =====
(function(){
  const isDark = true;
  const textColor='#8892a4', gridColor='rgba(255,255,255,.06)';
  const baseOpts = { responsive:true, maintainAspectRatio:false,
    plugins:{ legend:{ labels:{ color:textColor, boxWidth:14, padding:16 } } } };

  // Severity donut
  const counts = {critical:0,high:0,medium:0,low:0};
  (FINDINGS.findings||[]).forEach(f=>{ const s=(f.severity||'low').toLowerCase(); if(counts[s]!==undefined)counts[s]++; });
  new Chart(document.getElementById('chartSeverity'),{
    type:'doughnut',
    data:{ labels:['Critical','High','Medium','Low'],
      datasets:[{ data:[counts.critical,counts.high,counts.medium,counts.low],
        backgroundColor:['#fc8181','#f6ad55','#f6e05e','#68d391'], borderWidth:0, hoverOffset:6 }] },
    options:{...baseOpts, cutout:'65%', plugins:{...baseOpts.plugins,
      tooltip:{ callbacks:{ label:(c)=>' '+c.label+': '+c.raw } } } }
  });

  // Category bar
  const cats={};
  (FINDINGS.findings||[]).forEach(f=>{ const c=f.category||'other'; cats[c]=(cats[c]||0)+1; });
  const colors=['#7c6af7','#56cfb2','#f6ad55','#fc8181','#f6e05e','#68d391'];
  new Chart(document.getElementById('chartCategory'),{
    type:'bar',
    data:{ labels:Object.keys(cats),
      datasets:[{ label:'Findings', data:Object.values(cats),
        backgroundColor:Object.keys(cats).map((_,i)=>colors[i%colors.length]),
        borderRadius:6, borderWidth:0 }] },
    options:{...baseOpts, indexAxis:'y',
      plugins:{...baseOpts.plugins, legend:{display:false}},
      scales:{ x:{ ticks:{color:textColor}, grid:{color:gridColor}, beginAtZero:true },
               y:{ ticks:{color:textColor}, grid:{color:gridColor} } } }
  });
})();

// ===== HOTSPOT BARS ANIMATION =====
(function(){
  document.querySelectorAll('.finding').forEach(card=>{
    const bars = card.querySelectorAll('.hotspot-fill');
    if(!bars.length) return;
    let maxVal = 0;
    bars.forEach(b=>{ const v=parseFloat(b.dataset.target||0); if(v>maxVal)maxVal=v; });
    if(!maxVal) return;
    // Animate when card is opened
    card.addEventListener('click', ()=>{
      setTimeout(()=>{
        bars.forEach(b=>{
          const v = parseFloat(b.dataset.target||0);
          b.style.width = Math.max(2, (v/maxVal)*100)+'%';
        });
      }, 50);
    });
  });
})();

// ===== TOGGLE FINDINGS =====
function toggleFinding(hdr){
  hdr.classList.toggle('open');
  const body = hdr.nextElementSibling;
  body.classList.toggle('open');
}

// ===== FILTER =====
(function(){
  document.getElementById('filterBar').addEventListener('click',e=>{
    const btn = e.target.closest('.filter-btn');
    if(!btn) return;
    document.querySelectorAll('.filter-btn').forEach(b=>b.classList.remove('active'));
    btn.classList.add('active');
    const f = btn.dataset.f;
    document.querySelectorAll('#findingsList .finding').forEach(card=>{
      card.style.display = (f==='all'||card.dataset.sev===f)?'':'none';
    });
  });
})();

// ===== SCROLL FADE =====
(function(){
  const obs = new IntersectionObserver(entries=>{
    entries.forEach(e=>{ if(e.isIntersecting) e.target.classList.add('visible'); });
  },{threshold:0.05});
  document.querySelectorAll('.fade-up').forEach(el=>obs.observe(el));
})();
</script>
</body>
</html>
`
