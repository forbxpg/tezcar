# tezcar

Платформа каршеринга для Узбекистана: Go, PostgreSQL/PostGIS, Kafka, gRPC, Redis, ClickHouse.

## Запуск окружения

```bash
make up    # Postgres+PostGIS :5432, Kafka :9092, Redis :6379
make test
make down
```

## Раскладка

Один Go-модуль, монорепозиторий:

- `cmd/<service>/` — точка входа сервиса
- `internal/<service>/` — код сервиса, чужой `internal` не импортируем
- `proto/` — gRPC-контракты и события Kafka
