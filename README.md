# Cosmos Social App

Учебный pet-project: мини-соцсеть. Первая полноценная практика бэкенда на Go.

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

## Запуск
1. Поднять PostgreSQL и Redis  
2. Заполнить `.env` (пример — `.env.example`)  
3. Применить миграции  
4. `go run .` в `backend/`  
5. Frontend: `npm install && npm run dev` в `frontend/`

## Планы
- Rate limit на регистрацию
- Обновление Swagger
- Доработка/чистка части эндпоинтов
- (далее) обновление регистрации