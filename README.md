# Сервис для работы с ПВЗ
[Исходное тестовое задание](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-spring-2025/Backend-trainee-assignment-spring-2025.md)

Backend-сервис для сотрудников ПВЗ, который позволяет вести учёт приемок товаров и добавленных позиций.  


## Запуск проекта
1. Клонируйте репозиторий
```bash
git clone https://github.com/KrivosheevNikita/avito-backend-task.git
cd avito-backend-task
```
2. Настройте переменные окружения
```bash
cp .env.example .env
```
3. Запустите через Docker
```bash
docker-compose up --build
```

- HTTP сервис будет доступен на: `localhost:8080`
- gRPC: `localhost:3000`
- Метрики: `localhost:9000`

## Функциональность
Сервис реализует данный [API](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-spring-2025/swagger.yaml).
- Заведение ПВЗ
- Добавление товаров в рамках одной приёмки
- Удаление товаров в рамках не закрытой приёмки
- Закрытие приёмки
- Получение списка ПВЗ и всю информацию по ним при помощи пагинации и фильтрации по дате

Дополнительно:
- Авторизация и регистрация пользователей
- gRPC-метод, возвращающий список всех ПВЗ
- Сбор метрик с помощью prometheus
- Логирование
## Стек технологий
- Go
- PostgreSQL
- gRPC
- Prometheus
- Docker


## Тестирование

Запуск unit и интеграционного тестов:
```bash
go test ./... -cover
```
Запуск через Docker:
```bash
docker exec -it pvz-service go test ./... -cover
```
Интеграционный тест находится в `tests/integration_test.go` и содержит:

1. Создание ПВЗ
2. Начало приемки
3. Добавление 50 товаров
4. Закрытие приемки
