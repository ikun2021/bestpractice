# Common Issues and Solutions

## Installation Issues

### Issue: goctl command not found

**Symptoms:**
```bash
$ goctl --version
zsh: command not found: goctl
```

**Solution:**
```bash
# Install goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# Ensure $GOPATH/bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Verify installation
goctl --version
```

### Issue: go-zero version mismatch

**Symptoms:**
```
undefined: rest.RestConf
```

**Solution:**
```bash
# Update go-zero to latest
go get -u github.com/zeromicro/go-zero

# Clean module cache
go clean -modcache

# Tidy dependencies
go mod tidy
```

## Code Generation Issues

### Issue: goctl api generate fails

**Symptoms:**
```bash
$ goctl api go -api user.api -dir .
Error: invalid syntax at line 10
```

**Solution:**

Check API syntax - common mistakes:

```go
// ❌ Wrong: Missing syntax declaration
type Request {
    Name string `json:"name"`
}

// ✅ Correct: Include syntax
syntax = "v1"

type Request {
    Name string `json:"name"`
}

// ❌ Wrong: Using 'any' type
type Request {
    Data any `json:"data"`
}

// ✅ Correct: Use concrete types
type Request {
    Data map[string]string `json:"data"`
}

// ❌ Wrong: Missing return type
@handler GetUser
get /users/:id (GetUserRequest)

// ✅ Correct: Include returns
@handler GetUser
get /users/:id (GetUserRequest) returns (GetUserResponse)
```

### Issue: goctl rpc generate fails

**Symptoms:**
```bash
$ goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
Error: protoc-gen-go: program not found
```

**Solution:**
```bash
# Install required tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Ensure they're in PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Regenerate
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

## Runtime Issues

### Issue: Service won't start - port already in use

**Symptoms:**
```
Error: listen tcp :8888: bind: address already in use
```

**Solution:**
```bash
# Find process using the port
lsof -i :8888

# Kill the process
kill -9 <PID>

# Or change port in config
# etc/config.yaml
Port: 8889  # Use different port
```

## API Issues

### Issue: 404 Not Found for valid endpoint

**Symptoms:**
```bash
$ curl http://localhost:8888/api/users
404 page not found
```

**Solution:**

Check route registration:

```go
// API definition
@server(
    prefix: /api/v1  // ✅ Check prefix
    group: user
)
service user-api {
    @handler GetUsers
    get /users (GetUsersRequest) returns (GetUsersResponse)
}

// Actual URL will be: /api/v1/users (not /api/users)
```

Verify routes are registered:
```go
// In main function
func main() {
    // Enable route logging
    logx.DisableStat()

    // Routes will be logged on startup
    server.Start()
}
```

### Issue: Request body not parsed

**Symptoms:**
```go
// All request fields are empty/zero
req.Name == ""
req.Age == 0
```

**Solution:**

Check Content-Type header:

```bash
# ❌ Wrong: Missing Content-Type
curl -X POST http://localhost:8888/api/users \
  -d '{"name":"John"}'

# ✅ Correct: Include Content-Type
curl -X POST http://localhost:8888/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John"}'
```

Check JSON tags:

```go
// ❌ Wrong: Missing json tags
type Request struct {
    Name string
    Age  int
}

// ✅ Correct: Include json tags
type Request struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

### Issue: Path parameter not parsed

**Symptoms:**
```go
// Path parameter is 0 or empty
req.Id == 0
```

**Solution:**

Check tag in type definition:

```go
// ❌ Wrong: Using json tag for path parameter
type GetUserRequest struct {
    Id int64 `json:"id"`
}

// ✅ Correct: Use path tag
type GetUserRequest struct {
    Id int64 `path:"id"`
}

// API definition
@handler GetUser
get /users/:id (GetUserRequest) returns (GetUserResponse)
```

## RPC Issues

### Issue: RPC service not discovered

**Symptoms:**
```
Error: rpc error: code = Unavailable desc = connection error
```

**Solution:**

Check etcd configuration:

```yaml
# Server configuration (etc/user-rpc.yaml)
Name: user.rpc
ListenOn: 0.0.0.0:8080
Etcd:
  Hosts:
    - 127.0.0.1:2379  # ✅ Ensure etcd is running
  Key: user.rpc        # ✅ Consistent key

# Client configuration
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379  # ✅ Same etcd host
    Key: user.rpc        # ✅ Same key
```

Verify etcd:
```bash
# Check etcd is running
etcdctl version

# List registered services
etcdctl get --prefix user.rpc

# Start etcd if not running
etcd
```

### Issue: RPC call timeout

**Symptoms:**
```
Error: rpc error: code = DeadlineExceeded desc = context deadline exceeded
```

**Solution:**

Increase timeout:

```yaml
# Client configuration
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: user.rpc
  Timeout: 10000  # ✅ Increase to 10 seconds (was 5000)
```

Or pass context with timeout:

```go
// Create timeout context
ctx, cancel := context.WithTimeout(l.ctx, 10*time.Second)
defer cancel()

// Use timeout context
resp, err := l.svcCtx.UserRpc.GetUser(ctx, &user.GetUserRequest{
    Id: 123,
})
```

### Issue: gRPC status error not handled

**Symptoms:**
```go
// Error doesn't match expected type
if err == ErrNotFound {  // Never true
    // ...
}
```

**Solution:**

Use gRPC status:

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// ❌ Wrong: Direct error comparison
if err == model.ErrNotFound {
    return nil, err
}

// ✅ Correct: Check gRPC status
st, ok := status.FromError(err)
if ok {
    switch st.Code() {
    case codes.NotFound:
        return nil, ErrUserNotFound
    case codes.InvalidArgument:
        return nil, ErrInvalidInput
    default:
        return nil, err
    }
}
```

## Middleware Issues

### Issue: Middleware not applied

**Symptoms:**
```
Auth middleware not checking tokens
CORS not working
```

**Solution:**

Check middleware registration:

```go
// In API definition
@server(
    prefix: /api/v1
    group: user
    middleware: Auth  // ✅ Register middleware
)
service user-api {
    @handler GetUser
    get /users/:id (GetUserRequest) returns (GetUserResponse)
}

// In handler registration (generated)
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    server.AddRoutes(
        []rest.Route{
            {
                Method:  http.MethodGet,
                Path:    "/users/:id",
                Handler: GetUserHandler(serverCtx),
            },
        },
        rest.WithPrefix("/api/v1"),
        rest.WithMiddlewares([]rest.Middleware{
            serverCtx.Auth,  // ✅ Middleware applied
        }),
    )
}
```

### Issue: Middleware order wrong

**Symptoms:**
```
Auth middleware runs after logging
Context not available in middleware
```

**Solution:**

Control middleware order:

```go
// ✅ Correct order: Auth before other middlewares
@server(
    prefix: /api/v1
    group: user
    middleware: Auth, RateLimit, Logging  // Auth runs first
)
```

Or use chain:

```go
import "github.com/zeromicro/go-zero/rest/chain"

// Control explicit order
middlewares := chain.New(
    authMiddleware,      // Runs first
    rateLimitMiddleware, // Runs second
    loggingMiddleware,   // Runs last
)
```

## Configuration Issues

### Issue: Config not loaded

**Symptoms:**
```
Error: config file not found
All config values are zero/empty
```

**Solution:**

Check config file path:

```bash
# ❌ Wrong: Relative path may not work
./service -f config.yaml

# ✅ Correct: Use absolute path or explicit relative
./service -f etc/config.yaml
./service -f /app/etc/config.yaml

# Or set working directory
cd /app && ./service -f etc/config.yaml
```

### Issue: Config validation fails

**Symptoms:**
```
Error: invalid config: missing required field
```

**Solution:**

Check required fields:

```go
// In config struct
type Config struct {
    rest.RestConf  // ✅ Must embed

    DataSource string  // ✅ Required (no default, optional, or omitempty)

    // Optional field
    Optional string `json:",optional"`  // ✅ Can be missing

    // With default
    MaxSize int64 `json:",default=1048576"`  // ✅ Has default value
}
```

In YAML:

```yaml
# ✅ Required fields must be present
Name: user-api  # Required from RestConf
Host: 0.0.0.0   # Required from RestConf
Port: 8888      # Required from RestConf
DataSource: "..."  # Required from Config

# Optional fields can be omitted
# MaxSize will use default if not specified
```

## Performance Issues

### Issue: High memory usage

**Symptoms:**
```
OOM (Out of Memory) errors
Memory steadily increasing
```

**Solution:**

Check for goroutine leaks:

```go
// ❌ Wrong: Goroutine never stops
go func() {
    for {
        doWork()
        time.Sleep(time.Second)
    }
}()

// ✅ Correct: Goroutine respects context
go func() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            doWork()
        case <-ctx.Done():
            return  // ✅ Exit goroutine
        }
    }
}()
```

Check connection leaks:

```go
// Use connection pooling (automatic with go-zero)
// Limit concurrent connections
server := rest.MustNewServer(c.RestConf,
    rest.WithMaxConns(1000),  // ✅ Limit connections
)
```

## Logging Issues

### Issue: Logs not showing

**Symptoms:**
```
No log output despite logger calls
```

**Solution:**

Check log level:

```yaml
Log:
  Level: info  # ✅ Ensure level allows your logs
  # debug < info < error < severe
```

Check log mode:

```yaml
Log:
  Mode: console  # ✅ For development (stdout)
  # Mode: file   # For production (writes to file)
```

Ensure proper logger usage:

```go
// ❌ Wrong: Using standard library
log.Println("message")  // Doesn't use go-zero logging

// ✅ Correct: Use logx
l.Logger.Info("message")
logx.Info("message")  // Package level
```

For more help:
- [go-zero Documentation](https://go-zero.dev)
- [GitHub Discussions](https://github.com/zeromicro/go-zero/discussions)
- [GitHub Issues](https://github.com/zeromicro/go-zero/issues)
