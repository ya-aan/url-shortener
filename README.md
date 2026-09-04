# URL Shortener

REST API сервис для сокращения URL на Go + PostgreSQL.

## Запуск

```bash
git clone https://github.com/ya-aan/url-shortener.git
cd url-shortener

cp .env.example .env

docker compose up --build
```

После запуска API доступен по адресу:

```text
http://localhost:8080
```

PostgreSQL доступен с хоста на порту:

```text
localhost:5433
```

## Тесты

```bash
go test ./...
```
