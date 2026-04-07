# Server Agent

## Context
This agent handles the Go backend for Corteza Neuronic AI. It manages API endpoints, authentication, and business logic. The server is built with Go 1.24.1 using chi router for HTTP handling and supports communication via REST API, web sockets, and gRPC.

**Key Architecture:**
- **Database:** PostgreSQL and MySQL with abstract store layer
- **File Storage:** MinIO for local/cloud storage
- **Email:** SMTP relay with MailHog support for local testing
- **Logging:** Structured logging (JSON/human-readable) with Sentry integration
- **Monitoring:** Primitive resource usage logging
- **Scheduler:** Interval automation execution control
- **Health Checks:** Scheduler, email, and automation service checks

### Skills

#### Core Server Components
- **Go Modules**: Backend API implemented in Go 1.24.1 using chi router
- **Authentication**: Complete OAuth2 server/client implementation supporting authorization code flow, implicit flow, and client credentials flow. Features include:
  - Local credentials (email/password) authentication
  - External OIDC providers (Google, LinkedIn, GitHub, etc.)
  - Token management (access tokens, refresh tokens, token revocation)
  - Security context and role-based access control (RBAC) with scope-based authorization
  - Auth scope system with permitted, prohibited, and forced role memberships per client
  - Security features: Brute-force protection, CSRF protection, JWT token signing, cookie encryption
  - Session management with separate authentication sessions per client and configurable expiration times
- **Compose**: Low-code platform for building applications (namespaces, modules, records)
- **Automation**: Workflow management and execution engine
- **Discovery**: Search and discovery functionality
- **Federation**: Data federation and synchronization between instances
- **System**: Core system functionality (users, roles, templates)

#### Architecture & Infrastructure
- **Database Management**: PostgreSQL/MySQL with abstract store layer and schema migrations
- **File Storage**: MinIO integration for local/cloud storage
- **Email Relay**: SMTP service with MailHog support for local testing
- **Logging**: Structured JSON/human-readable logging with Sentry integration
- **Monitoring**: Resource usage tracking and health checks
- **Scheduler**: Interval automation execution control

#### API Development
- **REST API**: Standard RESTful endpoints with Swagger/OpenAPI 3 documentation
- **gRPC**: Communication between backend services
- **WebSockets**: Real-time communication for web applications
- **API Code Generation**: Auto-generated from rest.yaml definitions using openapi3-converter

### Key Commands

#### Code Generation
```bash
# Generate all server files from CUE schemas
cd server/codegen && make server

# Generate documentation (requires DOCS_DIR to be set)
cd server/codegen && make docs DOCS_DIR=/path/to/corteza-docs

# Generate both server files and documentation
cd server/codegen && make all DOCS_DIR=/path/to/corteza-docs
```

**IMPORTANT NOTE:** Never manually edit files with `gen.go` in their name. These are automatically generated files that will be overwritten when you run the code generation process. If you need to make changes to these files, you must modify the CUE schemas in `server/codegen/schema/` or the templates in `server/codegen/assets/templates/` and then regenerate the files.

#### Build & Development
```bash
# Run tests
cd server && make test

# Run linting
cd server && make lint

# Build server
cd server && make build

# Start development server with hot reload (requires gin)
cd server && gin --laddr localhost --build cmd/corteza --bin build/gin-bin -- --env-file .env serve-api
```

#### Database
```bash
# Run database migrations
cd server && make migrate

# Reset database
cd server && make reset-db
```

### Project Structure
```
server/
├── app/               # Application context and initialization (setup levels: Setup → DB → Services → Provision → Activate)
├── cmd/               # Main entry point and CLI commands
├── auth/              # Authentication and authorization (OAuth2 server/client, local/external providers)
├── compose/           # Low-code platform for building applications (namespaces, modules, records)
├── automation/        # Automation and business process management (workflows)
├── discovery/         # Search and discovery functionality
├── federation/        # Data federation and synchronization between instances
├── store/             # Database abstraction layer (PostgreSQL/MySQL support, migrations)
├── system/            # System configuration (users, roles, templates, resource management)
├── webapp/            # Static web app serving directly from server binary
├── tests/             # Integration tests categorized by service
├── docs/              # Generated API documentation (OpenAPI 3)
├── pkg/               # Generic packages used across the system (no service dependencies)
├── provision/         # Core resources (roles, access control, built-in extensions)
└── Makefile           # Server-specific Makefile
```
