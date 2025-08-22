# Демонстрационный сервис с Kafka, PostgreSQL и кешем

Система из двух сервисов. Основной сервис (application) для приема заказов через Kafka, сохранения в PostgreSQL, кеширования и предоставления REST API и web-интерфейса для их получения. Вспомогательный сервис (supplement) для генерации и отправки заказов в Kafka.

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

## Управление

### Запуск

Выполните в разных терминалах:

Старт kafka и pospgres

```
docker compose run
```

Старт application

```
go run ./cmd/app/main.go
```

Старт supplement

```
go run ./cmd/sup/main.go
```

### Остановка

Нажмите Ctrl + C во всех терминалах

## Использование

### API Endpoints

Получение заказа по UID

```
GET /order/<order_uid>
```

### WEB Interfaces

UI для получения заказов

```
http://localhost:8081/
```

UI для мониторинга Kafka

```
http://localhost:8080/
```
