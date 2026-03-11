## Описание проекта

веб‑сервер‑планировщик задач (to‑do / scheduler).  
Он предоставляет HTTP API и простой фронтенд для:
- добавления задач с датой, заголовком, комментарием и правилом повторения;
- получения ближайших задач;
- просмотра и редактирования существующих задач;
- хранения данных в локальной SQLite‑базе (`scheduler.db`).

Фронтенд находится в директории `web`, а автотесты для проверки API — в директории `tests`.

## Задания со звёздочкой

Задания повышенной трудности выполнялись частичо .

## Запуск проекта локально

### Требования

- Установлен Go (версия не ниже `1.24`).

### Шаги запуска

1. Клонировать репозиторий и перейти в папку проекта:
   ```bash
   git clone <repo-url>
   cd go_final_project
   ```
2. Установить зависимости :
   ```bash
   go mod tidy
   ```
3. Запустить сервер:
   ```bash
   go run ./...
   ```

### Параметры окружения

Сервер использует две переменные окружения:

- `TODO_DBFILE` — путь к файлу базы данных SQLite.
  - По умолчанию: `scheduler.db` в корне проекта.
- `TODO_PORT` — порт HTTP‑сервера.
  - По умолчанию — порт из `tests/settings.go`:
    ```go
    var Port = 7540
    ```
    То есть, если переменная не задана, сервер слушает `:7540`.

```bash
TODO_DBFILE=./scheduler.db TODO_PORT=:7540 go run ./...
```

### Адрес в браузере

После запуска доступен по адресу:

- `http://localhost:7540/`

## Запуск тестов

Файл настроек тестов: `tests/settings.go`:

```go
var Port   = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search       = false
var Token        = ``
```

Для запуска тестов:

```bash
go test ./tests 
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestTask$ ./tests
go test -run ^TestEditTask$ ./tests
```

## Сборка и запуск через Docker

```bash
docker build -t scheduler-app .
```
### Запуск контейнера

```bash
docker run --rm \
  -p 7540:7540 \
  -v "$(pwd)/scheduler.db:/data/scheduler.db" \
  -e TODO_DBFILE=/data/scheduler.db \
  -e TODO_PORT=:7540 \
  --name scheduler-app \
  scheduler-app
```
