# Go API Boilerplate

A production-ready Golang API boilerplate following clean architecture principles.

## Features

- ✨ Clean Architecture
- 🚀 RESTful API with versioning
- 💉 Dependency Injection
- 🔒 JWT Authentication
- 📝 Request Validation
- 📚 Swagger/OpenAPI Documentation
- 🗄️ SQLite Database with GORM
- 🔄 Database Migrations
- 📊 Prometheus Metrics
- 📈 Grafana Dashboards
- 🔍 Structured Logging
- ⚡ Rate Limiting
- 💾 Caching Layer
- 🏥 Health Check Endpoints
- 🔐 Security Best Practices
- 🧪 Unit and Integration Testing
- 🐳 Docker and Docker Compose
- 🔄 CI/CD Pipeline (GitHub Actions)
- 🎯 Repository Pattern
- 🛠️ Service Layer
- 🌐 CORS Support
- 🛑 Graceful Shutdown

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional)

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/mrfansi/go-api-boilerplate.git
cd go-api-boilerplate
```

2. Copy the environment file:

```bash
cp .env.example .env
```

3. Install dependencies:

```bash
go mod download
```

4. Run the application:

### Using Go:

```bash
go run cmd/api/main.go
```

### Using Docker Compose:

```bash
docker-compose up -d
```

## Project Structure

```
.
├── cmd/                    # Application entry points
│   └── api/               # API server
├── internal/              # Private application code
│   ├── domain/           # Enterprise business rules
│   │   ├── entity/       # Business objects
│   │   ├── repository/   # Repository interfaces
│   │   └── errors/       # Domain errors
│   ├── application/      # Application business rules
│   │   └── service/      # Use cases implementation
│   ├── infrastructure/   # External implementations
│   │   ├── config/       # Configuration
│   │   ├── container/    # Dependency injection
│   │   ├── database/     # Database implementations
│   │   └── repository/   # Repository implementations
│   └── interfaces/       # Interface adapters
│       └── http/         # HTTP layer
│           ├── handler/  # HTTP handlers
│           ├── middleware/ # HTTP middleware
│           └── router/   # HTTP router
├── pkg/                  # Public library code
├── scripts/             # Build/deploy scripts
└── test/               # Test files
```

## API Documentation

Once the application is running, you can access the Swagger documentation at:

```
http://localhost:8080/swagger/
```

## Monitoring

### Prometheus

Access Prometheus metrics and dashboard:

```
http://localhost:9090
```

### Grafana

Access Grafana dashboard (default credentials: admin/admin):

```
http://localhost:3000
```

## Testing

Run all tests:

```bash
go test -v ./...
```

Run tests with coverage:

```bash
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
