# Testing Guide

This directory contains test files for the Go API template.

## Structure

```
tests/
├── unit/           # Unit tests for individual components
├── integration/    # Integration tests that test multiple components together
└── README.md       # This file
```

## Running Tests

### Run all tests
```bash
go test ./tests/...
```

### Run unit tests only
```bash
go test ./tests/unit/...
```

### Run integration tests only
```bash
go test ./tests/integration/...
```

### Run tests with verbose output
```bash
go test -v ./tests/...
```

### Run tests with coverage
```bash
go test -cover ./tests/...
```

### Generate coverage report
```bash
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html
```

### Run specific test
```bash
go test -run TestAuthService_EdgeCases ./tests/unit/...
```

### Run tests with race detection
```bash
go test -race ./tests/...
```

## Test Structure

### Unit Tests
- Test individual functions and methods in isolation
- Use mocks for dependencies
- Fast execution
- Located in `tests/unit/`

### Integration Tests
- Test multiple components working together
- May use real dependencies (database, HTTP server)
- Slower execution but more realistic
- Located in `tests/integration/`

## Writing Tests

### Unit Test Example
```go
func TestMyFunction(t *testing.T) {
    // Given
    input := "test input"
    expected := "expected output"
    
    // When
    result := MyFunction(input)
    
    // Then
    assert.Equal(t, expected, result)
}
```

### Integration Test Example with Test Suite
```go
type MyTestSuite struct {
    suite.Suite
    server *httptest.Server
}

func (suite *MyTestSuite) SetupSuite() {
    // Setup code that runs once before all tests
}

func (suite *MyTestSuite) TearDownSuite() {
    // Cleanup code that runs once after all tests
}

func (suite *MyTestSuite) TestSomething() {
    // Test implementation
}

func TestMyTestSuite(t *testing.T) {
    suite.Run(t, new(MyTestSuite))
}
```

## Test Configuration

Tests use a separate configuration to avoid affecting the main application:
- Use in-memory databases when possible
- Use different ports for test servers
- Mock external dependencies

## Best Practices

1. **AAA Pattern**: Arrange, Act, Assert
2. **Descriptive Names**: Test names should describe what they test
3. **One Assertion Per Test**: Focus on testing one thing at a time
4. **Clean Up**: Always clean up resources in teardown methods
5. **Independent Tests**: Tests should not depend on each other
6. **Fast Tests**: Unit tests should run quickly
7. **Use Table-Driven Tests**: For testing multiple scenarios

## Dependencies

The tests use the following libraries:
- `testify/assert`: Assertions
- `testify/suite`: Test suites
- `testify/mock`: Mocking (when needed)
- `httptest`: HTTP testing utilities