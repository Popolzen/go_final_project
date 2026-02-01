# GophKeeper

Безопасный менеджер паролей с клиент-серверной архитектурой.

## Возможности

- 🔐 Хранение логинов/паролей, текстовых заметок, бинарных данных, банковских карт
- 🔒 End-to-end шифрование (AES-256-GCM)
- 🔄 Синхронизация между устройствами
- 💻 CLI клиент для Windows, Linux, macOS
- 🚀 Быстрый и безопасный

## Быстрый старт

### Требования

- Go 1.21+
- PostgreSQL 12+

### Установка

1. Клонировать репозиторий:
```bash
git clone https://github.com/yourusername/gophkeeper.git
cd gophkeeper
```

2. Запустить PostgreSQL:
```bash
make docker-up
```

3. Создать .env файл:
```bash
cp .env.example .env
```

4. Применить миграции:
```bash
make migrate-up
```

5. Запустить сервер:
```bash
make run-server
```

6. Собрать клиент:
```bash
make build-client
```

## Использование клиента

### Регистрация
```bash
./bin/gophkeeper-client register alice MySecurePass123
```

### Логин
```bash
./bin/gophkeeper-client login alice MySecurePass123
```

### Добавление данных

Логин/пароль:
```bash
./bin/gophkeeper-client add login gmail --login alice@gmail.com --password qwerty123
```

Текстовая заметка:
```bash
./bin/gophkeeper-client add text notes "My secret note"
```

Бинарный файл:
```bash
./bin/gophkeeper-client add binary docs --file /path/to/file.pdf
```

Банковская карта:
```bash
./bin/gophkeeper-client add card mycard --number 1234567812345678 --holder "ALICE SMITH" --cvv 123 --expiry 12/25
```

### Просмотр данных

Список всех секретов:
```bash
./bin/gophkeeper-client list
```

Получить конкретный секрет:
```bash
./bin/gophkeeper-client get login gmail
./bin/gophkeeper-client get text notes
./bin/gophkeeper-client get card mycard
```

### Обновление данных

```bash
./bin/gophkeeper-client update login  gmail --login new@gmail.com --password newpass
./bin/gophkeeper-client update text  notes "Updated note"
./bin/gophkeeper-client update card  mycard --number 9999888877776666 --holder "ALICE SMITH" --cvv 321 --expiry 01/27
```

### Удаление данных

```bash
./bin/gophkeeper-client delete login gmail
./bin/gophkeeper-client delete text notes
```

## Команды Make

- `make help` - показать все команды
- `make build-server` - собрать сервер
- `make build-client` - собрать клиент
- `make run-server` - запустить сервер
- `make migrate-up` - применить миграции
- `make migrate-down` - откатить миграции
- `make docker-up` - запустить PostgreSQL в Docker
- `make docker-down` - остановить PostgreSQL
- `make test` - запустить тесты
- `make clean` - удалить бинарники

## Архитектура

```
┌─────────┐      HTTPS/TLS      ┌─────────┐
│ Client  │ ◄─────────────────► │ Server  │
│  (CLI)  │   JWT + Encrypted   │  (API)  │
└─────────┘                     └────┬────┘
                                      │
                                      ▼
                                 ┌─────────┐
                                 │   DB    │
                                 │(PostgreSQL)│
                                 └─────────┘
```

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - регистрация
- `POST /api/v1/auth/login` - вход

### Секреты (требуется JWT)
- `POST /api/v1/secrets/` - создать секрет
- `GET /api/v1/secrets/` - список секретов
- `GET /api/v1/secrets/{id}` - получить по ID
- `GET /api/v1/secrets/by-name/{name}?type=login` - получить по имени
- `PUT /api/v1/secrets/{id}` - обновить секрет
- `DELETE /api/v1/secrets/{id}` - удалить секрет
- `GET /api/v1/secrets/changes?after=2024-01-01T00:00:00Z` - изменения после даты

## Лицензия

MIT