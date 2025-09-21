# Togglr Go SDK

Go SDK для работы с Togglr - системой управления фича-флагами.

## Установка

```bash
go get github.com/rom8726/togglr-sdk-go
```

## Быстрый старт

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/rom8726/togglr-sdk-go"
)

func main() {
    // Создание клиента с настройками по умолчанию
    client, err := togglr.NewClientWithDefaults("your-api-key-here",
        togglr.WithBaseURL("http://localhost:8090"),
        togglr.WithTimeout(1*time.Second),
        togglr.WithCache(1000, 10*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Создание контекста запроса
    ctx := togglr.NewContext().
        WithUserID("user123").
        WithCountry("US").
        WithDeviceType("mobile")

    // Оценка фича-флага
    value, enabled, found, err := client.Evaluate("new_ui", ctx)
    if err != nil {
        log.Fatal(err)
    }

    if found {
        fmt.Printf("Feature enabled: %t, value: %s\n", enabled, value)
    }
}
```

## Конфигурация

### Создание клиента

```go
// С настройками по умолчанию
client, err := togglr.NewClientWithDefaults("api-key")

// С кастомной конфигурацией
cfg := togglr.DefaultConfig("api-key")
cfg.BaseURL = "https://api.togglr.com"
cfg.Timeout = 2 * time.Second
cfg.Retries = 3

client, err := togglr.NewClient(cfg)
```

### Функциональные опции

```go
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithBaseURL("https://api.togglr.com"),
    togglr.WithTimeout(2*time.Second),
    togglr.WithRetries(3),
    togglr.WithCache(1000, 10*time.Second),
    togglr.WithCircuitBreaker(true),
)
```

## Использование

### Создание контекста запроса

```go
ctx := togglr.NewContext().
    WithUserID("user123").
    WithUserEmail("user@example.com").
    WithCountry("US").
    WithDeviceType("mobile").
    WithOS("iOS").
    WithOSVersion("15.0").
    WithBrowser("Safari").
    WithLanguage("en-US")
```

### Оценка фича-флагов

```go
// Полная оценка
value, enabled, found, err := client.Evaluate("feature_key", ctx)

// Простая проверка включенности
isEnabled, err := client.IsEnabled("feature_key", ctx)

// С значением по умолчанию
defaultValue := client.EvaluateBoolOrDefault("feature_key", ctx, false)
```

### Работа с контекстом

```go
// С контекстом отмены
ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
defer cancel()

value, enabled, found, err := client.EvaluateWithContext(ctx, "api-key", "feature_key", reqCtx)
```

## Кэширование

SDK поддерживает опциональное кэширование результатов оценки:

```go
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithCache(1000, 10*time.Second), // размер кэша и TTL
)
```

## Повторные попытки

SDK автоматически повторяет запросы при временных ошибках:

```go
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithRetries(3), // количество попыток
    togglr.WithBackoff(togglr.Backoff{
        BaseDelay: 100 * time.Millisecond,
        MaxDelay:  2 * time.Second,
        Factor:    2.0,
    }),
)
```

## Логирование и метрики

```go
// Кастомный логгер
logger := &MyLogger{}
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithLogger(logger),
)

// Кастомные метрики
metrics := &MyMetrics{}
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithMetrics(metrics),
)
```

## Обработка ошибок

```go
value, enabled, found, err := client.Evaluate("feature_key", ctx)
if err != nil {
    switch {
    case errors.Is(err, togglr.ErrUnauthorized):
        // Ошибка авторизации
    case errors.Is(err, togglr.ErrBadRequest):
        // Неверный запрос
    default:
        // Другая ошибка
    }
}
```

## Генерация клиента

Для обновления сгенерированного клиента из OpenAPI спецификации:

```bash
make generate
```

## Сборка и тестирование

```bash
# Сборка
make build

# Тестирование
make test

# Линтинг
make lint

# Очистка
make clean
```

## Примеры

Полные примеры использования находятся в `pkg/togglr/examples/`.
