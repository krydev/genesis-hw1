# Test project for Genesis school. Web API with JWT authentication and file-based storage written in Golang

## Used 3rd-party packages:
- Gin HTTP web framework for building API
- jwt-go for building token-based authentication
- GJson for accessing nested json fields to avoid creating a struct

## Points for improvement
- Adding refresh token logic to improve user experience
- Adding logout logic for destroying tokens manually without waiting for expiration
