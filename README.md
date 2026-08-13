# Cosmos Social App

Учебный пет-проект: мини-соцсеть. Первая полноценная практика бэкенда на Go.

## Возможности
- Регистрация и аутентификация (JWT access + refresh, ротация refresh-токенов)
- Посты (создание, удаление, загрузка изображений)
- Лента, лайки, комментарии
- Профиль пользователя

## Стек

**Backend:** Go, PostgreSQL, SQLX, Redis  
**Frontend:** React, TypeScript, Tailwind CSS, Vite (сгенерирован ИИ)

## Backend
- Слои: `handler → service → repository`
- Миграции
- Unit of Work для транзакций
- JWT: access / refresh, whitelist refresh-токенов в Redis
- Логирование: `slog` + `request_id`
- Интеграционные тесты репозиториев (отдельная БД)
- Unit-тесты сервисного слоя

Бэкенд написан мной. Swagger и часть тестов — с помощью ИИ, код проверен вручную.  
Изначально был GORM, позже миграция на SQLX.

## Запуск через Docker
## Требования: 
- Docker 
- Docker Compose v2

Иметь в корне проекта .env файл(пример в .env.example)
Запуск контейнеров через Docker
docker compose up --build -d

Провести миграции
for f in backend/migrations/*.up.sql; do
  docker compose exec -T db psql -U aegis -d my_db < "$f"
done

фронтенд откроется на localhost:5173