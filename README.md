# 📝 Микросервисное TODO-приложение

## Описание

Этот проект реализует микросервисную архитектуру для TODO-приложения с авторизацией по JWT.  
Состав системы:
- **api-gateway**: маршрутизация и аутентификация
- **auth**: сервис аутентификации пользователей (работает со своей БД)
- **todo**: сервис списков и задач пользователей (работает со своей БД)
- **PostgreSQL**: отдельные БД для auth и todo

Docker Compose собирает и запускает инфраструктуру в едином пространстве.  
Все сервисы написаны на Go.

---

## Структура проекта

```
.
├── api-gateway          # API Gateway (reverse proxy, auth middleware)
├── auth                 # Сервис аутентификации (своя БД)
├── todo                 # Сервис задач (своя БД)
├── compose.yml          # docker-compose для локального запуска
└── .env                 # переменные окружения (опции БД)
```

---

## Быстрый старт

1. **Создайте .env с настройками БД (пример):**
    ```ini
    # .env
    AUTH_DATASOURCE_USERNAME=authuser
    AUTH_DATASOURCE_PASSWORD=authpass
    AUTH_DATASOURCE_DB_NAME=authdb
    AUTH_DATASOURCE_PORT=5433

    TODO_DATASOURCE_USERNAME=todouser
    TODO_DATASOURCE_PASSWORD=todopass
    TODO_DATASOURCE_DB_NAME=tododb
    TODO_DATASOURCE_PORT=5434
    ```
    Если файла нет — значения возьмутся по умолчанию из config.yml каждого сервиса.

2. **Проверьте наличие файлов `configs/config.yml` для каждого сервиса.**
    - Пример `todo/configs/config.yml`:
        ```yaml
        postgres:
          host: db_todo
          port: 5432
          dbname: tododb
          username: todouser
          password: todopass
        ```

3. **Запустите в корне проекта:**
    ```bash
    docker compose up --build
    ```
    Это собирает проект, запускает базы, сервисы и накатывает миграции, если используется [migrate](https://github.com/golang-migrate/migrate).

---

## Миграции

- **Хранятся в директориях `auth/schema` и `todo/schema`.**
- Применять миграции можно вручную через docker (если не используется авто-миграция через compose):
    ```bash
    docker run --rm --network=host \\
        -v $(pwd)/todo/schema:/migrations \\
        migrate/migrate:v4.17.1 \\
        -path=/migrations \\
        -database \"postgres://todouser:todopass@localhost:5434/tododb?sslmode=disable\" up

    docker run --rm --network=host \\
        -v $(pwd)/auth/schema:/migrations \\
        migrate/migrate:v4.17.1 \\
        -path=/migrations \\
        -database \"postgres://authuser:authpass@localhost:5433/authdb?sslmode=disable\" up
    ```
- **CLI для создания миграций** (необходим [migrate](https://github.com/golang-migrate/migrate#usage)):
    ```bash
    migrate create -ext sql -dir todo/schema create_table_todos
    ```

---

## Использование

- **Api-gateway доступен на:** `http://localhost:8080`
- **Auth-сервис:**        `http://localhost:8081` (если проброшен порт)
- **TODO-сервис:**        `http://localhost:8082` (если проброшен порт)
- **Postgres для auth:**  `localhost:5433`
- **Postgres для todo:**  `localhost:5434`

### Пример авторизации (через gateway):

```http
POST /auth/sign-in
Content-Type: application/json

{
  "username": "demo",
  "password": "demo"
}
```

---

## Остановка и очистка

```bash
docker compose down -v
```
Флаг `-v` полностью удаляет тома с базой (полная очистка БД).

---

## Полезные команды

- **Посмотреть логи:**
    ```bash
    docker compose logs -f
    ```

- **Войти в контейнер:**
    ```bash
    docker compose exec todo sh
    ```

- **Провести миграцию одной БД вручную** (пример см. выше).

---

## Разработка

- Для быстрого перезапуска сервисов используйте volumes для `configs`:
    ```yaml
    volumes:
      - ./todo/configs:/app/configs
    ```
- Все переменные окружения можно указать через `.env` или напрямую в compose.yml

---
