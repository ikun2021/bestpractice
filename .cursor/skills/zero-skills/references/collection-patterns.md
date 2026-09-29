# Collection Patterns

go-zero's `core/collection` package provides a set of production-proven, high-performance data structures.

## Overview

| Data Structure | Purpose | Time Complexity |
|---------|------|-----------|
| LRU Cache | Cache frequently accessed data locally | O(1) lookup/update |
| Ring Buffer | Fixed-size circular buffer | O(1) add/read |
| TimingWheel | Efficient timer management | O(1) add/remove/execute |
| SafeMap | General-purpose concurrent map | O(1) lookup/update/delete |

---

## LRU Cache

LRU (Least Recently Used) evicts the data that has gone unused the longest when the cache is full.

### ✅ Basic Usage

```go
import "github.com/zeromicro/go-zero/core/collection"

// Create a cache with a capacity of 100.
cache, err := collection.NewCache(100)
if err != nil {
    log.Fatal(err)
}

// Set a value.
cache.Set("user:1", userData)

// Get a value.
if val, ok := cache.Get("user:1"); ok {
    user := val.(*User)
}

// Delete a value.
cache.Del("user:1")

// Get the cache size.
size := cache.Size()
```

### ✅ With an Eviction Callback

```go
cache, err := collection.NewCache(100, collection.WithEvict(func(key string, value interface{}) {
    // Clean up resources.
    if closer, ok := value.(io.Closer); ok {
        closer.Close()
    }
    log.Printf("evicted: %s", key)
}))
```

### ✅ Using It in ServiceContext

```go
type ServiceContext struct {
    Config     config.Config
    UserCache  *collection.Cache
    UsersModel model.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
    cache, _ := collection.NewCache(1000, collection.WithEvict(func(key string, value interface{}) {
        logx.Infof("cache evicted: %s", key)
    }))

    return &ServiceContext{
        Config:    c,
        UserCache: cache,
        UsersModel: model.NewUsersModel(sqlx.NewMysql(c.DataSource), c.Cache),
    }
}
```

### ✅ Using the Cache in the Logic Layer

```go
func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (*types.GetUserResponse, error) {
    cacheKey := fmt.Sprintf("user:%d", req.Id)

    // Check the local cache first.
    if val, ok := l.svcCtx.UserCache.Get(cacheKey); ok {
        user := val.(*model.Users)
        return &types.GetUserResponse{
            Id:    user.Id,
            Name:  user.Name,
            Email: user.Email,
        }, nil
    }

    // On a local cache miss, query the database.
    user, err := l.svcCtx.UsersModel.FindOne(l.ctx, req.Id)
    if err != nil {
        return nil, err
    }

    // Populate the local cache.
    l.svcCtx.UserCache.Set(cacheKey, user)

    return &types.GetUserResponse{
        Id:    user.Id,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}
```

### ❌ Common Mistakes

```go
// Wrong: an excessively large capacity can exhaust memory.
cache, _ := collection.NewCache(10000000)  // Too large!

// Wrong: the cache retains a resource without cleaning it up.
cache.Set("conn", dbConn)  // The cached connection cannot be released.

// Wrong: cache key collision.
cache.Set("user", userA)
cache.Set("user", userB)  // Overwrites userA.
```

### Best Practices

- Choose a sensible cache size based on available memory and the expected hit rate.
- Run expensive eviction-callback work asynchronously.
- Use meaningful key names to avoid collisions.
- Keep cached objects small and immutable when possible.

---

## Ring Buffer

A ring buffer is a fixed-size queue in which new elements overwrite the oldest elements.

### ✅ Basic Usage

```go
import "github.com/zeromicro/go-zero/core/collection"

// Create a ring buffer that retains the 100 most recent records.
ring := collection.NewRing(100)

// Add elements.
ring.Add(logEntry1)
ring.Add(logEntry2)

// Get all elements.
items := ring.Take()
for _, item := range items {
    log.Printf("%v", item)
}
```

### ✅ Log Buffer Pattern

```go
type LogBuffer struct {
    ring *collection.Ring
    mu   sync.RWMutex
}

func NewLogBuffer(size int) *LogBuffer {
    return &LogBuffer{
        ring: collection.NewRing(size),
    }
}

func (lb *LogBuffer) AddLog(entry LogEntry) {
    lb.mu.Lock()
    defer lb.mu.Unlock()
    lb.ring.Add(entry)
}

func (lb *LogBuffer) GetRecentLogs() []LogEntry {
    lb.mu.RLock()
    defer lb.mu.RUnlock()

    items := lb.ring.Take()
    logs := make([]LogEntry, 0, len(items))
    for _, item := range items {
        logs = append(logs, item.(LogEntry))
    }
    return logs
}
```

### ✅ Fixed-Window Metrics

```go
type MetricsCollector struct {
    requestTimes *collection.Ring
}

func NewMetricsCollector(windowSize int) *MetricsCollector {
    return &MetricsCollector{
        requestTimes: collection.NewRing(windowSize),
    }
}

func (mc *MetricsCollector) RecordRequest(duration time.Duration) {
    mc.requestTimes.Add(duration)
}

func (mc *MetricsCollector) GetAverageRequestTime() time.Duration {
    items := mc.requestTimes.Take()
    if len(items) == 0 {
        return 0
    }

    var total time.Duration
    for _, item := range items {
        total += item.(time.Duration)
    }
    return total / time.Duration(len(items))
}
```

### ❌ Common Mistakes

```go
// Wrong: an excessive size wastes memory.
ring := collection.NewRing(1000000)  // Too large!

// Wrong: storing large objects consumes too much memory.
ring.Add(largeFileContent)  // Do not store large objects.

// Wrong: assuming the wrong element order.
items := ring.Take()
// Elements are ordered from oldest to newest, not the reverse.
```

### Best Practices

- Choose a sensible size based on the data volume.
- Store lightweight objects or references.
- Read periodically to prevent data from accumulating.
- Use an additional index when strict ordering guarantees are required.

---

## TimingWheel

A timing wheel is an efficient timer implementation suited to managing large numbers of scheduled tasks.

### ✅ Basic Usage

```go
import "github.com/zeromicro/go-zero/core/collection"

// Create a timing wheel with 100 ms precision and 3,600 slots
// (supporting delays of up to 6 minutes).
tw, err := collection.NewTimingWheel(100*time.Millisecond, 3600)
if err != nil {
    log.Fatal(err)
}
defer tw.Stop()

// Set a timer.
tw.SetTimer("task-1", nil, 5*time.Second, func(key string, value interface{}) {
    log.Printf("timer triggered: %s", key)
})

// Cancel a timer.
tw.RemoveTimer("task-1")
```

### ✅ Choosing Parameters

```go
// Calculate the number of slots from the maximum delay.
maxDelay := 10 * time.Minute          // Support delays of up to 10 minutes.
interval := 100 * time.Millisecond   // 100 ms precision.
numSlots := int(maxDelay / interval) // 6,000 slots.

tw, _ := collection.NewTimingWheel(interval, numSlots)
```

**Parameter notes:**
- `interval`: tick precision; smaller values improve precision but increase CPU overhead.
- `numSlots`: number of slots; affects memory usage and the maximum delay.
- `maximum delay = interval * numSlots`

### ✅ Request Timeout Management

```go
type RequestManager struct {
    tw         *collection.TimingWheel
    requests   map[string]*Request
    mu         sync.RWMutex
}

func NewRequestManager() *RequestManager {
    tw, _ := collection.NewTimingWheel(100*time.Millisecond, 3600)
    return &RequestManager{
        tw:       tw,
        requests: make(map[string]*Request),
    }
}

func (rm *RequestManager) AddRequest(req *Request) {
    rm.mu.Lock()
    defer rm.mu.Unlock()

    rm.requests[req.ID] = req

    // Set a 30-second timeout.
    rm.tw.SetTimer(req.ID, nil, 30*time.Second, func(key string, value interface{}) {
        rm.handleTimeout(key)
    })
}

func (rm *RequestManager) CompleteRequest(id string) {
    rm.mu.Lock()
    defer rm.mu.Unlock()

    delete(rm.requests, id)
    rm.tw.RemoveTimer(id)
}

func (rm *RequestManager) handleTimeout(id string) {
    rm.mu.Lock()
    defer rm.mu.Unlock()

    if req, ok := rm.requests[id]; ok {
        req.Timeout()
        delete(rm.requests, id)
    }
}
```

### ✅ Delayed Task Scheduling

```go
type TaskScheduler struct {
    tw *collection.TimingWheel
}

func NewTaskScheduler() *TaskScheduler {
    tw, _ := collection.NewTimingWheel(time.Second, 3600)
    return &TaskScheduler{tw: tw}
}

func (s *TaskScheduler) ScheduleTask(taskID string, delay time.Duration, task func()) {
    s.tw.SetTimer(taskID, nil, delay, func(key string, value interface{}) {
        // Run expensive work asynchronously to avoid blocking the timing wheel.
        go task()
    })
}

func (s *TaskScheduler) CancelTask(taskID string) {
    s.tw.RemoveTimer(taskID)
}
```

### ✅ Heartbeat Monitoring

```go
type HeartbeatManager struct {
    tw        *collection.TimingWheel
    clients   map[string]time.Time
    timeout   time.Duration
    onTimeout func(clientID string)
}

func NewHeartbeatManager(timeout time.Duration, onTimeout func(string)) *HeartbeatManager {
    tw, _ := collection.NewTimingWheel(time.Second, int(timeout/time.Second)*2)

    hm := &HeartbeatManager{
        tw:        tw,
        clients:   make(map[string]time.Time),
        timeout:   timeout,
        onTimeout: onTimeout,
    }

    // Check periodically.
    tw.SetTimer("heartbeat-check", nil, time.Second, hm.checkHeartbeats)
    return hm
}

func (hm *HeartbeatManager) UpdateHeartbeat(clientID string) {
    hm.clients[clientID] = time.Now()
}

func (hm *HeartbeatManager) checkHeartbeats(key string, value interface{}) {
    now := time.Now()
    for clientID, lastBeat := range hm.clients {
        if now.Sub(lastBeat) > hm.timeout {
            go hm.onTimeout(clientID)
            delete(hm.clients, clientID)
        }
    }

    // Schedule the next check.
    hm.tw.SetTimer("heartbeat-check", nil, time.Second, hm.checkHeartbeats)
}
```

### ❌ Common Mistakes

```go
// Wrong: too few slots to support the required delay.
tw, _ := collection.NewTimingWheel(time.Second, 10)  // Supports only 10 seconds!

// Wrong: expensive callback work blocks the timing wheel.
tw.SetTimer("task", nil, delay, func(key string, value interface{}) {
    time.Sleep(10 * time.Second) // Blocks!
    processLargeData()           // Expensive operation!
})

// Wrong: forgetting to stop the timing wheel during shutdown.
// defer tw.Stop() // This call is required!
```

### Best Practices

- Use a goroutine for expensive work in callbacks.
- Always call `tw.Stop()` during shutdown.
- Choose suitable precision and slot counts for your requirements.
- Use meaningful keys to simplify debugging and management.

---

## SafeMap

A general-purpose concurrent map suited to read-heavy workloads with infrequent writes.

### ✅ Basic Usage

```go
import "github.com/zeromicro/go-zero/core/collection"

m := collection.NewSafeMap()

// Set a value.
m.Set("key", "value")

// Get a value.
if val, ok := m.Get("key"); ok {
    fmt.Println(val.(string))
}

// Delete a value.
m.Del("key")

// Iterate over all entries.
m.Range(func(key string, value interface{}) bool {
    fmt.Printf("%s: %v\n", key, value)
    return true // Continue iterating.
})
```

### ✅ Service Instance Registry

```go
type ServiceRegistry struct {
    services *collection.SafeMap
}

func NewServiceRegistry() *ServiceRegistry {
    return &ServiceRegistry{
        services: collection.NewSafeMap(),
    }
}

func (r *ServiceRegistry) Register(serviceID, address string) {
    r.services.Set(serviceID, &ServiceInfo{
        Address:   address,
        LastSeen:  time.Now(),
    })
}

func (r *ServiceRegistry) Deregister(serviceID string) {
    r.services.Del(serviceID)
}

func (r *ServiceRegistry) GetService(serviceID string) (*ServiceInfo, bool) {
    if val, ok := r.services.Get(serviceID); ok {
        return val.(*ServiceInfo), true
    }
    return nil, false
}

func (r *ServiceRegistry) ListServices() []*ServiceInfo {
    var services []*ServiceInfo
    r.services.Range(func(key string, value interface{}) bool {
        services = append(services, value.(*ServiceInfo))
        return true
    })
    return services
}
```

### ✅ Shared Configuration Store

```go
type ConfigStore struct {
    config *collection.SafeMap
}

func NewConfigStore() *ConfigStore {
    return &ConfigStore{
        config: collection.NewSafeMap(),
    }
}

func (s *ConfigStore) Set(key string, value interface{}) {
    s.config.Set(key, value)
}

func (s *ConfigStore) Get(key string) (interface{}, bool) {
    return s.config.Get(key)
}

func (s *ConfigStore) GetInt(key string, defaultVal int) int {
    if val, ok := s.config.Get(key); ok {
        return val.(int)
    }
    return defaultVal
}

func (s *ConfigStore) GetString(key string, defaultVal string) string {
    if val, ok := s.config.Get(key); ok {
        return val.(string)
    }
    return defaultVal
}
```

### ✅ Choosing Between SafeMap and sync.Map

| Scenario | Recommendation | Reason |
|------|------|------|
| Read-heavy, infrequent writes, iteration required | SafeMap | Simple and direct; convenient `Range` support |
| Write once, read many times | sync.Map | Optimized for this workload |
| Frequent writes | Neither | Consider a sharded-lock design |

### ⚠️ Memory Management Considerations

Go maps, including `sync.Map` and `SafeMap`, have an important property: **deleting entries does not automatically release their allocated memory**.

```go
// Problem scenario.
m := collection.NewSafeMap()

// Insert many entries.
for i := 0; i < 1000000; i++ {
    m.Set(fmt.Sprintf("key-%d", i), i)
}

// Delete many entries.
for i := 0; i < 1000000; i++ {
    m.Del(fmt.Sprintf("key-%d", i))
}

// Memory usage remains high because the backing storage does not shrink.
```

### ✅ Rebuild Periodically to Reclaim Memory

```go
type ServiceRegistry struct {
    services *collection.SafeMap
    mu       sync.Mutex
}

func (r *ServiceRegistry) Compact() {
    r.mu.Lock()
    defer r.mu.Unlock()

    // Create a new map and migrate live entries.
    newMap := collection.NewSafeMap()
    r.services.Range(func(key string, value interface{}) bool {
        newMap.Set(key, value)
        return true
    })
    r.services = newMap
}

// Compact automatically on a schedule.
func (r *ServiceRegistry) StartCompactor(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            r.Compact()
        }
    }()
}
```

### ❌ Common Mistakes

```go
// Wrong: modifying the map while iterating with Range.
m.Range(func(key string, value interface{}) bool {
    m.Del(key) // May cause problems!
    return true
})

// Wrong: frequent writes cause lock contention.
for i := 0; i < 10000; i++ {
    go m.Set(key, value) // SafeMap is not suitable for frequent writes.
}

// Wrong: inserting and then deleting many entries wastes memory.
for i := 0; i < 1000000; i++ {
    m.Set(key, value)
    m.Del(key) // The memory is not released!
}
```

### Best Practices

- Use it for read-heavy workloads with infrequent writes that require iteration.
- Avoid modifying the map inside `Range`.
- Rebuild it periodically when its size fluctuates significantly.
- Consider a sharded-lock design for write-heavy workloads.

---

## Complete Example

### Multi-Level Cache Architecture

```go
type CacheService struct {
    localCache *collection.Cache       // L1: local cache
    redisCache *redis.Redis            // L2: Redis cache
    db         sqlx.SqlConn            // L3: database
    requestLog *collection.Ring        // Request log
    timeoutMgr *collection.TimingWheel // Timeout management
}

func NewCacheService(cfg config.Config) *CacheService {
    localCache, _ := collection.NewCache(1000, collection.WithEvict(func(key string, value interface{}) {
        logx.Infof("local cache evicted: %s", key)
    }))

    return &CacheService{
        localCache:  localCache,
        redisCache:  redis.MustNewRedis(cfg.Redis),
        db:          sqlx.NewMysql(cfg.DataSource),
        requestLog:  collection.NewRing(100),
        timeoutMgr:  mustNewTimingWheel(),
    }
}

func (s *CacheService) Get(ctx context.Context, key string) (interface{}, error) {
    // L1: local cache.
    if val, ok := s.localCache.Get(key); ok {
        s.requestLog.Add("L1 hit: " + key)
        return val, nil
    }

    // L2: Redis cache.
    val, err := s.redisCache.Get(key)
    if err == nil {
        s.localCache.Set(key, val)
        s.requestLog.Add("L2 hit: " + key)
        return val, nil
    }

    // L3: database.
    val, err = s.queryDB(ctx, key)
    if err != nil {
        return nil, err
    }

    // Populate the caches.
    s.localCache.Set(key, val)
    s.redisCache.Set(key, val)

    return val, nil
}
```

---

## Best-Practice Summary

### ✅ DO

| Data Structure | Best Practices |
|---------|---------|
| LRU Cache | Choose a sensible size; clean up resources in eviction callbacks; run expensive work asynchronously |
| Ring Buffer | Use a fixed size; store lightweight objects; read periodically |
| TimingWheel | Use goroutines in callbacks; call `Stop()` during shutdown |
| SafeMap | Use for read-heavy workloads; rebuild periodically to reclaim memory |

### ❌ DON'T

| Data Structure | Avoid |
|---------|---------|
| LRU Cache | Excessive capacity; caching large objects; key collisions |
| Ring Buffer | Storing large objects; excessive capacity |
| TimingWheel | Blocking callbacks; too few slots; forgetting `Stop()` |
| SafeMap | Modifying during `Range`; frequent writes; ignoring memory usage |

---

## Related Resources

- [REST API Patterns](./rest-api-patterns.md) - API-layer caching strategies
- [Advanced Components](./advanced-components.md) - Bloom filter and related helpers
- [Resilience Patterns](./resilience-patterns.md) - Resilience patterns such as circuit breaking and rate limiting
