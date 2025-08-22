# Демонстрационный сервис с Kafka, PostgreSQL, кешем

Система для приема заказов через Kafka, сохранения в PostgreSQL, кеширования и предоставления REST API для их получения.

## Структура

```
demo-service/
├── cmd/
│   ├── app/            # Точка входа сервиса (application)
│   └── sup/            # Точка входа продюсера (supplement)
├── configs/            # Файлы конфигураций
├── internal/
│   ├── app/            # Запуск сервиса
│   ├── config/         # Конфигурации
│   ├── model/          # Модели данных
│   ├── repository/     # Работа с кешем и postgres
│   ├── service/        # Бизнес-логика
│   ├── sup/            # Запуск продюсера
│   ├── transport/      # Работа с http и kafka
│   └── util/           # Вспомогательные инструмены
├── migrations/         # Файлы миграций
├── models/             # Файлы моделей
├── web/                # Файлы фронтенда
└── docker-compose.yml  # Файл для запуска kafka и postgres
```

## Запуск

Выполните в разных терминалах

```
# Запуск kafka и pospgres
docker compose run

# Запуск сервиса
go run ./cmd/app/main.go

# Запуск продюсера
go run ./cmd/sup/main.go
```

## Использование

### API Endpoints

Получение заказа по UID

```
GET /order/<order_uid>
```

### WEB Interface

UI для получения заказов

```
http://localhost:8081/

```

UI для мониторинга Kafka

```
http://localhost:8080/
```
