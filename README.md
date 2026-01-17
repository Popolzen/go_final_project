cat > README.md << 'EOF'
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
- Docker (опционально)

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
make generate-key  # Сгенерировать ключ шифрования
# Добавить ключ в .env
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

## Использование

### Регистрация
```bash
./bin/gophkeeper-client register alice MySecurePass123
```

### Логин
```bash
./bin/gophkeeper-client login alice MySecurePass123
```

### Добавление пароля
```bash
./bin/gophkeeper-client add login gmail --username alice@gmail.com --password qwerty123
```

### Просмотр паролей
```bash
./bin/gophkeeper-client list logins
./bin/gophkeeper-client get login gmail
```

### Синхронизация
```bash
./bin/gophkeeper-client sync
```

## Разработка

### Команды Make

- `make help` - показать все команды
- `make dev` - запустить окружение для разработки
- `make test` - запустить тесты
- `make build-all` - собрать под все платформы

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

## Лицензия

MIT
EOF