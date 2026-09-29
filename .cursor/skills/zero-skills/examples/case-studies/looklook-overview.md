# Studying go-zero-looklook

## Purpose

[go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook) is a community-maintained, full-stack go-zero example. Use it to inspect how a larger repository organizes API, RPC, messaging, deployment, and observability concerns. Treat its source and documentation as a case study rather than as a canonical API reference; verify framework APIs against the go-zero version selected by your project.

## Repository Map

The current repository is organized around these top-level areas:

| Path | What to inspect |
|---|---|
| `app/` | Business services, including API, RPC, and asynchronous workloads |
| `pkg/` | Shared middleware, errors, context helpers, and utilities |
| `deploy/` | Infrastructure and deployment configuration |
| `data/` | Local data created by dependent middleware |
| `doc/` | Project-specific development and deployment documentation |
| `docker-compose.yml` | Local service stack |
| `docker-compose-env.yml` | Supporting middleware for local development |

Paths and service names can change. Inspect the checked-out revision instead of copying a directory tree from a blog post or older branch.

## What the Project Demonstrates

The project README describes a stack that includes go-zero, MySQL, Redis, Kafka, go-queue, asynq, Prometheus, Grafana, Jaeger, Filebeat, go-stash, Elasticsearch, Docker, Kubernetes, and CI/CD tooling.

Useful study topics include:

- API services aggregating RPC calls.
- Direct RPC endpoints for local development and discovery-oriented deployment for other environments.
- JWT authentication and request-scoped user identity.
- Kafka publish/subscribe with go-queue.
- Redis-backed asynchronous and scheduled work with asynq.
- Structured logs shipped through Filebeat and Kafka to go-stash and Elasticsearch.
- Prometheus metrics, Grafana dashboards, and distributed tracing.
- Docker Compose for local dependencies and Kubernetes-oriented deployment material.

The upstream README explicitly notes that development favors direct connections to reduce local service-discovery setup. Do not assume every environment uses etcd or that a copied RPC configuration matches the current branch.

## A Reliable Study Workflow

### 1. Pin the revision

```bash
git clone https://github.com/Mikaelemmmm/go-zero-looklook.git
cd go-zero-looklook
git rev-parse HEAD
```

Record the commit hash when documenting a pattern. This makes path and behavior claims reproducible.

### 2. Read the project documentation

Start with the root README and `doc/`. Identify the documented startup flow and the middleware required for the feature you want to study.

### 3. Trace one request end to end

For a selected feature:

1. Find its `.api` definition and generated route.
2. Follow the handler into the logic package.
3. Identify RPC clients and model calls in the service context.
4. Follow the corresponding `.proto` definition and RPC logic.
5. Inspect error conversion, context values, and logging along the path.
6. Compare configuration values with the local and deployment manifests.

This produces a more accurate understanding than copying an isolated handler or configuration fragment.

### 4. Trace one asynchronous flow

Search `app/` for `kq.NewPusher`, `kq.MustNewQueue`, and asynq task registration. Check:

- the event or task schema;
- the producer's message key;
- consumer idempotency;
- retry and dead-letter behavior;
- shutdown and flush behavior;
- how trace and request identifiers cross the queue boundary.

### 5. Compare operational configuration

Inspect `deploy/`, `docker-compose.yml`, and `docker-compose-env.yml`. Confirm ports, endpoints, credentials, and telemetry exporters against the application configuration at the same revision.

## Patterns Worth Reusing

- Keep generated transport code thin and business behavior in logic packages.
- Construct dependencies in the service context instead of creating clients inside handlers.
- Centralize consistent error and response handling, while preserving HTTP and RPC semantics.
- Treat message handlers as idempotent because queue delivery can be repeated.
- Keep secrets and environment-specific endpoints outside source code.
- Pass `context.Context` through RPC, model, and outbound calls.

## Patterns to Evaluate Before Reusing

- Shared packages can reduce duplication but can also couple otherwise independent services.
- A single repository and Docker Compose stack simplify learning but do not define the correct production boundary.
- Custom response envelopes and error codes must match your public API contract.
- Development connection strategy, discovery, and production deployment should be chosen independently.
- Telemetry configuration evolves; validate exporter names and endpoints against the go-zero version in `go.mod`.

## Verification Checklist

Before turning a project example into guidance, confirm all of the following:

- The cited path exists at the pinned commit.
- The code is copied exactly or clearly labeled as adapted pseudocode.
- The example uses the dependency versions in that commit's `go.mod`.
- Configuration keys exist in the corresponding Go config structs.
- The startup command matches the repository's current documentation.
- No sample credentials are presented as production defaults.

## References

- [go-zero-looklook repository](https://github.com/Mikaelemmmm/go-zero-looklook)
- [go-zero-looklook documentation](https://github.com/Mikaelemmmm/go-zero-looklook/tree/main/doc)
- [go-zero examples](https://github.com/zeromicro/go-zero/tree/master/example)
- [go-zero documentation](https://go-zero.dev/docs/)
