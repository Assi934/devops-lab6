# devops-lab6

Микросервис «Автомобили» (вариант 11) на Go kit: сущность `Car` (ID, Brand, Model, Year, Price, Sold), методы `Create`, `GetByID`, `MarkSold`.

Транспорт: HTTP (`:8080`) и gRPC (`:9090`), in-memory репозиторий.

![CI](https://github.com/Assi934/devops-lab6/actions/workflows/ci.yml/badge.svg)
![Docker](https://github.com/Assi934/devops-lab6/actions/workflows/docker.yml/badge.svg)

## Запуск

```sh
docker run -p 8080:8080 -p 9090:9090 ghcr.io/assi934/devops-lab6:main
```

## API (HTTP)

```sh
# Создать автомобиль
curl -X POST http://localhost:8080/cars \
  -H "Content-Type: application/json" \
  -d '{"car":{"brand":"Toyota","model":"Camry","year":2023,"price":2500000}}'

# Получить по ID
curl http://localhost:8080/cars/<id>

# Отметить проданным
curl -X PATCH http://localhost:8080/cars/<id>/sold
```

## gRPC

Схема: `pb/car.proto`, сервис `car.CarService`. Проверка:

```sh
grpcurl -plaintext -d '{"car":{"brand":"BMW","model":"X5","year":2024,"price":8500000}}' \
  localhost:9090 car.CarService/Create
```

## CI/CD

- `.github/workflows/ci.yml`: матрица Go 1.21/1.22, golangci-lint, тесты с `-race` и покрытием, артефакт `coverage.out`, кэш модулей.
- `.github/workflows/docker.yml`: multi-stage сборка образа и публикация в GHCR (`ghcr.io/assi934/devops-lab6`), теги `main`, `sha-<short>`, `v*` → семвер.
