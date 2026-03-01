# CI/CD Integration Templates

This directory contains templates and examples for integrating TriageProf into your CI/CD pipelines.

## Available Templates

### GitHub Actions

- **[github-actions-performance.yml](github-actions-performance.yml)** - Complete GitHub Actions workflow for performance analysis with configurable gates

## Quick Start

### GitHub Actions Integration

1. **Copy the template** to your repository:
   ```bash
   cp docs/ci-cd-templates/github-actions-performance.yml .github/workflows/performance.yml
   ```

2. **Customize thresholds** in the workflow file:
   ```yaml
   # Configure your performance thresholds here
   CRITICAL_THRESHOLD=3
   HIGH_THRESHOLD=8
   MEDIUM_THRESHOLD=15
   ```

3. **Commit and push** the workflow file to trigger performance analysis on every push/PR.

## Configuration Options

### Performance Gates

Configure performance thresholds based on your project requirements:

- **Critical Findings Threshold**: Maximum allowed critical severity findings
- **High Findings Threshold**: Maximum allowed high severity findings  
- **Medium Findings Threshold**: Maximum allowed medium severity findings (warning only)

### Customization

You can customize the workflow to:

- Run on specific branches only
- Adjust analysis duration
- Change performance thresholds
- Add additional reporting steps
- Integrate with other CI/CD tools

## Outputs

The workflow generates:

1. **Performance Artifacts**: Complete analysis output including profiles and findings
2. **Performance Report**: Markdown report with findings summary and recommendations
3. **GitHub Annotations**: Inline annotations showing performance gate status

## Advanced Usage

### Using with Your Own Codebase

To analyze your own Go project instead of the demo kit:

```yaml
- name: Run performance analysis
  run: |
    ./triageprof demo --repo . --out performance-output --duration 30
```

### Custom Performance Gates

Add custom gate logic for specific performance metrics:

```yaml
- name: Custom performance gates
  run: |
    # Example: Check for specific bottleneck patterns
    if grep -q "JSON hotspot" performance-output/findings.json; then
      echo "::warning::JSON hotspot detected - consider optimization"
    fi
```

## Troubleshooting

### Common Issues

1. **No findings generated**: Ensure your project has Go benchmarks or use `--duration` flag
2. **Performance gates failing**: Adjust thresholds or optimize critical code paths
3. **Build failures**: Check Go version compatibility and dependencies

### Debugging

Add debug steps to your workflow:

```yaml
- name: Debug performance output
  if: failure()
  run: |
    echo "Debugging performance analysis..."
    ls -la performance-output/
    cat performance-output/findings.json | jq '.findings[0:5]'  # Show first 5 findings
```

## Integration with Other CI/CD Systems

### GitLab CI/CD

The same principles apply - convert the GitHub Actions workflow to GitLab CI/CD syntax.

### Jenkins

Use shell scripts with similar logic in your Jenkins pipeline.

### CircleCI

Adapt the workflow steps to CircleCI configuration format.

## Best Practices

1. **Start with lenient thresholds** and tighten them over time
2. **Run performance analysis on main branch only** for PRs to avoid duplicate runs
3. **Cache dependencies** to speed up workflow execution
4. **Use artifacts** to preserve performance data for debugging
5. **Monitor trends** over time to catch performance regressions early