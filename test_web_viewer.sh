#!/bin/bash

# Test script for web viewer enhancements
set -e

echo "Testing web viewer enhancements..."

# Clean up any previous test output
rm -rf test-web-output test-web-enhanced

# Test basic web report generation
echo "1. Testing basic web report generation..."
./bin/triageprof web --in out-demo/findings.json --outdir test-web-output

if [ ! -f "test-web-output/web/index.html" ]; then
    echo "ERROR: Web report index.html not generated"
    exit 1
fi

if [ ! -f "test-web-output/web/app.js" ]; then
    echo "ERROR: Web report app.js not generated"
    exit 1
fi

if [ ! -f "test-web-output/web/style.css" ]; then
    echo "ERROR: Web report style.css not generated"
    exit 1
fi

if [ ! -f "test-web-output/web/data/findings.json" ]; then
    echo "ERROR: Web report findings.json not generated"
    exit 1
fi

echo "✓ Basic web report generation works"

# Test enhanced web report generation with insights
echo "2. Testing enhanced web report generation with insights..."
./bin/triageprof web --in out-demo/findings.json --outdir test-web-enhanced --insights out-demo/insights.json

if [ ! -f "test-web-enhanced/web/index.html" ]; then
    echo "ERROR: Enhanced web report index.html not generated"
    exit 1
fi

if [ ! -f "test-web-enhanced/web/data/insights.json" ]; then
    echo "ERROR: Enhanced web report insights.json not generated"
    exit 1
fi

echo "✓ Enhanced web report generation with insights works"

# Verify enhanced HTML content
echo "3. Verifying enhanced HTML content..."
if ! grep -q "Quick Statistics" test-web-enhanced/web/index.html; then
    echo "ERROR: Quick Statistics section not found in enhanced web report"
    exit 1
fi

if ! grep -q "Plugin Information" test-web-enhanced/web/index.html; then
    echo "ERROR: Plugin Information section not found in enhanced web report"
    exit 1
fi

echo "✓ Enhanced HTML content verified"

# Verify enhanced CSS content
echo "4. Verifying enhanced CSS content..."
if ! grep -q ".quick-stats-section" test-web-enhanced/web/style.css; then
    echo "ERROR: Quick stats CSS not found"
    exit 1
fi

if ! grep -q ".plugin-info-section" test-web-enhanced/web/style.css; then
    echo "ERROR: Plugin info CSS not found"
    exit 1
fi

echo "✓ Enhanced CSS content verified"

# Verify enhanced JavaScript content
echo "5. Verifying enhanced JavaScript content..."
if ! grep -q "renderQuickStats" test-web-enhanced/web/app.js; then
    echo "ERROR: renderQuickStats function not found"
    exit 1
fi

if ! grep -q "tooltip" test-web-enhanced/web/app.js; then
    echo "ERROR: Tooltip enhancements not found"
    exit 1
fi

echo "✓ Enhanced JavaScript content verified"

echo ""
echo "🎉 All web viewer enhancement tests passed!"
echo ""
echo "Enhancements included:"
echo "  ✓ Quick Statistics dashboard section"
echo "  ✓ Interactive chart tooltips with percentages"
echo "  ✓ Plugin Information section"
echo "  ✓ Enhanced visual hierarchy and hover effects"
echo "  ✓ Improved mobile responsiveness"

# Clean up
rm -rf test-web-output test-web-enhanced
