# Contributing to CypherGoat CLI

Thank you for your interest in contributing to CypherGoat CLI! This document provides guidelines and instructions for contributing.

---

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Building](#building)
- [Submitting Changes](#submitting-changes)
- [Security](#security)
- [Questions](#questions)

---

## Getting Started

### Prerequisites

- Go 1.25 or later
- Task runner (https://taskfile.dev)
- Git

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/moralpriest/cli.git
cd cli

# Install Task if not already installed
task install:task

# Verify setup
task help
```

---

## Development Setup

### 1. Fork the Repository

1. Go to https://github.com/moralpriest/cli
2. Click the "Fork" button
3. Clone your fork:

```bash
git clone https://github.com/YOUR-USERNAME/cli.git
cd cli
git remote add upstream https://github.com/moralpriest/cli.git
```

### 2. Create a Feature Branch

```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/your-feature-name
```

### 3. Make Changes

Make your changes following the coding standards below.

---

## Coding Standards

### Go Best Practices

- Use `context.Context` for all API calls
- Use `any` instead of `interface{}`
- Use `slices.SortFunc` for sorting
- Use `slices` and `maps` packages from stdlib
- Wrap errors with `fmt.Errorf("...: %w", err)`
- Use meaningful variable names
- Keep functions small and focused

### Code Style

- Run `task lint` before committing
- Follow Go's standard formatting (`go fmt`)
- Add comments for exported functions
- Keep line length reasonable (max ~120 characters)

### Example

```go
// Good
func FetchEstimateFromAPI(ctx context.Context, coin1, coin2 string, amount float64, best bool, network1, network2 string) ([]Estimate, error) {
    // Implementation
}

// Avoid - no context, unclear naming
func api(coin1, coin2 string, amount float64) ([]Estimate, error) {
    // Implementation
}
```

---

## Testing

### Running Tests

```bash
# Run all tests with race detector
task test

# Quick test (no race detector)
task test:short
```

### Writing Tests

- Write tests for new functionality
- Use table-driven tests when appropriate
- Test edge cases and error conditions
- Aim for meaningful coverage

### Example Test

```go
func TestFetchEstimate(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    int
        wantErr bool
    }{
        {
            name:    "valid response",
            input:   mockValidResponse,
            want:    3,
            wantErr: false,
        },
        {
            name:    "invalid json",
            input:   "invalid",
            want:    0,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

---

## Building

### Build Commands

```bash
# Build for current platform
task build

# Build for all platforms
task build:all

# Build for specific platform
task build:linux
task build:macos
task build:windows
task build:macos:arm64
```

### Verification

```bash
# Verify SLSA provenance (requires Cosign)
task verify

# Verify checksums
task verify:checksum
```

---

## Submitting Changes

### 1. Commit Your Changes

```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "feat: add support for privacy coin DERO

- Add DERO to coin ID mapping
- Update price service with DERO support
- Add tests for DERO price fetching

Closes #XX"
```

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `chore`: Maintenance tasks

### 2. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 3. Create Pull Request

1. Go to https://github.com/moralpriest/cli
2. Click "New Pull Request"
3. Select your branch
4. Fill in the PR template
5. Submit

### PR Requirements

- [ ] Tests pass (`task test`)
- [ ] Linting passes (`task lint`)
- [ ] Code is properly formatted
- [ ] Documentation updated
- [ ] Clear description of changes

---

## Security

### Sensitive Data

- Never commit API keys or secrets
- Use environment variables for sensitive data
- Add `.env` to `.gitignore`

### Security Considerations for Privacy Coins

When adding support for privacy coins:

- Ensure no user data is leaked
- Follow best practices for crypto handling
- Consider supply chain security
- Document any trust assumptions

### Reporting Security Issues

For security vulnerabilities, please:

1. **Do NOT** open a public issue
2. Email: security@your-domain.com
3. Or use GitHub's private vulnerability reporting

---

## Questions

### Getting Help

- Check existing [Issues](../../issues)
- Check [Documentation](../../blob/main/README.md)
- Open a new issue with your question

### Communication

- GitHub Issues for bugs and feature requests
- Discussions for questions and ideas
- Keep discussions public when possible

---

## Thank You!

Your contributions make CypherGoat CLI better for everyone. We appreciate your time and effort!
