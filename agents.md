# Corteza Neuronic AI - Complete Project Agent

## Context
This is the root agent for the complete Corteza Neuronic AI project. It has full control over both the frontend and backend components, ensuring that any changes made are properly coordinated across the entire system.

### Skills

#### Project Coordination
- **Complete System Control**: Manages both frontend (Vue.js) and backend (Go) components
- **Cross-Component Communication**: Coordinates API calls, data flow, and state management between frontend and backend
- **Dependency Management**: Ensures all libraries and dependencies are properly linked and versioned
- **Build Orchestration**: Handles the entire build process from code compilation to deployment

#### Full Stack Development

##### Frontend Development (Vue.js)
- **Vue 2.x**: Use [`vue2` skill](.kilocode/skills/vue2/SKILL.md) for Vue 2.x development including Options API, components, directives, lifecycle hooks, computed properties, watchers, and Vue Router.
- **Vuex (Vue 2.x)**: Use [`vuex-vue2` skill](.kilocode/skills/vuex-vue2/SKILL.md) for Vuex 2.x state management in Vue 2 applications including state, mutations, actions, getters, modules, and plugins.

##### Backend Development (Go)
- **Golang Pro**: Use [`golang-pro` skill](.kilocode/skills/golang-pro/SKILL.md) for modern Go 1.21+ development with advanced concurrency, performance optimization, and production-ready microservices.
- **Go Concurrency Patterns**: Use [`go-concurrency-patterns` skill](.kilocode/skills/go-concurrency-patterns/SKILL.md) for goroutines, channels, sync primitives, and context management.
- **gRPC in Go**: Use [`grpc-golang` skill](.kilocode/skills/golang-grpc/SKILL.md) for building production-ready gRPC services in Go with mTLS, streaming, and observability.

##### Database & Persistence
- **Database Migrations**: Use [`database-migrations` skill](.kilocode/skills/database-migrations/SKILL.md) for database migration best practices across PostgreSQL, MySQL, and common ORMs (golang-migrate for Go).

##### Documentation & Libraries
- **Context7 Docs Lookup**: Use [`context7-docs-lookup` skill](.kilocode/skills/context7-docs-lookup/SKILL.md) to fetch up-to-date library documentation using Context7.

#### System Architecture
- **Microservices Communication**: REST API, WebSockets, and gRPC between services
- **Security**: OAuth2 authentication with local and external (OIDC) providers
- **Logging & Monitoring**: Structured logging, Sentry integration, and health checks
- **Scheduler**: Interval automation execution control

### Key Commands

#### Full System Build
```bash
# Windows
.\build.ps1                    # Build entire system
.\build.ps1 -InstallAll        # Install dependencies and build
.\build.ps1 -SkipBuild        # Skip build and just serve
.\build.ps1 -ServeOnly        # Serve without rebuilding
.\build.ps1 -CodeGen          # Run server code generation only
.\build.ps1 -Dev              # Build and run in development mode
.\build.ps1 -Test             # Run all tests (frontend + backend)
.\build.ps1 -Lint             # Run linting on all code
.\build.ps1 -Fresh            # Clean and rebuild everything
.\build.ps1 -Audit            # Run security audits

# Linux/macOS
make dev                      # Build and run in development mode
make test                     # Run all tests (frontend + backend)
make lint                     # Run linting on all code
make fresh                    # Clean and rebuild everything
make audit                    # Run security audits
```

#### Component-Specific Build
```bash
# Build only backend
cd server && make build

# Build only frontend apps
cd client && make build

# Build only shared libraries
cd lib && make build

# Start development server with hot reload (requires gin)
cd server && gin --laddr localhost --build cmd/corteza --bin build/gin-bin -- --env-file .env serve-api
```

#### Docker Deployment
```bash
# Start entire system with Docker Compose
docker-compose -f docker-compose.local.yaml up

# Build Docker image
docker build -t corteza .
```

### Project Structure
```
corteza-neuronicai/
├── server/                  # Go backend
│   ├── app/                 # Application context and initialization
│   ├── cmd/                 # Main entry point and CLI commands
│   ├── auth/                # Authentication and authorization
│   ├── compose/             # Low-code platform
│   ├── automation/          # Workflow management
│   ├── discovery/           # Search and discovery
│   ├── federation/          # Data federation
│   ├── store/               # Database abstraction layer
│   ├── system/              # System configuration
│   ├── webapp/              # Static web app serving
│   ├── tests/               # Integration tests
│   ├── docs/                # API documentation
│   ├── pkg/                 # Generic packages
│   ├── provision/           # Core resources
│   └── Makefile             # Server build orchestrator
├── client/                  # Vue.js frontend apps
│   ├── web/                 # Web applications
│   │   ├── one/             # Main application
│   │   ├── compose/         # Low-code builder
│   │   ├── admin/           # Administration
│   │   ├── workflow/        # Workflow automation
│   │   ├── reporter/        # Reporting & analytics
│   │   ├── discovery/       # Search & discovery
│   │   └── privacy/         # Data privacy
│   └── Makefile             # Client build orchestrator
├── lib/                     # Shared JS/Vue libraries
│   ├── js/                  # JavaScript/TypeScript core library
│   │   ├── src/
│   │   │   ├── api-clients/ # Auto-generated REST API clients
│   │   │   ├── automation/  # Automation helpers
│   │   │   ├── compose/     # Compose types
│   │   │   ├── corredor/    # Script parsing
│   │   │   ├── eventbus/    # Client-side event bus
│   │   │   └── system/      # Core system types
│   │   └── tools/
│   │       └── codegen/     # API client code generator
│   ├── vue/                 # Vue.js components and utilities
│   │   ├── src/
│   │   │   ├── components/  # Generic Vue components
│   │   │   ├── corredor/    # Corredor-specific logic
│   │   │   ├── libs/        # Generic libraries
│   │   │   ├── mixins/      # Reusable Vue mixins
│   │   │   ├── plugins/     # Vue plugins
│   │   │   └── store/       # Vuex store definitions
│   │   └── rollup.config.js
│   ├── eslint-client/       # ESLint configuration
│   └── Makefile             # Library build orchestrator
├── build.ps1                # Windows build script
├── Makefile                 # Main build orchestrator
└── docker-compose.local.yaml # Local Docker Compose configuration
```

### Development Workflow
1. **Modify code** in either frontend (client/), backend (server/), or shared libraries (lib/)
2. **Run tests** to verify changes: `make test`
3. **Build and serve** the entire system: `make dev`
4. **Commit changes** with appropriate commit messages

### Important Notes
- **Complete System Control**: This repository has full control over both frontend and backend
- **Consistent Build Process**: Always use the root Makefile or build.ps1 to ensure all components are properly built
- **Dependency Management**: Shared libraries are built and linked automatically for all web applications
- **API Compatibility**: Changes to the backend API should always be reflected in the frontend code and API documentation
- **Database Migrations**: Always run migrations before serving the application

### Architecture Guidelines
- **Frontend**: Vue.js with Bootstrap-Vue, Vuex, and Vue Router
- **Backend**: Go 1.24.1 with chi router, OAuth2 authentication, and abstract database layer
- **Database**: PostgreSQL/MySQL with schema migrations and transaction management
- **Security**: Role-based access control (RBAC) with scope-based authorization
- **API**: REST with OpenAPI 3 documentation, WebSockets for real-time communication, and gRPC for service-to-service communication
