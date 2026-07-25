# Space Posts — Frontend

Современный React-фронтенд для твоего бэкенда.

## Стек

- **React 18** + TypeScript
- **Vite**
- **React Router v6**
- **Axios** (с авто-refresh токена)
- **Tailwind CSS** (космический тёмный дизайн)

## Структура

```
src/
├── api/           # Все запросы к бэкенду
│   ├── client.ts  # axios instance + interceptors + error helper
│   ├── auth.ts
│   ├── user.ts
│   └── post.ts
├── components/    # UI-компоненты
│   ├── Navbar.tsx
│   ├── PostCard.tsx
│   ├── Stars.tsx      # анимированный фон
│   ├── Spinner.tsx
│   └── ProtectedRoute.tsx
├── context/
│   └── AuthContext.tsx   # логин / регистрация / токен / текущий юзер
├── pages/
│   ├── Home.tsx          # заглушка (лента потом)
│   ├── Login.tsx
│   ├── Register.tsx
│   ├── CreatePost.tsx
│   └── Profile.tsx       # /profile/:username
├── App.tsx
├── main.tsx
└── index.css
```

## Запуск

```bash
cd space-posts-frontend
npm install
npm run dev
```

Откроется на http://localhost:5173

Бэкенд должен быть на http://localhost:8080 (CORS уже включён).

## Что умеет

- Регистрация / Логин
- Автоматический refresh токена при 401
- Создание постов
- Просмотр профиля (своего и чужого) + список постов
- Лайк / анлайк (одна кнопка, `is_liked` учитывается)
- Космический дизайн: тёмный фон, звёзды, glassmorphism, градиенты

## Что пока нет (по твоей просьбе)

- Общая лента всех постов (сделаем позже)

## API Base URL

Жёстко прописан в `src/api/client.ts`:

```ts
const API_BASE = 'http://localhost:8080'
```

Если нужно поменять — только там.
