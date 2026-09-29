# Distributed Transaction Patterns

## Scope

go-zero does not make a set of independent service calls atomic by itself. [DTM](https://github.com/dtm-labs/dtm) provides workflow, Saga, TCC, XA, two-phase message, and outbox-oriented patterns, with a separate [go-zero driver](https://github.com/dtm-labs/dtmdriver-gozero) for framework integration.

Distributed transactions add operational and correctness costs. Prefer a single local database transaction when one service and one datastore can own the invariant.

## Choose the Pattern Deliberately

| Pattern | Use when | Main design obligation |
|---|---|---|
| Workflow | A durable, code-oriented orchestration can coordinate forward and rollback branches | Deterministic workflow and tested rollback behavior |
| Saga | Each completed step has a compensating action | Compensation is business-correct and idempotent |
| TCC | Resources can be reserved, confirmed, and cancelled | Try/Confirm/Cancel semantics and timeout handling |
| Two-phase message | A local transaction must reliably cause later work | Query/preparation logic and idempotent message handling |
| XA | Participating databases and workload tolerate coordinator and lock costs | Driver support, timeout, and contention analysis |
| Outbox | Event publication should follow a local database commit | Relay reliability, deduplication, and ordering policy |

Do not choose a pattern only because its happy-path example is short. Start with the failure semantics the business can accept.

## HTTP Saga Skeleton

The DTM HTTP client lives in `github.com/dtm-labs/client/dtmcli`. Generate the transaction ID with `dtmcli.MustGenGid`, add forward and compensating endpoints in matching pairs, and submit the Saga:

```go
const dtmServer = "http://dtm:36789/api/dtmsvr"

gid := dtmcli.MustGenGid(dtmServer)
payload := &OrderRequest{
    OrderID:  orderID,
    ProductID: productID,
    Quantity: quantity,
}

err := dtmcli.NewSaga(dtmServer, gid).
    Add(orderService+"/orders/create", orderService+"/orders/create-compensate", payload).
    Add(stockService+"/stock/deduct", stockService+"/stock/deduct-compensate", payload).
    Submit()
if err != nil {
    return fmt.Errorf("submit order saga %s: %w", gid, err)
}
```

This is only the coordinator call. Every branch endpoint still needs authentication, input validation, idempotency, observability, and durable local state changes.

## Branch Barriers

DTM branch barriers address common retry anomalies such as duplicate execution, empty compensation, and suspension. Use the barrier APIs demonstrated by the DTM examples for the exact protocol and datastore in your application.

Important rules:

- Create the barrier from DTM's branch metadata, not from user-controlled transaction identifiers.
- Execute the barrier record and the business mutation in the same local transaction.
- Keep forward, compensation, confirm, and cancel handlers idempotent.
- Never perform an irreversible external side effect inside the local transaction without its own deduplication key.
- Install and migrate the DTM barrier table required by the selected database driver.

The SQL transaction type accepted by a barrier is driver-specific. Do not copy an example using `*sql.Tx` into a project that uses `sqlx.Session`, GORM, MongoDB, or Redis without selecting the matching DTM adapter.

## TCC Review Checklist

For each resource, define:

1. **Try** — validates the request and reserves the resource without making the final effect visible.
2. **Confirm** — commits the reservation idempotently.
3. **Cancel** — releases the reservation idempotently, including an empty cancel that arrives before Try.

Check these failure cases in tests:

- Confirm is delivered more than once.
- Cancel is delivered more than once.
- Cancel arrives before Try.
- Try succeeds locally but its response is lost.
- One branch times out after another branch has reserved resources.
- The coordinator is unavailable during submission or status polling.

## Saga Compensation Checklist

Compensation is not a database rollback. It is a new business action and may run much later.

- Store the original business identifiers needed by compensation.
- Make compensation safe after partial execution.
- Decide what happens when compensation itself fails repeatedly.
- Avoid lossy compensation, such as recreating inventory without preserving allocation details.
- Record the global transaction ID in logs and business audit data.
- Protect branch endpoints so arbitrary callers cannot trigger compensations.

## Two-Phase Message and Outbox Guidance

Use a two-phase message or outbox when the core invariant is “commit local state, then reliably publish work.” Do not describe either pattern as exactly-once delivery. Consumers must still deduplicate.

For an outbox:

1. Write the business record and outbox row in one local transaction.
2. Relay unpublished rows to the broker.
3. Mark publication using a recoverable state transition.
4. Give each event a stable identifier.
5. Make consumers idempotent and define ordering per aggregate.

For DTM two-phase messages, follow the current `dtm-labs/dtm-examples` implementation for `DoAndSubmit`, query preparation, and barriers. These APIs vary between the HTTP, gRPC, and workflow clients.

## Configuration and Security

```go
type Config struct {
    rest.RestConf
    DTM struct {
        Server string
    }
}
```

```yaml
DTM:
  Server: http://dtm:36789/api/dtmsvr
```

- Keep the DTM server and branch endpoints on trusted networks.
- Apply service authentication supported by the chosen protocol and driver.
- Set explicit timeouts for coordinator and branch calls.
- Do not log full payloads when they contain credentials or personal data.
- Monitor transactions that remain prepared, submitted, or compensating beyond expected durations.

## Verification Strategy

Integration tests should run a real DTM server and the same database family used in production. Inject failures at branch boundaries and assert both business state and transaction state.

At minimum, test:

- forward success;
- forward failure and successful compensation;
- duplicate branch delivery;
- response loss after local commit;
- service timeout;
- coordinator restart;
- compensation failure followed by retry;
- concurrent transactions for the same business resource.

## References

- [DTM repository](https://github.com/dtm-labs/dtm)
- [DTM Go client](https://github.com/dtm-labs/client)
- [DTM examples](https://github.com/dtm-labs/dtm-examples)
- [DTM go-zero driver](https://github.com/dtm-labs/dtmdriver-gozero)
