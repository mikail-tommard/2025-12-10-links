# link-checker-service

Сервис для проверки доступности интернет-ресурсов и генерации PDF-отчетов

## Описание

Веб-сервер в который пользователь может отправлять ссылки на интернет-ресурсы, как по одной ссылке так и сразу несколько.
Сервер отвечает пользователю доступен/недоступен тот или иной интернет-ресурс.

Так же пользователь может отправлять запрос со списком номеров ранее отправленных ссылок, а сервис должен будет вернуть PDF-файл с отчетом о статусе интернет ресурсов входящий в этот список.

## API

Базовый url-сервера: `http://localhost:8080`

## Проверка ссылок

Проверка ссылок происходит по url `http://localhost:8080/links`, метод `POST`.
Тело запроса: `json{"links_list": ["google.com", "malformedlink.gg"]}`

Пример ответ:
```json
{
  "next_id": 2,
  "batches": [
    {
      "ID": 1,
      "Links": [{ "URL": "https://go.dev" }, { "URL": "https://google.com" }],
      "Results": [
        {
          "Link": { "URL": "https://go.dev" },
          "Status": "unavailable",
          "Error": "Get \"https://go.dev\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)"
        },
        { "Link": { "URL": "https://google.com" }, "Status": "available", "Error": "" }
      ],
      "Status": "done",
      "CreatedAt": "2025-12-07T17:53:50.628485+03:00",
      "UpdatedAt": "2025-12-07T17:53:55.635446+03:00"
    }
  ]
}
```

## Генерация PDF-файла

Генерация PDF-файла ранее отправленных ссылок происходит по url `http://localhost:8080/report`, метод `POST`.
Тело запроса: `json{"links_num": [1]}`
После того как вы отправили запрос на данный url вам вернеться отчет в виде PDF-файла.

## Структура проекта
```text
.
├── cmd
│   └── linkschecker          # точка входа приложения
│       └── main.go           # запуск HTTP-сервера
├── data
│   └── state.json            # сохранённое состояние/бэтчи ссылок
└── internal
    ├── adapter               # слой адаптеров (внешние интерфейсы)
    │   ├── checker           # адаптер проверки ссылок
    │   │   └── checker.go
    │   ├── httpapi           # HTTP API (роуты, хэндлеры)
    │   │   └── server.go
    │   └── report            # генерация PDF-отчётов
    │       └── pdfreport.go
    ├── storage               # работа с файловым хранилищем
    │   └── filestorage.go
    ├── domain                # доменные сущности (модели)
    │   └── links.go
    └── usecase               # бизнес-логика (юзкейсы)
        ├── links.go          # сценарии проверки ссылок
        └── report.go         # сценарии формирования отчёта
```

## Запуск проекта

1. Клонируйте репозиторий:
  ```bash
    https://github.com/mikail-tommard/2025-12-10-links.git
  ```
2. Запуск:
  ```bash
    go run cmd/linkschecker/main.go
  ```