# Contributing to Go-Sieve

Thank you for considering contributing to the Go-Sieve cache implementation! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you are expected to uphold our Code of Conduct. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

- Before creating a bug report, check the issue tracker to see if the problem has already been reported
- When creating a bug report, include a clear title and description, along with as much relevant information as possible
- If possible, include steps to reproduce, expected behavior, and actual behavior

### Suggesting Enhancements

- Before creating an enhancement suggestion, check the issue tracker to see if it has already been suggested
- Provide a clear description of the enhancement, along with any specific implementation details you can offer
- Explain why this enhancement would be useful to most Go-Sieve users

### Pull Requests

- Fill in the required template
- Do not include issue numbers in the PR title
- Include screenshots and animated GIFs in your pull request whenever possible
- End all files with a newline
- Avoid platform-dependent code
- Make sure all tests pass
- Document new code based on the existing Go documentation style

## Style Guidelines

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests after the first line
- Consider starting the commit message with an applicable emoji:
  - ✨ (`:sparkles:`) when adding a new feature
  - 🐛 (`:bug:`) when fixing a bug
  - 📚 (`:books:`) when adding or updating documentation
  - 🧪 (`:test_tube:`) when adding tests
  - 🔧 (`:wrench:`) when dealing with the build system
  - ⬆️ (`:arrow_up:`) when upgrading dependencies
  - ⬇️ (`:arrow_down:`) when downgrading dependencies

### Go Style

- Follow the standard Go style guidelines
- Run `go fmt` on your code before submitting
- Use meaningful variable names
- Document all exported functions, types, and constants
- Write comprehensive tests for new functionality

## Development Process

1. Fork the repository
2. Create a new branch for your feature or bugfix (`git checkout -b feature/my-new-feature`)
3. Make your changes
4. Run tests to ensure they pass (`go test ./...`)
5. Commit your changes (`git commit -am 'Add some feature'`)
6. Push to the branch (`git push origin feature/my-new-feature`)
7. Create a new Pull Request

## Testing

- Write unit tests for all new functionality
- Make sure all existing tests pass before submitting a pull request
- Aim for high test coverage, especially for critical parts of the codebase

## Questions?

If you have any questions, feel free to open an issue with the "question" label or reach out to the maintainers directly.

Thank you for your contributions!