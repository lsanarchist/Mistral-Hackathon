# Contributing to TriageProf

## Table of Contents

- [Welcome](#welcome)
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Plugin Development](#plugin-development)
- [Issue Reporting](#issue-reporting)
- [Feature Requests](#feature-requests)
- [Pull Request Process](#pull-request-process)
- [Community](#community)

## Welcome

Thank you for your interest in contributing to TriageProf! We welcome contributions from everyone, regardless of experience level. This guide will help you get started.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please read it to understand the expected behavior.

## Getting Started

### Prerequisites

- Go 1.20+
- Git
- Make
- Docker (for testing)

### Setup

```bash
# Clone the repository
git clone https://github.com/triageprof/triageprof.git
cd triageprof

# Build the project
make build

# Run tests
go test ./...
```

### Project Structure

```
.
├── cmd/                # CLI commands
├── internal/            # Core libraries
│   ├── analyzer/       # Analysis logic
│   ├── core/            # Core pipeline
│   ├── llm/             # LLM integration
│   ├── model/           # Data models
│   ├── plugin/          # Plugin system
│   ├── report/          # Report generation
│   └── webserver/       # Web server
├── plugins/            # Plugin system
│   ├── bin/            # Plugin binaries
│   ├── manifests/      # Plugin manifests
│   └── src/            # Plugin source code
├── web/                # Web assets
├── docs/               # Documentation
└── examples/           # Example projects
```

## Development Workflow

### Branching Strategy

- `main`: Stable branch
- `dev`: Development branch
- `feature/*`: Feature branches
- `fix/*`: Bug fix branches
- `docs/*`: Documentation updates

### Creating a Feature Branch

```bash
# Create and checkout a new feature branch
git checkout -b feature/your-feature-name

# Make your changes
# Commit with descriptive messages
git commit -m "feat: add new analysis rule for CPU hotpaths"

# Push to origin
git push origin feature/your-feature-name
```

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>
[optional body]
[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Build process or auxiliary tool changes

**Examples:**
```
feat(analyzer): add GC pressure detection rule
fix(core): handle nil pointer in benchmark detection
docs: update user guide with troubleshooting section
test(llm): add validation tests for insights generation
```

## Coding Standards

### Go Standards

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `go vet` for static analysis
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Project-Specific Standards

1. **Error Handling**: Use structured error contexts
2. **Logging**: Use `log` package with context
3. **Configuration**: Use environment variables with defaults
4. **Testing**: Write table-driven tests
5. **Documentation**: Document all public functions

### Code Formatting

```bash
# Format all Go code
go fmt ./...

# Check for common mistakes
go vet ./...

# Run linter
golangci-lint run
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/analyzer/

# Run with coverage
go test -cover ./...

# Run with race detector
go test -race ./...
```

### Writing Tests

```go
func TestDeterministicAnalyzer(t *testing.T) {
    tests := []struct {
        name     string
        profile  string
        expected []Finding
        wantErr  bool
    }{
        {
            name: "CPU hotpath detection",
            profile: "testdata/cpu-hotpath.prof",
            expected: []Finding{
                {
                    ID: "cpu-hotpath-1",
                    Title: "CPU Hotpath in FunctionX",
                    Category: "cpu",
                    Severity: "high",
                    Confidence: 0.95,
                },
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            analyzer := NewDeterministicAnalyzer()
            findings, err := analyzer.Analyze(tt.profile)

            if (err != nil) != tt.wantErr {
                t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !reflect.DeepEqual(findings, tt.expected) {
                t.Errorf("Analyze() = %v, want %v", findings, tt.expected)
            }
        })
    }
}
```

### Test Coverage

Aim for 80%+ test coverage for new features.

```bash
# Check coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

### Documentation Standards

1. **User-Facing**: Clear, concise, and practical
2. **Developer-Facing**: Detailed and technical
3. **Code Comments**: Explain why, not what
4. **Examples**: Include practical usage examples

### Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | Project overview and quick start |
| `docs/USER_GUIDE.md` | Comprehensive user guide |
| `docs/API_DOCUMENTATION.md` | Plugin API documentation |
| `docs/CLI_REFERENCE.md` | CLI command reference |
| `docs/CONTRIBUTING.md` | Contribution guidelines |
| `COMPASS.md` | Project direction and architecture |
| `AGENTS.md` | Development guidelines |

### Updating Documentation

```bash
# Edit documentation files
vim docs/USER_GUIDE.md

# Preview Markdown
markdown-preview docs/USER_GUIDE.md

# Check for broken links
markdown-link-check docs/**/*.md
```

## Plugin Development

### Plugin Requirements

1. **JSON-RPC 2.0**: Must implement the protocol
2. **Stdio Communication**: Read/write via stdin/stdout
3. **Manifest**: Must have a valid manifest file
4. **Error Handling**: Proper JSON-RPC error codes

### Creating a Plugin

```bash
# Create plugin directory
mkdir plugins/src/my-plugin
cd plugins/src/my-plugin

# Initialize Go module
go mod init github.com/triageprof/my-plugin

# Implement JSON-RPC handler
# See docs/API_DOCUMENTATION.md for details
```

### Plugin Testing

```bash
# Test plugin initialization
triageprof plugins test my-plugin --method initialize

# Test profile collection
triageprof plugins test my-plugin --method collectProfile --params '{"profileType": "cpu"}'

# Test with real analysis
triageprof demo --repo ./test-app --plugin my-plugin --out analysis/
```

### Plugin Best Practices

1. **Error Handling**: Return proper JSON-RPC error codes
2. **Timeout Handling**: Respect timeout parameters
3. **Resource Cleanup**: Implement proper shutdown
4. **Validation**: Validate all input parameters
5. **Logging**: Use stderr for logging (not stdout)
6. **Performance**: Optimize for minimal overhead

## Issue Reporting

### Before Reporting

1. Check existing issues
2. Verify with latest version
3. Review documentation
4. Search for similar problems

### Creating an Issue

Use the issue template and provide:

1. **Description**: Clear problem description
2. **Steps to Reproduce**: Detailed reproduction steps
3. **Expected Behavior**: What should happen
4. **Actual Behavior**: What actually happens
5. **Environment**: OS, Go version, TriageProf version
6. **Logs**: Relevant log output
7. **Screenshots**: If applicable

### Issue Labels

| Label | Meaning |
|-------|---------|
| `bug` | Confirmed bug |
| `enhancement` | Feature request |
| `documentation` | Documentation issue |
| `good first issue` | Good for beginners |
| `help wanted` | Needs community help |
| `question` | Question or discussion |

## Feature Requests

### Before Requesting

1. Check existing feature requests
2. Review project roadmap in `COMPASS.md`
3. Consider implementing it yourself

### Creating a Feature Request

Provide:

1. **Problem**: What problem does this solve?
2. **Solution**: Proposed solution
3. **Alternatives**: Alternative approaches considered
4. **Impact**: Who benefits from this feature?
5. **Examples**: Usage examples if applicable

## Pull Request Process

### Before Submitting

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Write tests
5. Update documentation
6. Run all tests

### Submitting a Pull Request

1. **Title**: Clear and descriptive
2. **Description**: Explain what and why
3. **Related Issues**: Reference related issues
4. **Checklist**: Complete all items
5. **Reviewers**: Assign appropriate reviewers

### Pull Request Template

```markdown
## Description

[Clear description of changes]

## Related Issues

Fixes #123
Resolves #456

## Changes Made

- [ ] Code changes
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Breaking changes documented

## Testing

[Description of testing performed]

## Checklist

- [ ] Code follows project standards
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
- [ ] Ready for review
```

### Review Process

1. **Initial Review**: Code quality and standards
2. **Functional Review**: Correctness and completeness
3. **Testing Review**: Test coverage and quality
4. **Documentation Review**: Completeness and accuracy
5. **Approval**: Ready to merge

## Community

### Ways to Contribute

1. **Code**: Implement features and fixes
2. **Documentation**: Improve docs and examples
3. **Testing**: Write tests and report bugs
4. **Plugins**: Develop new plugins
5. **Support**: Help others in discussions
6. **Feedback**: Provide constructive feedback

### Communication

- **GitHub Issues**: For bug reports and feature requests
- **GitHub Discussions**: For questions and ideas
- **Pull Requests**: For code contributions

### Recognition

All contributors are recognized in:
- `CONTRIBUTORS.md` file
- Release notes
- GitHub contributors list

## Resources

- [Go Documentation](https://golang.org/doc/)
- [JSON-RPC 2.0 Spec](https://www.jsonrpc.org/specification)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Effective Go](https://golang.org/doc/effective_go.html)

## Thank You!

Your contributions help make TriageProf better for everyone. We appreciate your time and effort!
