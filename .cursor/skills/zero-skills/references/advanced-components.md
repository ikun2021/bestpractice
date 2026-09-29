# Advanced Components

## Overview

go-zero includes reusable concurrency and data-structure packages under `core/`. This guide covers the current APIs for Bloom filters, timing wheels, MapReduce, batch executors, and singleflight.

Check the project's selected go-zero version before copying examples. Several of these APIs became generic in newer releases.

## Bloom Filter

`core/bloom` uses Redis to provide a probabilistic membership test. A negative result means an item is definitely absent; a positive result may be a false positive and must be verified against authoritative storage.

```go
store := redis.MustNewRedis(redis.RedisConf{
    Host: "redis:6379",
    Type: redis.NodeType,
})
filter := bloom.New(store, "users", 10_000_000)

key := []byte(strconv.FormatInt(userID, 10))
exists, err := filter.ExistsCtx(ctx, key)
if err != nil {
    return nil, fmt.Errorf("check user bloom filter: %w", err)
}
if !exists {
    return nil, ErrUserNotFound
}

// A positive result still requires a cache or database lookup.
return userModel.FindOne(ctx, userID)
```

When creating a record, update the authoritative store first and add the key to the filter after the commit:

```go
if _, err := userModel.Insert(ctx, user); err != nil {
    return err
}
if err := filter.AddCtx(ctx, key); err != nil {
    // Reconcile this failure; otherwise the filter can cause a false negative.
    return fmt.Errorf("add user to bloom filter: %w", err)
}
```

Size the filter from expected item count and acceptable false-positive rate. Do not use a Bloom filter when exact deletion or exact membership is required.

## Timing Wheel

`core/collection.TimingWheel` efficiently schedules many in-process delayed callbacks.

```go
wheel, err := collection.NewTimingWheel(
    time.Second,
    3600,
    func(key, value any) {
        logx.Infof("timer fired: key=%v value=%v", key, value)
    },
)
if err != nil {
    return err
}
defer wheel.Stop()

if err := wheel.SetTimer("order:123", orderID, 30*time.Minute); err != nil {
    return err
}
if err := wheel.MoveTimer("order:123", 45*time.Minute); err != nil {
    return err
}
if err := wheel.RemoveTimer("order:123"); err != nil {
    return err
}
```

`NewTimingWheel` starts its scheduler internally; there is no `Start` method. A timing wheel is in-memory, so it is not appropriate for jobs that must survive process restarts or move between replicas. Use a durable queue or scheduler for those jobs.

## MapReduce

`core/mr.Finish` runs independent functions concurrently and returns when all finish or one returns an error:

```go
var user *user.User
var stock *inventory.Stock

err := mr.Finish(
    func() (err error) {
        user, err = userRPC.GetUser(ctx, &userpb.GetUserReq{UserId: userID})
        return err
    },
    func() (err error) {
        stock, err = inventoryRPC.GetStock(ctx, &inventorypb.GetStockReq{ProductId: productID})
        return err
    },
)
if err != nil {
    return nil, err
}
```

For typed pipelines, current go-zero releases expose generic `mr.MapReduce`:

```go
validIDs, err := mr.MapReduce(
    func(source chan<- int64) {
        for _, id := range userIDs {
            source <- id
        }
    },
    func(id int64, writer mr.Writer[int64], cancel func(error)) {
        valid, err := checkUser(ctx, id)
        if err != nil {
            cancel(err)
            return
        }
        if valid {
            writer.Write(id)
        }
    },
    func(pipe <-chan int64, writer mr.Writer[[]int64], _ func(error)) {
        var result []int64
        for id := range pipe {
            result = append(result, id)
        }
        writer.Write(result)
    },
    mr.WithWorkers(16),
)
```

Call `cancel(err)` for mapper or reducer failures and ensure a non-void reducer writes exactly one result.

## Batch Executors

`core/executors.BulkExecutor` flushes when either the task-count limit or interval is reached:

```go
executor := executors.NewBulkExecutor(
    func(tasks []any) {
        if err := insertBatch(tasks); err != nil {
            // Execute cannot return an error. Persist or emit enough state to
            // retry failed batches; logging alone may lose data.
            logx.Errorf("insert batch: %v", err)
        }
    },
    executors.WithBulkTasks(1000),
    executors.WithBulkInterval(3*time.Second),
)

if err := executor.Add(event); err != nil {
    return err
}
defer func() {
    executor.Flush()
    executor.Wait()
}()
```

`ChunkExecutor` uses caller-supplied byte sizes instead of task counts:

```go
executor := executors.NewChunkExecutor(
    func(tasks []any) {
        sendBatch(tasks)
    },
    executors.WithChunkBytes(1024*1024),
)

if err := executor.Add(payload, len(payload)); err != nil {
    return err
}
```

The executor callback does not return an error. Use it only when callback failures have an explicit recovery path. Flush and wait during shutdown so buffered work is not abandoned.

## SingleFlight

`core/syncx.SingleFlight` coalesces concurrent calls that share a key:

```go
calls := syncx.NewSingleFlight()

value, err := calls.Do(cacheKey, func() (any, error) {
    return userModel.FindOne(ctx, userID)
})
if err != nil {
    return nil, err
}
return value.(*User), nil
```

`DoEx` additionally reports whether this caller executed the function:

```go
value, fresh, err := calls.DoEx(cacheKey, loadUser)
```

Here, `fresh` is `true` for the caller that performed the work and `false` for callers that shared its result. `SingleFlight` does not cache results after the in-flight call completes; combine it with a real cache if repeated calls should reuse a stored value.

Use a consistent, bounded key space. A panic inside the callback is not converted into an error, so callbacks should handle expected failures normally.

## Selection Guide

| Need | Component |
|---|---|
| Reject definitely absent keys before storage lookup | `bloom.Filter` |
| Schedule many best-effort in-process timers | `collection.TimingWheel` |
| Run independent calls concurrently | `mr.Finish` |
| Build a typed concurrent pipeline | `mr.MapReduce` |
| Flush by task count or interval | `executors.BulkExecutor` |
| Flush by payload size or interval | `executors.ChunkExecutor` |
| Coalesce concurrent duplicate work | `syncx.SingleFlight` |

## Checklist

- Propagate `context.Context` where context-aware methods exist.
- Handle constructor and scheduling errors.
- Define recovery behavior for asynchronous callback failures.
- Flush buffered executors during shutdown.
- Do not treat Bloom filters or singleflight as authoritative caches.
- Use durable infrastructure when work must survive restarts.
