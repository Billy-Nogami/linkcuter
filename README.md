# LinkCuter

HTTP‑сервис для сокращения ссылок. Поддерживает два хранилища: in‑memory и Postgres.

## Запуск

### Через Go

```bash
go run ./cmd/linkcuter
```

По умолчанию используется `configs/local.yaml`. Если файла нет, он будет создан.

### Через Docker

```bash
docker build -t linkcuter .
docker run --rm -p 8080:8080 linkcuter
```

### Через docker-compose (Postgres)

```bash
docker compose up --build
```

## API

### POST /api/shorten

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

Ответ:
```json
{
  "code": "AbCdEfGhIj",
  "short_url": "http://localhost:8080/AbCdEfGhIj"
}
```

### GET /{code}

```bash
curl -i http://localhost:8080/AbCdEfGhIj
```

## Конфигурация

Файл: `configs/local.yaml`

```yaml
server:
  addr: ":8080"
storage:
  mode: "memory"
  database_url: "postgres://postgres:postgres@localhost:5432/linkcuter?sslmode=disable"
```

Переменные окружения перекрывают файл:
- `CONFIG_PATH`
- `ADDR`
- `STORAGE`
- `DATABASE_URL`

## Миграции

Для Postgres миграции вшиты в бинарник и накатываются автоматически при старте.
