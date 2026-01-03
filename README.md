# GoCart API: Engineering a Production-Grade E-Commerce Nexus

Welcome to **GoCart API**. This isn't just another CRUD-heavy e-commerce skeleton; it's a technical blueprint for building resilient, high-concurrency distributed systems using the Go ecosystem. As a Engineer, I've designed this system to showcase the intersection of Clean Architecture, Event-Driven reliability, and hybrid API design.

---

## 🏛 The Architectural Blueprint: Beyond Layered Design

At its core, GoCart API implements a strictly decoupled **Clean Architecture**. This isn't about dogmatic adherence to rules, but about maximizing **testability** and **adaptability**.

- **Handlers Layer**: Pure interface with the outside world (REST via Gin, GraphQL via gqlgen).
- **Service Layer**: The "Brain." Orchestrates domain logic, transaction boundaries, and event emission. This layer remains agnostic of whether the data came from a JSON POST or a GraphQL Mutation.
- **Repository Layer**: The "Hands." Direct interaction with GORM/PostgreSQL, ensuring the service layer never leaks SQL-specific logic.

### Why this matters?
This abstraction allows me to rotate the persistence layer or swap out the event bus (e.g., from LocalStack SQS to production AWS SNS) without touching a single line of business logic.

---

## 🛠 Tech Stack & Tooling: My Arsenal

A project is only as strong as its foundations. I've hand-picked a suite of high-performance libraries and tools:

### Core Frameworks
- **[Go (1.24.4+)](https://go.dev/)**: The backbone of this high-concurrency service.
- **[Gin Gonic](https://github.com/gin-gonic/gin)**: A high-performance HTTP web framework.
- **[GORM](https://gorm.io/)**: The fantastic ORM library for Go, providing a clean API for PostgreSQL interactions.

### API & Communication
- **[gqlgen](https://github.com/99designs/gqlgen)**: Type-safe GraphQL code generation.
- **[Watermill](https://github.com/ThreeDotsLabs/watermill)**: A Go library for efficiently handling message streams (Pub/Sub).
- **[Swag](https://github.com/swaggo/swag)**: Automated OpenAPI 2.0/Swagger documentation generation.

### Infrastructure & Cloud
- **[PostgreSQL](https://www.postgresql.org/)**: Primary relational data store.
- **[AWS SDK v2 for Go](https://github.com/aws/aws-sdk-go-v2)**: For robust interaction with S3 and SQS.
- **[LocalStack](https://localstack.cloud/)**: A fully functional local AWS cloud stack for local development and testing.
- **[Docker & Compose](https://www.docker.com/)**: For reproducible environment orchestration.

### Security & Observability
- **[JWT (golang-jwt)](https://github.com/golang-jwt/jwt)**: For stateless, secure authentication.
- **[Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)**: Industry-standard password hashing.
- **[Zerolog](https://github.com/rs/zerolog)**: Zero-allocation, JSON-structured logging.

### Engineering & Testing
- **[Testify](https://github.com/stretchr/testify)**: A toolkit for assertions and mock objects.
- **[GoMock](https://github.com/golang/mock)**: For generating mock implementations of Go interfaces.
- **[golangci-lint](https://golangci-lint.run/)**: A fast Go linters runner.

---

## 🏎 Hybrid API Ecosystem: REST + GraphQL

Modern frontends (React, Mobile, SSR) have varying data requirements. I provide a dual-stack interface:

1.  **RESTful API**: Standardized, cache-friendly endpoints for traditional resource management.
2.  **GraphQL (gqlgen)**: For high-efficiency graph traversal. 
    - **Type-Safe Generation**: Utilizing Go's code generation to ensure resolvers and models are always in sync with the `.graphqls` schemas.
    - **Identity-Aware Resolvers**: Using Go's `context.Context` to propagate JWT identity deep into the graph, enabling fine-grained field-level authorization.

---

## 📬 Event-Driven Reliability with Watermill

In a real-world commerce system, consistency is king. GoCart uses the **Watermill** library to implement a robust event-driven backbone.

- **Pub/Sub Abstraction**: I've decoupled the service logic from the transport. In local dev, I emulate AWS SQS via **LocalStack**.
- **Asynchronous Workflows**: Critical side effects (notifications, inventory sync, analytics) are processed out-of-band, reducing the critical path of user requests and drastically improving latency.

---

## 🛡 Security Infrastructure: Hardened JWT & RGBAC

Security isn't a bolt-on; it's a foundation.

- **JWT Rotation**: Implements a short-lived Access Token + long-lived Refresh Token strategy with database-backed revocation.
- **RGBAC (Role-Group Based Access Control)**: Middleware-driven authorization that checks roles (Admin/Customer) before the request even hits the service layer.
- **Bcrypt Hashing**: Industry-standard password salting and hashing ensures user data remains secure even in the event of a partial DB leak.

---

## 🛠 Engineering DX & Environment Parity

I follow a "Local-First" philosophy to eliminate the "works on my machine" anti-pattern.

- **Containerized Dependencies**: PostgreSQL and LocalStack are orchestrated via `docker-compose`.
- **Dependency Injection (DI)**: Using Go interfaces to swap between `LocalUploadProvider` and `S3Provider` based on environment configuration.
- **Linting & Quality**: A strict `golangci-lint` configuration ensures code consistency across contributors.

---

## 🚦 Getting Started

I've detailed a complete [Local Setup Guide](file:///Users/kuldeepsingh/Desktop/Golang/gocart-api/LocalSetup.md) that covers:
- Dependency orchestration via Docker.
- Environment variable injection.
- Testing the REST, GraphQL, and Documentation (Swagger/RapiDoc) layers.

---

## 🧪 Testing Strategy: The Staff Engineer's Bar

I don't just "hit the endpoints." My testing suite includes:
- **Unit Tests**: Mocking dependencies with `Gomock` to test services in isolation.
- **GraphQL Integration Tests**: Automated resolver tests ensuring my schema-to-logic mapping is flawless.
- **Postman Collection**: A pre-baked collection for rapid manual verification of edge cases.

---

### Final Word
GoCart API is built for the engineer who wants to see how a "toy" project scales into a production-grade system. Dive into the code, check the `internal/services`, and see how I handle concurrency, contexts, and clean abstractions.

**Happy Engineering. 🛒🚀**
