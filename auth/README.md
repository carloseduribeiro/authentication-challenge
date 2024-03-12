# auth-service

## Migrations

### golang-migrate

**Migrate** reads migrations from sources and applies them in correct order to
a database.

Docs: https://github.com/golang-migrate/migrate

### commands:

```shell
# create a migration
migrate create -ext=sql -dir=internal/infra/database/migrations -seq create_users_table

# execute migrations
migrate -path=internal/infra/database/migrations -database "postgres://root:root@localhost:5432/auth_db?sslmode=disable" -verbose up

# restore migrations
migrate -path=internal/infra/database/migrations -database "postgres://root:root@localhost:5432/auth_db?sslmode=disable" -verbose down
```

## sqlc

"**sqlc** generates fully type-safe idiomatic Go code from SQL."

Docs: https://docs.sqlc.dev/en/latest/#

### commands:

```shell
# Generate sqlc code:
sqlc generate
```