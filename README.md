# tezcar

**Carsharing platform for Uzbekistan.** Book a car in the app, unlock it with your
phone, pay by the minute or by the day. The backend is Go, PostgreSQL with PostGIS,
Kafka, gRPC, Redis and ClickHouse.

> **Status: pre-alpha.** Nothing runs yet. The first milestone is the telematics
> gateway (Teltonika Codec 8E) and a fleet simulator.

## Run the local stack

You need Go (the version in `go.mod`) and Docker.

```bash
make up     # PostgreSQL + PostGIS :5432, Kafka :9092, Redis :6379
make test
make down
```

## Layout

One Go module, one repository:

- `cmd/<service>/` — a service's entry point
- `internal/<service>/` — its code; a service never imports another service's `internal`
- `proto/` — gRPC contracts and Kafka event schemas

## Contributing and security

See [CONTRIBUTING.md](CONTRIBUTING.md); contributions are voluntary and unpaid. Report vulnerabilities privately as described in
[SECURITY.md](SECURITY.md).

## License

tezcar is source-available under the [Business Source License 1.1](LICENSE). You may
read, modify and run it for non-production purposes. Production use requires a
commercial license from the Licensor. Each version converts to Apache License 2.0 four
years after its publication.
