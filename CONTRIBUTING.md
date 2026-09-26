# Contributing

## Set up

You need Go (the version in `go.mod`), Docker, and
[golangci-lint](https://golangci-lint.run/) v2.

```bash
make up     # local PostgreSQL + PostGIS, Kafka, Redis
```

## Checks

All of these must pass; CI runs the same:

```bash
make lint test
```

`make lint` runs golangci-lint with the formatters (`gofmt`, `goimports`);
`make test` runs `go test -race ./...`.

## Conventions

- English everywhere: code, comments, docs, commit messages.
- Every Go file starts with `// SPDX-License-Identifier: BUSL-1.1`.
- Every package has a doc comment; every exported name has one.
- `context.Context` is the first parameter of anything that does I/O.
- Errors are wrapped with `%w` and context (`fmt.Errorf("load rental %s: %w", id, err)`),
  never swallowed.
- SQL is written by hand (`pgx`, `sqlc`); no ORM.
- Money is an integer number of tiyin (`int64`), never a float.
- A service never imports another service's `internal` package; services talk over gRPC
  or Kafka.

## Commits and pull requests

Commit messages and PR titles follow [Conventional Commits](https://www.conventionalcommits.org/)
(`feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `build:`, `ci:`, `chore:`; `!` for
breaking changes). A scope names the service: `feat(telematics): parse Codec 8E`. The PR
title becomes the squashed commit.

Every commit is signed off (`git commit -s`), see below.

## Licensing of contributions

tezcar is licensed under the [Business Source License 1.1](LICENSE), and the Licensor also
licenses it commercially. To keep both possible, a contribution is accepted only on these
terms:

1. **Developer Certificate of Origin.** Your sign-off (`Signed-off-by: Name <email>`, added
   by `git commit -s`) certifies the [Developer Certificate of Origin 1.1](https://developercertificate.org/):
   you wrote the contribution or otherwise have the right to submit it.
2. **License grant.** By submitting a contribution you grant forbxpg a perpetual,
   worldwide, non-exclusive, royalty-free, irrevocable license to use, reproduce, modify,
   distribute, sublicense and relicense it, under the Business Source License, under the
   Change License, or under commercial terms; and a patent license on the same terms for
   any patent claims you can license that the contribution necessarily infringes. You keep
   the copyright in your contribution.

If you cannot agree to these terms, open an issue instead of a pull request.

## Security

Found a vulnerability? Do not open an issue: follow [SECURITY.md](SECURITY.md).
