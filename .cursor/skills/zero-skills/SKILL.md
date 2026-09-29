---
name: zero-skills
description: Use when building, reviewing, or troubleshooting go-zero services, including REST and RPC code generation, middleware, service discovery, collection data structures, resilience, observability, distributed transactions, and message queues. Do not use for database/model/SQL/MongoDB/Redis persistence guidance. Trigger for goctl, .api or .proto files, or projects importing github.com/zeromicro/go-zero.
version: 1.0.0
license: MIT
allowed-tools:
  - Read
  - Grep
  - Glob
---

# go-zero Skills for AI Agents

This skill provides comprehensive go-zero microservices framework knowledge, optimized for AI agents helping developers build production-ready services. It covers REST APIs, RPC services, resilience patterns, high-performance data structures, and troubleshooting. Database / model / SQL / MongoDB / Redis persistence patterns are out of scope for this installation.

## 🎯 When to Use This Skill

Invoke this skill when working with go-zero:
- **Creating services**: REST APIs, gRPC services, or microservices architectures
- **Data structures**: LRU Cache, Ring Buffer, TimingWheel, or concurrent maps
- **Production hardening**: Circuit breakers, rate limiting, or error handling
- **Debugging**: Understanding errors, fixing configuration, or resolving issues
- **Learning**: Understanding go-zero patterns and best practices

Do **not** use this skill for database schema design, `goctl model`, sqlx/Mongo CRUD, or Redis cache-aside persistence.

## 📚 Knowledge Structure

This skill organizes go-zero knowledge into focused modules. **Load specific guides as needed** rather than reading everything at once:

### Quick Start Guide
**Link**: [Official go-zero Documentation](https://go-zero.dev/docs/quick-start)
**Contains**: Installation, first API service, basic commands, hello-world examples (refer to official docs)

### Pattern Guides (Detailed Reference)

#### 1. REST API Patterns
**File**: [references/rest-api-patterns.md](references/rest-api-patterns.md)
**When to load**: Creating HTTP endpoints, implementing CRUD operations, adding middleware
**Contains**:
- Handler → Logic → Context three-layer architecture
- Request/response handling with proper types
- Middleware (auth, logging, metrics, CORS)
- Error handling with `httpx.Error()` and `httpx.OkJson()`
- Complete CRUD examples with ✅ correct vs ❌ incorrect patterns

#### 2. RPC Service Patterns
**File**: [references/rpc-patterns.md](references/rpc-patterns.md)
**When to load**: Building gRPC services, service-to-service communication
**Contains**:
- Protocol Buffers definition and code generation
- Service discovery with etcd/consul/kubernetes
- Load balancing strategies
- Client configuration and interceptors
- Error handling in RPC contexts

#### 3. Resilience Patterns
**File**: [references/resilience-patterns.md](references/resilience-patterns.md)
**When to load**: Production hardening, handling failures, managing system load
**Contains**:
- Circuit breaker configuration (Breaker)
- Rate limiting and API throttling
- Load shedding under pressure
- Timeout and retry strategies
- Graceful shutdown and degradation

#### 4. Collection Patterns
**File**: [references/collection-patterns.md](references/collection-patterns.md)
**When to load**: Using high-performance data structures, local caching, timeout management
**Contains**:
- LRU Cache for local hot data caching (O(1) operations)
- Ring Buffer for fixed-size circular buffers
- TimingWheel for efficient timer management (O(1) add/remove/execute)
- SafeMap for concurrent-safe map operations
- Memory management best practices and pitfalls

#### 5. goctl Command Reference
**File**: [references/goctl-commands.md](references/goctl-commands.md)
**When to load**: Generating code with goctl, setting up new services, post-generation steps
**Contains**:
- goctl installation and detection
- API/RPC generation commands with exact flags (no `goctl model` / database model generation)
- Post-generation pipeline (mod tidy, import fixing, build verification)
- Config templates (API, RPC, production)
- Deployment templates (Dockerfile, Kubernetes, Docker Compose)
- Middleware and error handler templates
- API spec patterns (CRUD, JWT, mixed auth)

#### 6. Swagger / OpenAPI Documentation
**File**: [references/swagger-patterns.md](references/swagger-patterns.md)
**When to load**: Generating Swagger/OpenAPI docs from `.api` files, adding API documentation annotations
**Contains**:
- Built-in `goctl api swagger` command reference (goctl ≥ 1.8.4, no plugin needed)
- `info` block metadata: title, description, host, basePath, schemes, useDefinitions, wrapCodeMsg
- `@server` tags for Swagger UI grouping
- `@doc` for endpoint summary/description
- Field tags: `example`, `options` (enum), `range`, `default`, `optional`
- Security definitions (`securityDefinitionsFromJson` + `authType`)
- Business error codes (`bizCodeEnumDescription`)
- Complete multi-file `.api` structure example with best practices

#### 7. Distributed Transaction Patterns
**File**: [references/distributed-transactions.md](references/distributed-transactions.md)
**When to load**: Cross-service data consistency, DTM integration, SAGA/TCC patterns
**Contains**:
- Pattern-selection guidance for Workflow, Saga, TCC, XA, two-phase message, and outbox
- Verified DTM HTTP Saga skeleton
- Barrier, idempotency, compensation, security, and failure-testing checklists

#### 8. Observability Patterns
**File**: [references/observability.md](references/observability.md)
**When to load**: Production monitoring, tracing, alerting setup
**Contains**:
- Prometheus metrics configuration
- Custom metrics implementation
- Distributed tracing with OpenTelemetry/Jaeger
- Structured logging with logx
- Grafana dashboards and alerting rules
- ELK integration for log aggregation

#### 9. Message Queue Patterns
**File**: [references/message-queue.md](references/message-queue.md)
**When to load**: Async processing, delayed tasks, event streaming
**Contains**:
- go-queue dq (Beanstalkd) for delayed tasks
- go-queue kq (Kafka) for high-throughput messaging
- Current producer and consumer APIs
- Lifecycle, retry, idempotency, and shutdown guidance

#### 10. Advanced Components
**File**: [references/advanced-components.md](references/advanced-components.md)
**When to load**: Performance optimization, concurrent processing, caching
**Contains**:
- Bloom filter for cache penetration prevention
- TimingWheel for delayed task scheduling
- MapReduce for parallel processing
- Executors for batch task buffering
- SharedCalls (SingleFlight) for duplicate prevention

### Supporting Resources

#### Best Practices
**File**: [best-practices/overview.md](best-practices/overview.md)
**When to load**: Production deployment, code review, optimization
**Contains**: Configuration management, logging, monitoring, security, performance

#### Troubleshooting
**File**: [troubleshooting/common-issues.md](troubleshooting/common-issues.md)
**When to load**: Debugging errors, configuration issues, runtime problems
**Contains**: Common error messages, solutions, configuration pitfalls, debugging tips

#### Claude Code Integration
**File**: [getting-started/claude-code-guide.md](getting-started/claude-code-guide.md)
**When to load**: Setting up Claude Code for zero-skills usage
**Contains**: Installation, invocation methods, advanced features (subagents, dynamic context)

#### Design Principles
**File**: [design/principles.md](design/principles.md)
**When to load**: Understanding framework philosophy, architecture decisions
**Contains**: Three-layer architecture rationale, design decisions, performance considerations, anti-patterns

#### Case Study: go-zero-looklook
**File**: [examples/case-studies/looklook-overview.md](examples/case-studies/looklook-overview.md)
**When to load**: Learning from production-scale example, real-world architecture
**Contains**: Source-grounded repository map, study workflow, reusable patterns, and verification checklist

#### Tool Integration Guides
**File**: [getting-started/README.md](getting-started/README.md)
**When to load**: Setting up zero-skills with Cursor, GitHub Copilot, Windsurf, or Codex
**Contains**: Feature comparison table, per-tool setup instructions (Claude Code, Cursor, Copilot, Windsurf, Codex)

## 🚀 Common Workflows

These workflows guide you through typical go-zero development tasks:

### Creating a New REST API Service

**Steps:**
1. Define API specification in `.api` file with types and routes
2. Generate code: `goctl api go -api user.api -dir .`
3. Implement business logic in `internal/logic/` layer
4. Add validation and error handling with `httpx`
5. Test endpoints with proper request/response handling

**Detailed guide**: [references/rest-api-patterns.md](references/rest-api-patterns.md#complete-rest-api-workflow)

### Adding Middleware

**Steps:**
1. Create middleware function in `internal/middleware/` directory
2. Define middleware in `.api` file or register programmatically
3. Implement authentication/authorization logic
4. Pass validated data through `r.Context()`
5. Handle errors with appropriate HTTP status codes

**Detailed guide**: [references/rest-api-patterns.md](references/rest-api-patterns.md#middleware-patterns)

### Building an RPC Service

**Steps:**
1. Define service in `.proto` file with messages and RPCs
2. Generate code: `goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.`
3. Implement service logic in `internal/logic/`
4. Configure service discovery (etcd/consul/kubernetes)
5. Test with RPC client and handle errors

**Detailed guide**: [references/rpc-patterns.md](references/rpc-patterns.md#complete-rpc-workflow)

### Generating Swagger Documentation

**Steps:**
1. Add `info` block with swagger metadata (title, host, basePath, schemes) to entry `.api` file
2. Add `tags` in `@server` blocks for Swagger UI grouping
3. Add `@doc` summary to each endpoint
4. Add `example`/`options`/`range` tags to request/response fields
5. Run `goctl api swagger --api entry.api --dir docs/swagger --filename api`

**Detailed guide**: [references/swagger-patterns.md](references/swagger-patterns.md)

### Using High-Performance Data Structures

**Steps:**
1. Choose appropriate data structure based on use case:
   - LRU Cache for local hot data caching
   - Ring Buffer for fixed-size log/message buffers
   - TimingWheel for timeout/delayed task management
   - SafeMap for concurrent-safe key-value storage
2. Import from `github.com/zeromicro/go-zero/core/collection`
3. Configure capacity/parameters based on memory constraints
4. Handle eviction callbacks for resource cleanup (LRU Cache)
5. Monitor memory usage and rebuild when necessary (SafeMap)

**Detailed guide**: [references/collection-patterns.md](references/collection-patterns.md)

## ⚡ Key Principles

When generating or reviewing go-zero code, always apply these principles:

### ✅ Always Follow

- **Three-layer separation**: Keep Handler (routing) → Logic (business) → Model (data) distinct
- **Structured errors**: Use `httpx.Error(w, err)` for HTTP errors, not `fmt.Errorf`
- **Configuration**: Load with `conf.MustLoad(&c, *configFile)` and inject via ServiceContext
- **Context propagation**: Pass `ctx context.Context` through all layers for tracing and cancellation
- **Type safety**: Define request/response types in `.api` files, generate with goctl
- **goctl generation**: Always use `goctl` to generate boilerplate, never hand-write handlers/routes

### ❌ Never Do

- Put business logic directly in handlers (violates three-layer architecture)
- Return raw errors with `w.Write()` or `fmt.Fprintf()` instead of using httpx helpers
- Hard-code configuration values (ports, hosts, database credentials)
- Skip validation of user inputs or forget to check `err != nil`
- Modify generated code (customize via `logic` layer instead)
- Bypass ServiceContext injection (leads to tight coupling and testing issues)

## 📖 Progressive Learning Path

Follow this path based on your needs:

### 🟢 New to go-zero?

1. **Start here**: [Official go-zero Quick Start](https://go-zero.dev/docs/quick-start)
   Install go-zero, create your first API, understand basic concepts

2. **Build REST/RPC**: [references/rest-api-patterns.md](references/rest-api-patterns.md) and [references/rpc-patterns.md](references/rpc-patterns.md)
   Three-layer architecture, middleware, service discovery

### 🟡 Building production services?

1. **Review best practices**: [best-practices/overview.md](best-practices/overview.md)
   Configuration, logging, monitoring, security checklist

2. **Add resilience**: [references/resilience-patterns.md](references/resilience-patterns.md)
   Circuit breakers, rate limiting, graceful degradation

3. **Optimize data structures**: [references/collection-patterns.md](references/collection-patterns.md)
   LRU Cache, Ring Buffer, TimingWheel for performance

4. **Check common pitfalls**: [troubleshooting/common-issues.md](troubleshooting/common-issues.md)
   Avoid typical mistakes and know how to debug issues

5. **Set up observability**: [references/observability.md](references/observability.md)
   Prometheus metrics, distributed tracing, structured logging

### 🔴 Advanced scenarios?

1. **Distributed transactions**: [references/distributed-transactions.md](references/distributed-transactions.md)
   DTM integration, SAGA/TCC patterns, data consistency

2. **Message queues**: [references/message-queue.md](references/message-queue.md)
   go-queue for async processing, delayed tasks, Kafka integration

3. **Performance optimization**: [references/advanced-components.md](references/advanced-components.md)
   Bloom filter, MapReduce, TimingWheel, SharedCalls

4. **Understand design**: [design/principles.md](design/principles.md)
   Framework philosophy, architecture decisions, anti-patterns

5. **Learn from examples**: [examples/case-studies/looklook-overview.md](examples/case-studies/looklook-overview.md)
   Production-scale architecture, real-world patterns

### 🔵 Extending capabilities?

1. **Use with Claude Code**: [getting-started/claude-code-guide.md](getting-started/claude-code-guide.md)
   Learn advanced features like subagents, dynamic context, and argument passing
   Run demo projects to validate your environment

2. **Verify knowledge**: [examples/verify-tutorial.sh](examples/verify-tutorial.sh)
   Script to check if examples work correctly

## 🔗 Integration with go-zero AI Ecosystem

This skill is part of a two-layer ecosystem for AI-assisted go-zero development:

| Tool | Purpose | Best For |
|------|---------|----------|
| **[ai-context](https://github.com/zeromicro/ai-context)** | Concise workflow instructions (~5KB) | GitHub Copilot, Cursor, Windsurf |
| **zero-skills** (this repo) | Comprehensive knowledge base + goctl reference (~45KB) | All AI tools, deep learning, reference |

The AI runs `goctl` directly in the terminal for code generation — no separate MCP server needed. See [references/goctl-commands.md](references/goctl-commands.md) for the complete command reference.

**Usage in Claude Code:**
- This skill loads automatically when working with go-zero projects
- Use `/zero-skills` to invoke manually for go-zero guidance
- AI runs goctl commands directly in the terminal for code generation
- Reference specific pattern files when needed (Claude loads them on demand)

See [getting-started/claude-code-guide.md](getting-started/claude-code-guide.md) for detailed usage instructions.

## 🌐 Additional Resources

- **Official docs**: [go-zero.dev](https://go-zero.dev) - Latest API reference and guides
- **GitHub**: [zeromicro/go-zero](https://github.com/zeromicro/go-zero) - Source code and examples
- **Community**: Discussions, issues, and contributions welcome in the main repository

## 📝 Version Compatibility

- **Stable target**: go-zero v1.10.3 (Go 1.24 or later required)
- **Master compatibility**: Use Go 1.25 or later when testing against the go-zero `master` branch
- **Updates**: Patterns updated regularly to reflect framework evolution
- **Breaking changes**: Check official docs for API changes between versions

---

**Quick invocation**: Use `/zero-skills` or ask "How do I [task] with go-zero?"
**Need help?** Reference the specific pattern guide for detailed examples and explanations.
