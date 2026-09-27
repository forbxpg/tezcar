# 1. Service boundaries

- Status: accepted
- Date: 2026-09-27

## Context

tezcar has to hold long-lived TCP connections to thousands of cars, ingest their
telemetry, run rentals, charge cards, verify identities and report on the fleet. We had to
decide how many deployable services to build.

Options considered:

1. **Monolith.** One process, one database. Simplest to build and operate; one transaction
   covers finishing a rental and charging for it.
2. **Modular monolith.** One process with enforced module boundaries; any module can be
   extracted later.
3. **Microservices.** Separate processes, each owning its data, talking over gRPC and Kafka.

A module earns its own service when it has at least one of:

- a different load profile, so it must scale separately;
- different failure requirements: its failure must not take down the rest, or the reverse;
- a separate security or data boundary;
- a different release cadence.

At pilot scale (50 cars, 150–250 rentals a day) no part of the system needs separate
services for load. tezcar is also a learning project: building sagas, the transactional
outbox and idempotent consumers is one of its goals. We record that reason openly rather
than dress it up as a technical one.

## Decision

We build six services, split along real seams rather than one service per entity:

| Service | Owns | Why it is separate |
|---|---|---|
| Telematics gateway | TCP connections to cars, protocol adapters, vehicle commands | Load profile: thousands of long-lived connections. Failure and release: restarting it drops every car, so it must deploy independently of everything else. |
| Digital twin | Last known state of every car, nearby-car search | Load profile: scales with the telemetry stream, not with users. |
| Core | Rentals, pricing, zones | Not split further: booking and finishing a rental must be one transaction across all three. |
| Payments | Card tokens, holds, charges, refunds, debts, the ledger | Learning goal: sagas, outbox, idempotency. Security boundary: only this service talks to payment gateways and holds their tokens. |
| Identity (KYC) | Documents, selfies, verification results | Data boundary: the law requires biometric data to stay on servers in Uzbekistan, and only one process may touch it. |
| Analytics | Telemetry history and business events in ClickHouse | Load profile and failure: heavy batch writes; may lag or fail without affecting rentals. |

Each service owns its database or schema; no service reads another's tables. Services
exchange events through Kafka using the transactional outbox, and call each other over gRPC
only when the caller needs an answer now.

## Consequences

- Finishing a rental and charging for it are no longer one transaction. Core finishes the
  rental and writes a charge event to its outbox in the same transaction; Payments charges
  later. The client sees the rental as finished and the receipt shortly after.
- Every consumer must be idempotent: events are delivered at least once. Payments deduplicates
  on a unique operation id enforced by the database, not by a check in code.
- A failed charge cannot undo a finished trip. It becomes a debt that blocks new rentals
  and is retried.
- A bug can span several services, so distributed tracing across gRPC and Kafka is required
  from the first service, not added later.
- Six services cost six pipelines, dashboards and on-call runbooks. We accept this cost for
  the learning goal; if it outgrows the benefit, Payments merges back into Core first.
