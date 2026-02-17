# Stempo Backend

Backend API для Stempo застосунку, написаний на Go з використанням чистої архітектури (Clean Architecture).

## Архітектура

Проєкт використовує **Clean Architecture** з чітким розділенням на шари:

```
backend/
├── cmd/
│   └── server/                    # Точка входу додатку
├── internal/                      # Внутрішній код (не експортується)
│   ├── domain/                    # Доменний шар
│   │   ├── entity/                # Бізнес-сутності (entities)
│   │   └── repository/            # Інтерфейси репозиторіїв
│   ├── usecase/                   # Шар бізнес-логіки
│   ├── repository/                # Реалізація репозиторіїв
│   │   └── postgres/              # PostgreSQL репозиторії
│   ├── delivery/                  # Шар доставки (presentation)
│   │   └── http/                  # HTTP handlers (REST API)
│   └── infrastructure/            # Інфраструктурний шар
│       ├── database/              # Підключення до БД
│       ├── middleware/            # HTTP middleware
│       └── config/                # Конфігурація
└── pkg/                           # Публічні пакети (можуть використовуватися іншими проєктами)
    └── utils/                     # Утиліти
```

### Принципи Clean Architecture:

1. **Domain Layer** (`internal/domain/`)
   - Незалежний від зовнішніх залежностей
   - Містить бізнес-сутності та інтерфейси
   - Не залежить від фреймворків або БД

2. **Use Case Layer** (`internal/usecase/`)
   - Реалізує бізнес-логіку
   - Залежить тільки від domain layer
   - Використовує інтерфейси з domain layer

3. **Repository Layer** (`internal/repository/`)
   - Реалізує інтерфейси з domain layer
   - Інкапсулює логіку роботи з БД
   - Може бути замінено без змін в usecase

4. **Delivery Layer** (`internal/delivery/`)
   - HTTP handlers для REST API
   - Залежить від usecase layer
   - Обробляє HTTP запити та відповіді

5. **Infrastructure Layer** (`internal/infrastructure/`)
   - Технічні деталі (БД, middleware, конфігурація)
   - Може бути замінено без змін в інших шарах

6. **Public Packages** (`pkg/`)
   - Утиліти, які можуть використовуватися іншими проєктами
   - Стабільний API

## Запуск

### З Docker Compose (Production)

```bash
docker-compose up -d
```

### Локально (Development)

1. Створіть `.env` файл на основі `.env.example`
2. Запустіть PostgreSQL:
```bash
docker-compose -f docker-compose.dev.yml up -d postgres
```
3. Запустіть міграції:
```bash
make migrate
```
або
```bash
go run ./cmd/server migrate
```
4. Запустіть сервер:
```bash
make run
```
або
```bash
go run ./cmd/server
```

### Make команди

- `make build` - Зібрати бінарний файл
- `make run` - Запустити сервер локально
- `make migrate` - Запустити міграції БД
- `make docker-up` - Запустити Docker Compose
- `make docker-down` - Зупинити Docker Compose
- `make docker-build` - Зібрати Docker образи
- `make docker-logs` - Переглянути логи backend
- `make test` - Запустити тести
- `make test-coverage` - Запустити тести з покриттям коду

## Тестування

Проєкт містить unit тести для:
- Use cases (бізнес-логіка)
- HTTP handlers
- Repositories
- Middleware

Запуск тестів:
```bash
make test
```

Запуск тестів з покриттям:
```bash
make test-coverage
```

## API Endpoints

### Auth
- `POST /api/auth/register` - Реєстрація
- `POST /api/auth/login` - Вхід
- `POST /api/auth/logout` - Вихід
- `GET /api/auth/me` - Поточний користувач

### Businesses
- `GET /api/businesses` - Список бізнесів
- `GET /api/businesses/featured` - Рекомендовані бізнеси
- `GET /api/businesses/:id` - Отримати бізнес
- `POST /api/businesses` - Створити бізнес
- `PUT /api/businesses/:id` - Оновити бізнес
- `DELETE /api/businesses/:id` - Видалити бізнес

### Bonus Programs
- `GET /api/bonus-programs` - Список програм
- `GET /api/bonus-programs/my-business?businessId=:id` - Програми мого бізнесу
- `GET /api/bonus-programs/:id` - Отримати програму
- `POST /api/bonus-programs` - Створити програму
- `PUT /api/bonus-programs/:id` - Оновити програму
- `DELETE /api/bonus-programs/:id` - Видалити програму

### Loyalty Cards
- `GET /api/loyalty-cards` - Список карток користувача
- `GET /api/loyalty-cards/:id` - Отримати картку
- `POST /api/loyalty-cards` - Створити картку
- `POST /api/loyalty-cards/:id/stamps` - Додати штамп
- `PUT /api/loyalty-cards/:id` - Оновити картку

## Технології

- **Go 1.21+**
- **Gin** - HTTP web framework
- **GORM** - ORM для роботи з PostgreSQL
- **JWT** - Авторизація
- **PostgreSQL 16** - База даних
- **Docker & Docker Compose** - Контейнеризація
- **Testify** - Testing framework

## Залежності між шарами

```
cmd/server
    ↓
internal/delivery/http
    ↓
internal/usecase
    ↓
internal/domain (interfaces)
    ↑
internal/repository/postgres (implements)
    ↓
internal/infrastructure/database
```

**Правило**: Залежності завжди спрямовані всередину (inward). Зовнішні шари залежать від внутрішніх, але не навпаки.
