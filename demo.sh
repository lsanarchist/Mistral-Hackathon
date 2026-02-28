#!/bin/bash

# 🎬 TriageProf Demo Script - "Wow" Factor Demonstration
# This script showcases TriageProf's killer feature: AI-powered performance triage

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║   🎬 TRIAGEPROF DEMO: AI-Powered Performance Triage           ║"
echo "║   From Data to Insights in One Command!                       ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Configuration
DEMO_DIR="demo-output"
DURATION=20
CONCURRENCY=8
PLUGIN="go-pprof-http"
SERVER_URL="http://localhost:6060"

# Clean up previous demo
echo "🧹 Cleaning up previous demo output..."
rm -rf "$DEMO_DIR"
mkdir -p "$DEMO_DIR"
echo "✅ Demo directory ready: $DEMO_DIR"
echo ""

# Build the system
echo "🔨 Building TriageProf..."
make build > /dev/null 2>&1
echo "✅ Build complete"
echo ""

# Start demo server
echo "🚀 Starting enhanced demo server..."
cd examples/demo-server
go run main.go > server.log 2>&1 &
SERVER_PID=$!
echo "📝 Server PID: $SERVER_PID"
echo "⏳ Waiting for server to start..."
sleep 3
echo "✅ Server ready"
echo ""

# Generate load
echo "🔥 Generating realistic load..."
cd ../..
./examples/load.sh "$SERVER_URL" "$DURATION" "$CONCURRENCY"
echo ""

# 🎯 WOW MOMENT 1: Automatic Plugin Discovery
echo "🎯 WOW MOMENT 1: Automatic Plugin Discovery"
echo "─────────────────────────────────────────────"
echo ""
echo "📋 Available plugins:"
./bin/triageprof plugins
echo ""
echo "✨ Magic: Plugins are automatically discovered from manifests!"
echo ""

# 🎯 WOW MOMENT 2: End-to-End Workflow
echo "🎯 WOW MOMENT 2: Single Command Analysis"
echo "─────────────────────────────────────────────"
echo ""
echo "🚀 Running: bin/triageprof run --plugin $PLUGIN --target-url $SERVER_URL --duration $DURATION --outdir $DEMO_DIR"
echo ""

./bin/triageprof run --plugin "$PLUGIN" --target-url "$SERVER_URL" --duration "$DURATION" --outdir "$DEMO_DIR"

echo ""
echo "✅ Analysis complete! Files generated:"
ls -lh "$DEMO_DIR" | grep -E '\.(json|md|pb\.gz|txt)$'
echo ""

# 🎯 WOW MOMENT 3: Professional Report
echo "🎯 WOW MOMENT 3: Professional Markdown Report"
echo "────────────────────────────────────────────────"
echo ""
echo "📊 Executive Summary:"
echo "────────────────────────────────────────────────"
head -n 15 "$DEMO_DIR/report.md"
echo "..."
echo ""
echo "📈 Top Findings:"
echo "────────────────────────────────────────────────"
grep -A 5 "## " "$DEMO_DIR/report.md" | head -n 20
echo "..."
echo ""

# 🎯 WOW MOMENT 4: LLM Insights (Optional)
echo "🎯 WOW MOMENT 4: AI-Powered Insights (Optional)"
echo "─────────────────────────────────────────────────"
echo ""
echo "🤖 LLM augmentation explains WHY issues exist and HOW to fix them"
echo ""
echo "📝 To enable LLM insights, set MISTRAL_API_KEY and run:"
echo "   export MISTRAL_API_KEY='your-key-here'"
echo "   ./bin/triageprof run --plugin $PLUGIN --target-url $SERVER_URL --duration $DURATION --outdir $DEMO_DIR --llm"
echo ""
echo "💡 LLM provides:"
echo "   • Executive summary with severity assessment"
echo "   • Root cause analysis for each finding"
echo "   • Concrete suggestions with code examples"
echo "   • Confidence scores and caveats"
echo ""

# 🎯 WOW MOMENT 5: Comparison with Traditional Tools
echo "🎯 WOW MOMENT 5: Traditional vs. TriageProf"
echo "─────────────────────────────────────────────────"
echo ""
echo "🔴 Traditional Profilers:"
echo "   • Show WHAT is slow (hotspots)"
echo "   • Require manual analysis"
echo "   • No actionable insights"
echo ""
echo "🟢 TriageProf:"
echo "   • Shows WHAT is slow (deterministic analysis)"
echo "   • Explains WHY it's slow (LLM insights)"
echo "   • Suggests HOW to fix it (actionable recommendations)"
echo "   • Professional reports for stakeholders"
echo ""

# Verification
echo "✅ DEMO VERIFICATION"
echo "─────────────────────"
echo ""

# Check for expected files
MISSING_FILES=0
for file in bundle.json findings.json report.md cpu.pb.gz heap.pb.gz mutex.pb.gz block.pb.gz goroutine.txt; do
    if [ ! -f "$DEMO_DIR/$file" ]; then
        echo "❌ Missing: $file"
        MISSING_FILES=$((MISSING_FILES + 1))
    fi
done

if [ $MISSING_FILES -eq 0 ]; then
    echo "✅ All expected output files present"
else
    echo "❌ $MISSING_FILES files missing"
fi

# Check report content
if grep -q "Performance Triage Report" "$DEMO_DIR/report.md"; then
    echo "✅ Report contains expected header"
else
    echo "❌ Report header missing"
fi

if grep -q "Executive Summary" "$DEMO_DIR/report.md"; then
    echo "✅ Report contains executive summary"
else
    echo "❌ Executive summary missing"
fi

# Check findings
FINDING_COUNT=$(grep -c "## " "$DEMO_DIR/report.md" || echo "0")
if [ "$FINDING_COUNT" -gt 0 ]; then
    echo "✅ Found $FINDING_COUNT performance findings"
else
    echo "❌ No findings detected"
fi

# Cleanup
echo ""
echo "🧹 Cleaning up demo server..."
kill "$SERVER_PID" > /dev/null 2>&1 || true
echo "✅ Demo server stopped"

# Summary
echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                    🎉 DEMO COMPLETE!                          ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "📁 Results available in: $DEMO_DIR/"
echo ""
echo "📊 Key Achievements:"
echo "   ✅ Automatic plugin discovery and validation"
echo "   ✅ Complete profile collection (CPU, heap, mutex, block, goroutine)"
echo "   ✅ Deterministic performance analysis"
echo "   ✅ Professional markdown report generation"
echo "   ✅ Ready for LLM augmentation (optional)"
echo ""
echo "🚀 Next Steps:"
echo "   1. Review the report: cat $DEMO_DIR/report.md"
echo "   2. Try with LLM: export MISTRAL_API_KEY='your-key' && ./bin/triageprof run --llm"
echo "   3. Explore findings: cat $DEMO_DIR/findings.json | jq ."
echo "   4. Check bundle: cat $DEMO_DIR/bundle.json | jq ."
echo ""
echo "💡 Remember: TriageProf transforms raw data into actionable insights!"
echo ""
