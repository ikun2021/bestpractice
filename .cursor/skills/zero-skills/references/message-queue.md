# Message Queue Patterns

## Overview

[go-queue](https://github.com/zeromicro/go-queue) provides two queue implementations commonly used with go-zero services:

| Package | Backend | Best fit |
|---|---|---|
| `dq` | Beanstalkd + Redis | Scheduled and delayed jobs |
| `kq` | Kafka | Event streams and consumer groups |

Both systems can deliver a message more than once. Make handlers idempotent and treat malformed or permanently failing messages explicitly rather than silently retrying forever.

## Delayed Queue (`dq`)

`dq` stores delayed jobs in Beanstalkd and uses Redis to coordinate consumption across nodes.

### Configuration

```yaml
Name: order-jobs

DqConf:
  Beanstalks:
    - Endpoint: beanstalkd-1:11300
      Tube: orders
    - Endpoint: beanstalkd-2:11300
      Tube: orders
  Redis:
    Host: redis:6379
    Type: node
    Pass: ""
```

```go
type Config struct {
    service.ServiceConf
    DqConf dq.DqConf
}
```

### Producer

```go
type OrderProducer struct {
    producer dq.Producer
}

func NewOrderProducer(c dq.DqConf) *OrderProducer {
    return &OrderProducer{producer: dq.NewProducer(c.Beanstalks)}
}

func (p *OrderProducer) ScheduleTimeout(orderID int64, delay time.Duration) error {
    body, err := json.Marshal(struct {
        Type    string `json:"type"`
        OrderID int64  `json:"order_id"`
    }{Type: "order_timeout", OrderID: orderID})
    if err != nil {
        return err
    }

    _, err = p.producer.Delay(body, delay)
    return err
}

func (p *OrderProducer) ScheduleAt(body []byte, at time.Time) error {
    _, err := p.producer.At(body, at)
    return err
}
```

Use the identifier returned by `Delay` or `At` when the application needs to track the scheduled job.

### Consumer

`Consumer.Consume` starts the worker group and blocks. Run it as the process's main workload or in a goroutine whose lifecycle is owned by the process.

```go
consumer := dq.NewConsumer(c.DqConf)
consumer.Consume(func(body []byte) {
    var task OrderTask
    if err := json.Unmarshal(body, &task); err != nil {
        logx.Errorf("discarding malformed order task: %v", err)
        return
    }

    if err := handler.Process(context.Background(), task); err != nil {
        // Decide whether to retry, send to a dead-letter flow, or record the
        // task for manual recovery. Keep Process idempotent.
        logx.Errorf("process order task %d: %v", task.OrderID, err)
    }
})
```

The current `dq.Consumer` interface exposes `Consume` but no `Stop` method. Do not wrap it in a fake `service.Service` whose `Stop` method cannot stop the underlying consumers. Use process-level shutdown and verify the go-queue version's lifecycle behavior before embedding it in a larger service group.

## Kafka Queue (`kq`)

### Configuration

```yaml
Name: order-events
Brokers:
  - kafka-1:9092
  - kafka-2:9092
Group: order-processor
Topic: orders
Offset: first
Conns: 1
Consumers: 8
Processors: 8
ForceCommit: true
```

Load this directly into `kq.KqConf`, or nest it in an application configuration field.

### Consumer

`kq.MustNewQueue` returns a `queue.MessageQueue`, which implements `Start` and `Stop`. The handler receives the message context, key, and value.

```go
type OrderEventHandler struct {
    svcCtx *svc.ServiceContext
}

func (h *OrderEventHandler) Consume(ctx context.Context, key, value string) error {
    var event OrderEvent
    if err := json.Unmarshal([]byte(value), &event); err != nil {
        // Returning nil acknowledges input that can never be decoded. Record
        // it first if the system requires a dead-letter trail.
        logx.WithContext(ctx).Errorf("invalid order event key=%q: %v", key, err)
        return nil
    }

    return h.svcCtx.OrderEvents.Apply(ctx, event)
}

handler := &OrderEventHandler{svcCtx: svcCtx}
q := kq.MustNewQueue(c.KqConsumerConf, handler)
defer q.Stop()
q.Start()
```

For a small handler, `kq.WithHandle` avoids defining a type:

```go
q := kq.MustNewQueue(c.KqConsumerConf, kq.WithHandle(
    func(ctx context.Context, key, value string) error {
        return processEvent(ctx, key, value)
    },
))
defer q.Stop()
q.Start()
```

### Producer

Use `kq.Pusher`; there is no `kq.Producer` type.

```go
pusher := kq.NewPusher(c.Brokers, c.Topic)
defer func() {
    if err := pusher.Close(); err != nil {
        logx.Errorf("close Kafka pusher: %v", err)
    }
}()

body, err := json.Marshal(event)
if err != nil {
    return err
}

// A stable domain key keeps events for the same aggregate on one partition.
return pusher.PushWithKey(ctx, strconv.FormatInt(event.OrderID, 10), string(body))
```

`NewPusher` buffers by default. Use `kq.WithSyncPush()` only when the caller must observe the broker write synchronously, and always close the pusher so buffered messages are flushed.

## Retry and Idempotency

Retries are application policy, not a substitute for idempotency:

```go
func (h *OrderEventHandler) Consume(ctx context.Context, key, value string) error {
    event, err := decodeOrderEvent(value)
    if err != nil {
        return nil // malformed input is not transient; record it as needed
    }

    applied, err := h.svcCtx.ProcessedEvents.TryStart(ctx, event.ID)
    if err != nil {
        return err
    }
    if !applied {
        return nil // duplicate delivery
    }

    return h.svcCtx.OrderEvents.Apply(ctx, event)
}
```

The idempotency record and business update should share a transaction where the storage system supports it. If they cannot, design a recoverable state transition instead of assuming exactly-once delivery.

## Choosing Between `dq` and `kq`

| Requirement | Prefer |
|---|---|
| Execute at a future time | `dq` |
| Consumer groups and replay | `kq` |
| Per-key ordering | `kq` with stable message keys |
| Simple scheduled background jobs | `dq` |
| High-throughput event stream | `kq` |

## Checklist

- Keep handlers idempotent.
- Bound processing time and propagate `context.Context` where the API provides it.
- Distinguish malformed messages from transient failures.
- Define retry and dead-letter behavior for the application.
- Use stable Kafka keys when ordering by aggregate matters.
- Close `kq.Pusher` and stop `kq` consumers during shutdown.
- Load broker addresses, credentials, topics, and groups from configuration.

## References

- [go-queue repository](https://github.com/zeromicro/go-queue)
- [go-queue `dq` package](https://pkg.go.dev/github.com/zeromicro/go-queue/dq)
- [go-queue `kq` package](https://pkg.go.dev/github.com/zeromicro/go-queue/kq)
