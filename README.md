# todo-app

Учебный pet-проект: REST API для управления задачами на Go, без веб-фреймворков — только стандартный `net/http` (роутинг с методами и path-параметрами из Go 1.22+) и `pgx` для PostgreSQL.

## Стек

- Go 1.25
- PostgreSQL 18
- [jackc/pgx v5](https://github.com/jackc/pgx) — драйвер и пул соединений
- [golang-migrate](https://github.com/golang-migrate/migrate) — миграции
- Docker / Docker Compose

## Структура

```
cmd/todo-app/          точка входа, сборка зависимостей, graceful shutdown
internal/
  core/
    app-error/         доменные ошибки
    db/                подключение к postgres
    domains/task/      доменная модель Task и её правила
  features/task/
    repository/        SQL-запросы, маппинг строк в модель
    service/           бизнес-логика поверх репозитория
    transport/         http-сервер, роуты, хендлеры, DTO
migrations/            SQL-миграции
```

Слои связаны через интерфейсы, объявленные на стороне потребителя: `transport.Service` и `service.Repository`. Зависимости собираются в `main`.

## Конфигурация

Скопировать пример и подставить свои значения:

```sh
cp .env.example .env
```

| Переменная | Назначение |
| --- | --- |
| `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` | учётные данные БД |
| `PGX_CONN_LOCAL` | DSN для запуска с хоста (`127.0.0.1:5432`) — используется миграциями из `Makefile` |
| `PGX_CONN_DOCKER` | DSN внутри compose-сети (хост `postgres`) — его читает приложение |

## Запуск

### Docker Compose (рекомендуемый способ)

```sh
docker compose up --build
```

Поднимаются три сервиса: `postgres` (с healthcheck), `migrate` (накатывает миграции и завершается), `app` (стартует после успешных миграций). API доступен на `http://localhost:9091`, данные БД лежат в `./out/pgdata`.

### Локально

```sh
make run          # go run cmd/todo-app/main.go
make migrate-up   # накатить миграции  (нужен бинарник migrate)
make migrate-down # откатить
make migrate-version
```

`Makefile` подхватывает переменные из `.env`. Приложение читает DSN из `PGX_CONN_DOCKER`, поэтому при запуске с хоста в эту переменную нужно положить локальный DSN (`127.0.0.1:5432`).

## API

Базовый URL: `http://localhost:9091`

| Метод | Путь | Описание |
| --- | --- | --- |
| `POST` | `/tasks` | создать задачу |
| `GET` | `/tasks` | список задач, фильтр `?completed=true` / `?completed=false` |
| `GET` | `/tasks/{id}` | задача по ID |
| `PATCH` | `/tasks/{id}` | отметить выполненной / снять отметку |
| `PUT` | `/tasks/{id}` | обновить заголовок и описание |
| `DELETE` | `/tasks/{id}` | удалить задачу |

### POST /tasks

`title` и `description` обязательны — пустые значения отклоняются доменной моделью (`400`).

```sh
curl -X POST localhost:9091/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title": "купить хлеб", "description": "по дороге домой"}'
```

`201 Created`:

```json
{
  "ID": 1,
  "Title": "купить хлеб",
  "Description": "по дороге домой",
  "Completed": false,
  "CreatedAt": "2026-08-03T12:00:00Z",
  "CompletedAt": null
}
```

### GET /tasks

```sh
curl localhost:9091/tasks
curl 'localhost:9091/tasks?completed=true'
curl 'localhost:9091/tasks?completed=false'
```

Без параметра возвращаются все задачи; некорректное значение `completed` игнорируется.

### GET /tasks/{id}

```sh
curl localhost:9091/tasks/1
```

`404`, если задачи нет.

### PATCH /tasks/{id}

```sh
curl -X PATCH localhost:9091/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{"completed": true}'
```

При `true` проставляется `completed_at`, при `false` — сбрасывается в `null`. В ответе — обновлённая задача.

### PUT /tasks/{id}

```sh
curl -X PUT localhost:9091/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{"title": "новый заголовок", "description": "новое описание"}'
```

### DELETE /tasks/{id}

```sh
curl -X DELETE localhost:9091/tasks/1
```

`204 No Content`.

## Схема данных

```sql
CREATE TABLE tasks (
    id           INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title        VARCHAR(200)  NOT NULL,
    description  VARCHAR(2000) NOT NULL,
    completed    BOOLEAN       NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);
```

## TODO

- [ ] структурированный логгер вместо `fmt.Println`
- [ ] json-теги в DTO ответов (сейчас поля отдаются с именами полей Go)
- [ ] маппинг доменных ошибок в HTTP-коды в одном месте
- [ ] валидация тела запроса на уровне transport
- [ ] тесты (unit для сервиса, интеграционные для репозитория)
- [ ] конфиг вместо прямых `os.Getenv`, адрес сервера из переменной окружения
