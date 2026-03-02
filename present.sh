#!/usr/bin/env bash
# =============================================================================
#  TriageProf — Hackathon Video Presentation Script
#  Usage:  bash present.sh
#  Press ENTER to advance each step. Ctrl+C to quit at any time.
# =============================================================================

set -uo pipefail

# ── colours ──────────────────────────────────────────────────────────────────
BOLD=$'\e[1m'
DIM=$'\e[2m'
RESET=$'\e[0m'
CYAN=$'\e[36m'
GREEN=$'\e[32m'
YELLOW=$'\e[33m'
MAGENTA=$'\e[35m'
RED=$'\e[31m'
BLUE=$'\e[34m'

# ── helpers ───────────────────────────────────────────────────────────────────

# Print a section banner
banner() {
    local text="$1"
    local width=62
    local line
    line=$(printf '─%.0s' $(seq 1 $width))
    echo ""
    echo "${CYAN}${BOLD}┌${line}┐${RESET}"
    printf "${CYAN}${BOLD}│  %-${width}s│${RESET}\n" "$text"
    echo "${CYAN}${BOLD}└${line}┘${RESET}"
    echo ""
}

# Typewriter effect — fast enough to look live, not annoying
type_cmd() {
    echo -n "${GREEN}${BOLD}\$ ${RESET}"
    local text="$1"
    local i
    for (( i=0; i<${#text}; i++ )); do
        printf '%s' "${text:$i:1}"
        sleep 0.03
    done
    echo ""
}

# Print a dimmed comment line
comment() {
    echo "${DIM}# $1${RESET}"
}

# Wait for ENTER — show a subtle prompt
pause() {
    echo ""
    printf "${DIM}[ press ENTER to continue ]${RESET}"
    read -r _
}

# Run a command and stream its output
run_live() {
    eval "$1" || true
}

# Print a key/value info line
info() {
    printf "  ${CYAN}%-22s${RESET} %s\n" "$1" "$2"
}

# ── cleanup trap (registered early so it always fires) ───────────────────────
DEMO_PID=""
HTTP_PID=""
LOAD_PID=""
cleanup() {
    [[ -n "$LOAD_PID" ]] && kill "$LOAD_PID" 2>/dev/null || true
    [[ -n "$DEMO_PID" ]] && kill "$DEMO_PID" 2>/dev/null || true
    [[ -n "$HTTP_PID" ]] && kill "$HTTP_PID" 2>/dev/null || true
}
trap cleanup EXIT

# ── load API key ──────────────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if [[ -f apikey.swaga ]]; then
    export MISTRAL_API_KEY
    MISTRAL_API_KEY=$(cat apikey.swaga)
elif [[ -z "${MISTRAL_API_KEY:-}" ]]; then
    echo "${RED}ERROR: No API key found.${RESET}"
    echo "  Put your Mistral API key in ./apikey.swaga  or  export MISTRAL_API_KEY=..."
    exit 1
fi

DEMO_OUT="./demo-output"
DEMO_SERVER_URL="http://localhost:6060"

# ── pre-flight: build if needed ───────────────────────────────────────────────
if [[ ! -f ./bin/triageprof ]]; then
    echo "${YELLOW}Binary not found — building first...${RESET}"
    make build
fi

# ── clear cache and old output for a fresh run ────────────────────────────────
echo "${DIM}Clearing LLM cache and previous output...${RESET}"
rm -rf /tmp/triageprof-insights-cache/
rm -rf "${DEMO_OUT:?}"/*
mkdir -p "$DEMO_OUT"
echo "${GREEN}✓ Cache cleared${RESET}"
sleep 0.5

# =============================================================================
#  SLIDE 1 — Title
# =============================================================================
clear
echo ""
echo "${BOLD}${MAGENTA}"
cat << 'EOF'
  ████████╗██████╗ ██╗ █████╗  ██████╗ ███████╗██████╗ ██████╗  ██████╗ ███████╗
     ██╔══╝██╔══██╗██║██╔══██╗██╔════╝ ██╔════╝██╔══██╗██╔══██╗██╔═══██╗██╔════╝
     ██║   ██████╔╝██║███████║██║  ███╗█████╗  ██████╔╝██████╔╝██║   ██║█████╗
     ██║   ██╔══██╗██║██╔══██║██║   ██║██╔══╝  ██╔═══╝ ██╔══██╗██║   ██║██╔══╝
     ██║   ██║  ██║██║██║  ██║╚██████╔╝███████╗██║     ██║  ██║╚██████╔╝██║
     ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝
EOF
echo "${RESET}"
echo "${BOLD}        AI-Powered Go Performance Profiling  ·  Mistral Hackathon 2025${RESET}"
echo ""
echo "  ${DIM}Built with  Mistral AI  ·  Go 1.24  ·  pprof${RESET}"
echo ""
echo "  ${CYAN}What it does:${RESET}"
echo "    ${BOLD}1.${RESET} Collects CPU / heap / alloc / mutex / block profiles from any live Go service"
echo "    ${BOLD}2.${RESET} Runs deterministic bottleneck analysis — 8+ rule-based patterns, scored findings"
echo "    ${BOLD}3.${RESET} Enriches findings with ${BOLD}mistral-large-latest${RESET} — root causes, fix suggestions,"
echo "       effort estimates, before/after metrics"
echo "    ${BOLD}4.${RESET} Produces a self-contained interactive HTML report, then serves it in your browser"
echo ""
pause

# =============================================================================
#  SLIDE 2 — Project structure
# =============================================================================
clear
banner "📁  Project Structure"

comment "What's inside the repo"
echo ""
type_cmd "ls -1"
echo ""
run_live "ls -1"
echo ""
pause

# =============================================================================
#  SLIDE 3 — Build
# =============================================================================
clear
banner "🔨  Step 1 — Build"

comment "One command builds the main binary + all profiler plugins"
echo ""
type_cmd "make build"
echo ""
run_live "make build"
echo ""
echo "${GREEN}${BOLD}✓ Built:${RESET}"
info "Main binary"   "./bin/triageprof"
info "Go plugin"     "./plugins/bin/go-pprof-http"
echo ""
pause

# =============================================================================
#  SLIDE 4 — Available plugins
# =============================================================================
clear
banner "🔌  Step 2 — Plugins"

comment "triageprof uses a JSON-RPC plugin architecture"
comment "Plugins are separate executables — easy to extend"
echo ""
type_cmd "./bin/triageprof plugins"
echo ""
run_live "./bin/triageprof plugins" || true
echo ""
pause

# =============================================================================
#  SLIDE 5 — Start demo server
# =============================================================================
clear
banner "🚀  Step 3 — Start the demo Go service"

comment "A real Go HTTP server with intentional performance problems:"
comment "  · CPU hotspot in a tight hash loop"
comment "  · Allocation churn via large []byte creation"
comment "  · pprof endpoint exposed on :6060"
echo ""
type_cmd "examples/demo-server/main &"
echo ""

# Kill any existing demo server on :6060
pkill -f 'demo-server/main' 2>/dev/null || true
sleep 0.3

./examples/demo-server/main &
DEMO_PID=$!
echo "${DIM}  (PID $DEMO_PID)${RESET}"
sleep 1

# Verify it's up
if curl -sf "${DEMO_SERVER_URL}/debug/pprof/" -o /dev/null; then
    echo "${GREEN}${BOLD}✓ Demo server is up at ${DEMO_SERVER_URL}${RESET}"
else
    echo "${YELLOW}  Server may still be starting — continuing...${RESET}"
fi
echo ""
pause

# =============================================================================
#  SLIDE 6 — Generate load
# =============================================================================
clear
banner "📈  Step 4 — Generate continuous load"

comment "Start a background load loop — keeps the server busy during profiling"
echo ""
type_cmd "while true; do curl -sf ${DEMO_SERVER_URL}/api/process -o /dev/null; done &"
echo ""

# Continuous load loop hitting the real CPU-heavy endpoints
(while true; do
    curl -sf "${DEMO_SERVER_URL}/api/process"   -o /dev/null 2>/dev/null
    curl -sf "${DEMO_SERVER_URL}/api/analytics" -o /dev/null 2>/dev/null
    curl -sf "${DEMO_SERVER_URL}/api/search"    -o /dev/null 2>/dev/null
    curl -sf "${DEMO_SERVER_URL}/api/users"     -o /dev/null 2>/dev/null
done) &
LOAD_PID=$!

echo "${GREEN}${BOLD}✓ Load running in background (PID $LOAD_PID)${RESET}"
echo "${DIM}  CPU hotspot and allocation churn will show up clearly in profiles${RESET}"
echo ""
pause

# =============================================================================
#  SLIDE 7 — Run triageprof (collect + analyze only, fast)
# =============================================================================
clear
banner "🔍  Step 5 — Collect & Analyse  (deterministic, no LLM yet)"

comment "Collect 10s of profiles, run rule-based analysis, produce findings.json"
echo ""
type_cmd "./bin/triageprof run --plugin go-pprof-http --target-url ${DEMO_SERVER_URL} --duration 10 --outdir ${DEMO_OUT}"
echo ""
mkdir -p "$DEMO_OUT"
run_live "./bin/triageprof run --plugin go-pprof-http --target-url ${DEMO_SERVER_URL} --duration 10 --outdir ${DEMO_OUT}" || true
echo ""
pause

# =============================================================================
#  SLIDE 8 — Show raw findings
# =============================================================================
clear
banner "📋  Step 6 — Findings  (deterministic, zero hallucination)"

comment "Pure pprof-backed findings — specific functions, real percentages"
echo ""
type_cmd "cat ${DEMO_OUT}/findings.json | python3 -m json.tool | head -60"
echo ""
python3 -m json.tool "${DEMO_OUT}/findings.json" 2>/dev/null | head -60 || true
echo "${DIM}  ... (truncated for display)${RESET}"
echo ""

# Stop the load generator now — profiling is done
[[ -n "$LOAD_PID" ]] && kill "$LOAD_PID" 2>/dev/null || true
LOAD_PID=""

pause

# =============================================================================
#  SLIDE 9 — Run with LLM
# =============================================================================
clear
banner "🧠  Step 7 — Mistral AI Enrichment"

comment "Now pass the findings to mistral-large-latest:"
comment "  · root cause analysis per finding"
comment "  · prioritised fix recommendations with effort + complexity"
comment "  · code examples, before/after metrics, validation steps"
comment "  · executive summary with confidence score"
echo ""
type_cmd "./bin/triageprof run --plugin go-pprof-http --target-url ${DEMO_SERVER_URL} --duration 10 --outdir ${DEMO_OUT} --llm --llm-timeout 90"
echo ""

# Run in background, capture output to tmp file
LLM_LOG=$(mktemp /tmp/triageprof-llm-XXXXXX.log)
./bin/triageprof run \
    --plugin go-pprof-http \
    --target-url "${DEMO_SERVER_URL}" \
    --duration 10 \
    --outdir "${DEMO_OUT}" \
    --llm \
    --llm-timeout 90 \
    >"$LLM_LOG" 2>&1 &
LLM_RUN_PID=$!

# ── ASCII animation while we wait ────────────────────────────────────────────
FRAMES=(
"  ·  ·  ·"
"  ●  ·  ·"
"  ●  ●  ·"
"  ●  ●  ●"
"  ·  ●  ●"
"  ·  ·  ●"
)
BRAINFRAMES=(
"   (  ^  ^  )"
"   ( *  ^  )"
"   ( *  *  )"
"   ( ~  *  )"
"   ( ~  ~  )"
"   ( ^  ~  )"
)
STAGES=(
    "Collecting CPU profiles          "
    "Collecting heap profiles         "
    "Collecting alloc profiles        "
    "Running deterministic analysis   "
    "Sending findings to Mistral API  "
    "Waiting for mistral-large-latest "
    "Receiving AI insights            "
    "Parsing recommendations          "
    "Generating HTML report           "
)
STAGE_DELAYS=(4 3 3 3 5 30 20 5 5)

tput civis 2>/dev/null || true   # hide cursor

stage_idx=0
frame=0
elapsed=0

while kill -0 "$LLM_RUN_PID" 2>/dev/null; do
    stage=${STAGES[$stage_idx]}
    delay=${STAGE_DELAYS[$stage_idx]}

    fi=${FRAMES[$(( frame % ${#FRAMES[@]} ))]}
    bf=${BRAINFRAMES[$(( frame % ${#BRAINFRAMES[@]} ))]}

    printf "\r  ${CYAN}${BOLD}%s${RESET}  ${MAGENTA}%s${RESET}  ${DIM}%s${RESET}  " \
        "$fi" "$bf" "$stage"

    sleep 0.15
    frame=$(( frame + 1 ))
    elapsed=$(( elapsed + 1 ))

    # Advance stage label roughly on schedule
    if (( elapsed >= delay * 7 )) && (( stage_idx < ${#STAGES[@]} - 1 )); then
        stage_idx=$(( stage_idx + 1 ))
        elapsed=0
    fi
done

tput cnorm 2>/dev/null || true   # restore cursor
printf "\r%-80s\r" " "           # clear animation line

wait "$LLM_RUN_PID" || true

# Show captured output
cat "$LLM_LOG"
rm -f "$LLM_LOG"
# ─────────────────────────────────────────────────────────────────────────────

echo ""
pause

# =============================================================================
#  SLIDE 10 — Show insights
# =============================================================================
clear
banner "💡  Step 8 — Mistral AI Insights"

comment "insights.json — structured, grounded AI analysis:"
comment "  · references real function names from pprof data"
comment "  · never invents metrics — all numbers come from findings.json"
echo ""
type_cmd "cat ${DEMO_OUT}/insights.json | python3 -m json.tool | head -80"
echo ""
python3 -m json.tool "${DEMO_OUT}/insights.json" 2>/dev/null | head -80 || true
echo "${DIM}  ... (truncated for display)${RESET}"
echo ""
pause

# =============================================================================
#  SLIDE 11 — Show report size
# =============================================================================
clear
banner "📊  Step 9 — The Report"

comment "A single self-contained HTML file — no dependencies, no server needed"
echo ""
type_cmd "ls -lh ${DEMO_OUT}/report.html && wc -l ${DEMO_OUT}/report.html"
echo ""
ls -lh "${DEMO_OUT}/report.html"
wc -l "${DEMO_OUT}/report.html"
echo ""
echo "  ${CYAN}Contains:${RESET}"
echo "    ${BOLD}·${RESET} Overall health score gauge"
echo "    ${BOLD}·${RESET} Severity breakdown chart (critical / high / medium / low)"
echo "    ${BOLD}·${RESET} AI Executive Summary with confidence score"
echo "    ${BOLD}·${RESET} Top risks under load"
echo "    ${BOLD}·${RESET} Prioritised recommendations — effort, complexity, code examples, validation"
echo "    ${BOLD}·${RESET} Per-finding AI cards — root causes ↔ suggestions, before/after metrics"
echo ""
pause

# =============================================================================
#  SLIDE 12 — Serve report
# =============================================================================
clear
banner "🌐  Step 10 — Open the Report"

comment "triageprof auto-offers to serve the report after every run"
comment "Serving it now on a free local port..."
echo ""

# Find a free port
PORT=$(python3 -c "import socket; s=socket.socket(); s.bind(('',0)); print(s.getsockname()[1]); s.close()")
REPORT_URL="http://127.0.0.1:${PORT}/report.html"

type_cmd "# Serving ${DEMO_OUT}/ on ${REPORT_URL}"
echo ""

# Start background HTTP server
python3 -m http.server "$PORT" --directory "$DEMO_OUT" >/dev/null 2>&1 &
HTTP_PID=$!

# Wait until port is actually accepting connections (up to 3s)
for i in $(seq 1 30); do
    if curl -sf "http://127.0.0.1:${PORT}/" -o /dev/null 2>/dev/null; then
        break
    fi
    sleep 0.1
done

echo "${GREEN}${BOLD}✓ Report served at:${RESET}"
echo ""
echo "    ${BOLD}${BLUE}${REPORT_URL}${RESET}"
echo ""

# Open browser
xdg-open "$REPORT_URL" 2>/dev/null || open "$REPORT_URL" 2>/dev/null || true

echo "${DIM}  (browser opened — server stays alive until you press ENTER to exit at the end)${RESET}"
echo ""
pause

# =============================================================================
#  SLIDE 13 — Architecture diagram
# =============================================================================
clear
banner "🏗️   Architecture"

echo ""
echo "  ${BOLD}Input:${RESET}  any Go service exposing  ${CYAN}import _ \"net/http/pprof\"${RESET}"
echo ""
echo "  ${CYAN}┌─────────────────────────────────────────────────────────────┐${RESET}"
echo "  ${CYAN}│                      triageprof run                         │${RESET}"
echo "  ${CYAN}│                                                             │${RESET}"
echo "  ${CYAN}│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │${RESET}"
echo "  ${CYAN}│  │  go-pprof    │    │ Deterministic │    │  Mistral AI  │  │${RESET}"
echo "  ${CYAN}│  │  -http       │───▶│  Analyser    │───▶│  Enrichment  │  │${RESET}"
echo "  ${CYAN}│  │  (plugin)    │    │  8+ rules    │    │  mistral-    │  │${RESET}"
echo "  ${CYAN}│  │  JSON-RPC    │    │  scored      │    │  large-      │  │${RESET}"
echo "  ${CYAN}│  └──────────────┘    └──────────────┘    │  latest      │  │${RESET}"
echo "  ${CYAN}│                                           └──────┬───────┘  │${RESET}"
echo "  ${CYAN}│                                                  │          │${RESET}"
echo "  ${CYAN}│                                    ┌─────────────▼────────┐ │${RESET}"
echo "  ${CYAN}│                                    │  HTML Report         │ │${RESET}"
echo "  ${CYAN}│                                    │  findings.json       │ │${RESET}"
echo "  ${CYAN}│                                    │  insights.json       │ │${RESET}"
echo "  ${CYAN}│                                    │  report.md           │ │${RESET}"
echo "  ${CYAN}│                                    └──────────────────────┘ │${RESET}"
echo "  ${CYAN}└─────────────────────────────────────────────────────────────┘${RESET}"
echo ""
echo "  ${BOLD}Key property:${RESET} Mistral only adds ${BOLD}why/how${RESET} — all numbers come from real pprof data."
echo "  No hallucinated metrics. Grounded AI analysis."
echo ""
pause

# =============================================================================
#  SLIDE 14 — Closing
# =============================================================================
clear
banner "✅  Summary"

echo ""
echo "  ${BOLD}${GREEN}What we just saw:${RESET}"
echo ""
echo "    ${GREEN}✓${RESET}  ${BOLD}make build${RESET}                 — single command, binary + plugins ready"
echo "    ${GREEN}✓${RESET}  ${BOLD}triageprof plugins${RESET}          — extensible JSON-RPC plugin system"
echo "    ${GREEN}✓${RESET}  ${BOLD}triageprof run${RESET}              — collect, analyse, enrich, report"
echo "    ${GREEN}✓${RESET}  ${BOLD}mistral-large-latest${RESET}        — grounded root-cause analysis"
echo "    ${GREEN}✓${RESET}  ${BOLD}Self-contained HTML report${RESET}  — dark theme, charts, AI cards, ~300KB"
echo "    ${GREEN}✓${RESET}  ${BOLD}Auto browser serve${RESET}          — zero manual steps after profiling"
echo ""
echo "  ${CYAN}Repo:${RESET}  ${BOLD}https://github.com/lsanarchist/Mistral-Hackathon${RESET}"
echo ""
echo ""
echo "  ${MAGENTA}${BOLD}TriageProf — built for the Mistral AI Hackathon 🚀${RESET}"
echo ""

echo "${DIM}  (demo server and HTTP server will stop when this script exits)${RESET}"
echo ""
printf "${DIM}[ press ENTER to exit ]${RESET}"
read -r _
echo ""
