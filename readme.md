# Shorty — сервис сокращения ссылок

Shorty — это полнофункциональный сервис по сокращению URL-адресов, созданный как учебный и практический проект.  
Он позволяет:

- создавать короткие ссылки;
- просматривать статистику переходов;
- вести историю всех созданных ссылок;
- автоматически увеличивать количество кликов;
- запускать весь проект через Docker (Backend + Frontend + PostgreSQL).

Проект демонстрирует навыки backend-разработки на Go, работы с PostgreSQL, Docker, а также клиентской части на React.

---

##  Стек технологий

### **Backend**
- Go 1.22+
- Gin (HTTP framework)
- pgx/v5 (PostgreSQL драйвер)
- Zerolog (структурированное логирование)
- Clean Architecture (handlers → services → repositories)
- Docker / Docker Compose

### **Frontend**
- React + TypeScript
- Vite
- TailwindCSS

### **Database**
- PostgreSQL 15

---
### Поток данных

1. Frontend отправляет запрос на `POST /api/v1/shorten`.
2. Handler вызывает слой Service.
3. Service генерирует короткий код и вызывает Repository.
4. Repository сохраняет запись в PostgreSQL.
5. При открытии `/{shortcode}` backend делает redirect и увеличивает счётчик.
6. История ссылок загружается через `GET /api/v1/urls`.

---

## 🗄️ API эндпоинты

### **Создать короткую ссылку**
`POST /api/v1/shorten`
```json
{
  "url": "https://example.com"
}

Получить статистику

GET /api/v1/stats/:code

Получить историю всех ссылок

GET /api/v1/urls

Редирект по короткому коду

GET /:code

---

Запуск проекта через Docker

Убедитесь, что Docker и Docker Compose установлены.

1️⃣ Клонировать репозиторий
2️⃣ Запустить проект
docker-compose up -d --build
3️⃣ Открыть приложение
Frontend: http://localhost:3000
Backend:  http://localhost:8080
PostgreSQL: localhost:5432

---

## Очистка базы данных (очистить историю ссылок)

Если нужно полностью очистить историю:

1. Остановить контейнеры
- docker-compose down
2. Удалить volume PostgreSQL
- docker volume rm shorty_pgdata
3. Запустить заново
-docker-compose up -d --build

История ссылок будет полностью пустой.

## Возможные улучшения
    •	Авторизация пользователей
    •	Срок жизни ссылок
    •	Ограничение API rate limit
    •	Поддержка кастомных short-кодов
    •	Подробная логистика кликов (IP, User-Agent)
    •	Админ-панель