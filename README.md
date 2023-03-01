<h1>Kawaii Shop</h1>

<h2>Start PostgreSQLon Docker 🐋</h2>

```bash
docker run --name kawaii_db_test -e POSTGRES_USER=kawaii -e POSTGRES_PASSWORD=123456 -p 5432:5432 -d postgres:alpine
```

<h2>Execute a container and CREATE a new database</h2>

```bash
docker exec -it kawaii_db_test bash
psql -h localhost -U kawaii -W
CREATE DATABASE kawaii_db_test;
\l
```

<h2>Migrate command</h2>

```bash
# Migrate up
migrate -database 'postgres://kawaii:123456@0.0.0.0:5432/kawaii_db_test?sslmode=disable' -source file://D:/path-to-migrate -verbose up

# Migrate down
migrate -database 'postgres://kawaii:123456@0.0.0.0:5432/kawaii_db_test?sslmode=disable' -source file://D:/path-to-migrate -verbose down
```