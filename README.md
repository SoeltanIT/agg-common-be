# Agg Common Backend

[![Go Reference](https://pkg.go.dev/badge/github.com/SoeltanIT/agg-common-be.svg)](https://pkg.go.dev/github.com/SoeltanIT/agg-common-be)

A comprehensive Go package providing common utilities and helpers for building microservices in the Dino Aggregator ecosystem. This package includes standardized response handling, error management, request validation, and pagination support.

## Installation

```bash
go get -u github.com/SoeltanIT/agg-common-be
```

## Features

- **Response Handling**: Standardized JSON response format for APIs
- **Error Management**: Custom error types with HTTP status codes and error codes
- **Request Validation**: Integration with go-playground/validator with custom translations
- **Pagination**: Built-in support for paginated API responses
- **Environment Variables**: Access to environment variables with fallback values
- **User Context**: Extracting user information from JWT tokens stored in context

## Examples
- [Custom Error](https://github.com/SoeltanIT/agg-common-be/blob/main/_examples/custom-error/main.go)
- [Response Format](https://github.com/SoeltanIT/agg-common-be/blob/main/_examples/response/main.go)
- [Validator](https://github.com/SoeltanIT/agg-common-be/blob/main/_examples/validator/main.go)
- [Environment Variables](https://github.com/SoeltanIT/agg-common-be/blob/main/_examples/env/main.go)
- [User Context](https://github.com/SoeltanIT/agg-common-be/blob/main/_examples/user-context/main.go)

## Response Format

### Success Response

```json
{
  "status": "success",
  "message": "Operation completed successfully",
  "data": {
    "id": 1,
    "name": "John Doe"
  },
  "pagination": {
    "total": 100,
    "page": 2,
    "next": "https//example.com/api/users?page=3&pageSize=10",
    "prev": "https//example.com/api/users?page=1&pageSize=10"
  }
}
```

### Error Response

```json
{
  "status": "failed",
  "code": 4002000,
  "message": "Validation failed",
  "errors": [
    "Field 'email' is required",
    "Field 'username' must be at least 3 characters"
  ]
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Fiber v2](https://gofiber.io/) - Web framework
- [Validator v10](https://github.com/go-playground/validator) - Struct validation
- [JWT v5](https://github.com/golang-jwt/jwt) - JWT implementation