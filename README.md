# test-rockets-project

## Summary

Small PoC service to receive rocket events (webhook), process them asynchronously with per-channel ordering guarantees, and expose a REST API to read current rocket state. The implementation focuses on correctness for ordering and idempotency; storage is in-memory to keep the submission self-contained and easy to run.

> Go version: declared in `go.mod` (`go 1.25.3`).

## Endpoints

- `POST /messages`  — webhook for incoming events
- `GET /rockets`    — list rockets
- `GET /rockets/{channel}` — get rocket by channel

Example webhook payload:

```json
{
  "metadata": {
    "channel": "193270a9-c9cf-404a-8f83-838e71d9ae67",
    "messageNumber": 1,
    "messageTime": "2022-02-02T19:39:05.86337+01:00",
    "messageType": "RocketLaunched"
  },
  "message": {
    "type": "Falcon-9",
    "launchSpeed": 500,
    "mission": "ARTEMIS"
  }
}
```

## How to run (local)

Prerequisites: Go installed.

Run with `go run` (PowerShell):

```powershell
cd .\test-rockets-project
go run .\
```

Or build and run (bash):

```bash
cd test-rockets-project
go build -o test-rockets-project
./test-rockets-project
```

The service listens by default at `http://localhost:8080`.

## Using the provided `rockets` test runner

Run the platform-appropriate executable from the ZIP and point it to the webhook URL.

PowerShell example:

```powershell
.\rockets.exe launch "http://localhost:8080/messages" --message-delay=500ms --concurrency-level=1
```

Bash example:

```bash
./rockets launch "http://localhost:8080/messages" --message-delay=500ms --concurrency-level=1
```

Note: the challenge README uses `:8088` as an example; this implementation uses `:8080`. The runner accepts any URL.

## Key implemented behaviour

- Incoming events are stored in an `EventStore` and marked `pending`.
- A background consumer polls pending events and only processes an event when its `messageNumber` equals the expected next number for that `channel` (tracked in memory). This enforces per-channel ordering.
- The composite key `channel-messageNumber` is used for deduplication (idempotency) — duplicates are ignored.
- `Rocket` model methods (`ApplyLaunchEvent`, `ApplySpeedIncreasedEvent`, etc.) encapsulate domain rules (e.g. no speed changes after explosion).

## Design decisions and trade-offs (concise)

- Interfaces (`EventStore`, `RocketStore`): chosen to allow swapping in-memory stores for durable implementations (Postgres/Redis) with minimal changes.
- Composite id key for idempotency (`channel-messageNumber`): simple, reliable given the message contract.
- In-memory `next[channel]` ordering: straightforward for single-instance consumer; to scale horizontally consider partitioning by channel (e.g. Kafka partitions) or DB-based locking.
- Polling consumer (1s interval): simple and reliable for PoC; in production a push-based worker or broker reduces latency.

## Limitations / Not implemented (intentional)

- In-memory persistence — events and state are lost on restart.
- No DLQ or bounded retry policy — events that always fail remain `pending` indefinitely.
- Webhook is unauthenticated (acceptable for local testing only).
- `messageTime` is stored as string and not used for ordering; ordering relies on `messageNumber`.
- No pagination or sorting params on `GET /rockets` (the store returns an unsorted slice).

These items are documented as deliberate trade-offs due to time constraints; they are straightforward to implement given the current abstractions.

## How to validate behaviour quickly

1. Start service (`go run .`).
2. Run `rockets` launcher pointing to `http://localhost:8080/messages`.
3. Inspect `GET /rockets` and `GET /rockets/{channel}` to verify final state.

## Notes for reviewers / interview defenders

The repository contains clear separation of concerns: handlers (ingest), stores (persistence), consumer (ordering), and model (domain rules). The most important guarantees — per-channel ordering and idempotency — are implemented. Other operational hardening (durability, DLQ, auth, graceful shutdown, metrics) are documented as next steps and can be added without changing domain logic thanks to the interfaces.

If you want, I can add a short `DELIVERY.md` with a one-page checklist for evaluators.