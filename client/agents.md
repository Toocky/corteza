# Client Agent

## Context
This agent handles the frontend applications for Corteza Neuronic AI. It manages multiple Vue.js web apps that provide the user interface for the platform.

### Skills

#### Web Applications
- **One**: Main web application
- **Admin**: Administration interface
- **Compose**: Low-code application builder
- **Workflow**: Workflow automation interface
- **Reporter**: Reporting and analytics
- **Discovery**: Search and discovery
- **Privacy**: Data privacy and GDPR compliance

#### Frontend Development
- **Vue.js 2.7**: Progressive JavaScript framework
- **Bootstrap-Vue**: UI component library
- **Vuex**: State management
- **Vue Router**: Client-side routing
- **Webpack**: Module bundler
- **Babel**: JavaScript transpiler
- **C3 Component Catalog**: Tool for building UI components in an insulated environment (similar to Storybook)
- **Internationalization**: i18next for multi-language support

#### Authentication & Authorization
- OAuth2 flow with Corteza Auth server
- Local and external (OIDC) identity providers
- Token management and refresh
- Security context and scope-based authorization

### Key Commands

#### Build
```bash
# Build all web applications
cd client && make build

# Install dependencies and link libraries
cd client && make dev

# Clean and rebuild applications
cd client && make fresh
```

#### Development
```bash
# Run all tests
cd client && make test

# Run linting
cd client && make lint
```

#### Audit
```bash
# Run dependency audits
cd client && make audit
```

#### Individual App Commands
```bash
# Build a specific app (replace {app} with one/admin/compose/workflow/reporter/discovery/privacy)
cd client/web/{app} && yarn build

# Serve a specific app (replace {app} with one/admin/compose/workflow/reporter/discovery/privacy)
cd client/web/{app} && yarn serve

# Run tests for a specific app
cd client/web/{app} && yarn test:unit

# Link libraries for a specific app
cd client/web/{app} && yarn cdeps
```

#### Component Development with C3
```bash
# Start C3 component catalog (development mode only)
cd client/web/{app} && yarn dev-env
```

**C3 Component Catalog:**
- **Purpose**: Tool for building UI components in an insulated environment, making it easier to develop independent components
- **Inspiration**: Inspired by Storybook but with "hooking into" components instead of defining stories to outline functionality
- **Configuration Files**:
  - `src/dev-env.js`: Defines the insulated Vue application configuration for C3
  - `src/components/C3.js`: Main component registry that defines which components are available in C3
  - `*.c3.js`: Component-specific configuration files defining props, controls, and scenarios
- **Key Features**:
  - **Props Definition**: Define initial component properties
  - **Controls**: Interactive controls for modifying component properties (checkbox, input, textarea, select, generic)
  - **Scenarios**: Predefined states to showcase component behavior with different property combinations
  - **Grouping**: Components can be organized into groups for better navigation
- **Usage Flow**:
  1. Create a new Vue component (e.g., `MyComponent.vue`)
  2. Create a corresponding configuration file (e.g., `MyComponent.c3.js`)
  3. Register the component in the C3 registry (`src/components/C3.js`)
  4. Start the C3 dev environment with `yarn dev-env`
  5. View and interact with the component in the C3 UI

### Project Structure
```
client/
├── web/              # Web applications
│   ├── one/          # Main application
│   │   ├── public/   # Static assets (config, logo)
│   │   ├── src/      # Application source
│   │   │   ├── components/    # Vue components
│   │   │   ├── i18n/         # Internationalization assets
│   │   │   ├── mixins/       # Vue mixins
│   │   │   ├── plugins/      # Vue plugins
│   │   │   ├── store/        # Vuex store
│   │   │   ├── themes/       # Styling assets
│   │   │   └── views/        # Application views & router
│   │   └── package.json
│   ├── compose/      # Low-code application builder
│   ├── admin/        # Administration interface
│   ├── workflow/     # Workflow automation
│   ├── reporter/     # Reporting & analytics
│   ├── discovery/    # Search & discovery
│   └── privacy/      # Data privacy & GDPR compliance
└── Makefile          # Client build orchestrator
```

### Application Dependencies
- **@cortezaproject/corteza-js**: Core API library
- **@cortezaproject/corteza-vue**: Vue components
- **bootstrap-vue**: Bootstrap 4 components
- **echarts**: Charting library
- **@fullcalendar/***: Calendar components
- **@fortawesome/***: Icon library
- **vue-echarts**: Vue wrapper for ECharts
- **vue-grid-layout**: Grid layout system
- **vuedraggable**: Drag and drop
- **portal-vue**: Portal components

### Build Target
All built applications are deployed to `server/webapp/public/` directory for serving.
