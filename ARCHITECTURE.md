# Clean Architecture в Stempo Backend

## Огляд

Цей проєкт використовує **Clean Architecture** (також відома як Hexagonal Architecture або Ports & Adapters), яка забезпечує:

- **Незалежність від фреймворків**: Бізнес-логіка не залежить від зовнішніх бібліотек
- **Тестованість**: Бізнес-логіку можна тестувати без UI, БД, веб-сервера
- **Незалежність від UI**: UI можна легко замінити без змін в бізнес-логіці
- **Незалежність від БД**: Можна замінити PostgreSQL на MongoDB без змін в use cases
- **Незалежність від зовнішніх агентів**: Бізнес-логіка не знає нічого про зовнішній світ

## Структура шарів

### 1. Domain Layer (`internal/domain/`)

**Найвнутрішній шар** - не залежить ні від чого.

```
internal/domain/
├── entity/              # Бізнес-сутності
│   ├── user.go
│   ├── business.go
│   ├── bonus_program.go
│   └── loyalty_card.go
└── repository/          # Інтерфейси репозиторіїв (Ports)
    ├── user_repository.go
    ├── business_repository.go
    ├── bonus_program_repository.go
    └── loyalty_card_repository.go
```

**Характеристики:**
- Містить тільки бізнес-логіку та правила
- Не залежить від зовнішніх бібліотек
- Може містити тільки інтерфейси, не реалізації

### 2. Use Case Layer (`internal/usecase/`)

**Шар бізнес-логіки** - залежить тільки від domain layer.

```
internal/usecase/
├── auth_usecase.go
├── business_usecase.go
├── bonus_program_usecase.go
└── loyalty_card_usecase.go
```

**Характеристики:**
- Реалізує конкретні use cases додатку
- Використовує інтерфейси з domain layer
- Не знає про HTTP, БД, файлову систему
- Може бути протестований з mock репозиторіями

### 3. Repository Layer (`internal/repository/`)

**Шар доступу до даних** - реалізує інтерфейси з domain layer.

```
internal/repository/
└── postgres/
    ├── user_repository.go
    ├── business_repository.go
    ├── bonus_program_repository.go
    └── loyalty_card_repository.go
```

**Характеристики:**
- Реалізує інтерфейси з `internal/domain/repository/`
- Інкапсулює логіку роботи з БД
- Може бути замінено на іншу реалізацію (MongoDB, файли, тощо)
- Залежить від infrastructure layer для підключення до БД

### 4. Delivery Layer (`internal/delivery/`)

**Шар представлення** - залежить від use case layer.

```
internal/delivery/
└── http/
    ├── router.go
    ├── auth_handler.go
    ├── business_handler.go
    ├── bonus_program_handler.go
    └── loyalty_card_handler.go
```

**Характеристики:**
- Обробляє HTTP запити
- Конвертує HTTP запити в виклики use cases
- Конвертує результати use cases в HTTP відповіді
- Може бути замінено на gRPC, GraphQL, тощо

### 5. Infrastructure Layer (`internal/infrastructure/`)

**Технічний шар** - реалізує технічні деталі.

```
internal/infrastructure/
├── database/
│   └── postgres.go      # Підключення до БД
├── middleware/
│   └── auth.go          # JWT middleware
└── config/
    └── config.go        # Конфігурація
```

**Характеристики:**
- Містить технічні деталі реалізації
- Підключення до БД, файлової системи, зовнішніх API
- Може бути замінено без змін в інших шарах

### 6. Public Packages (`pkg/`)

**Публічні утиліти** - можуть використовуватися іншими проєктами.

```
pkg/
└── utils/
    ├── response.go      # Стандартні HTTP відповіді
    └── validator.go     # Валідація даних
```

**Характеристики:**
- Стабільний API
- Може використовуватися іншими проєктами
- Не залежить від internal пакетів

## Потік залежностей

```
┌─────────────────────────────────────────┐
│           cmd/server                    │  ← Entry point
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│      internal/delivery/http            │  ← HTTP handlers
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│         internal/usecase                │  ← Business logic
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│      internal/domain                    │  ← Entities & Interfaces
└─────────────────┬───────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌──────────────┐   ┌────────────────────┐
│  repository  │   │  infrastructure    │
│  (implements)│   │  (technical)       │
└──────────────┘   └────────────────────┘
```

## Правила залежностей

1. **Залежності спрямовані всередину**: Зовнішні шари залежать від внутрішніх
2. **Domain незалежний**: `internal/domain` не залежить ні від чого
3. **Use cases залежать тільки від domain**: Використовують тільки інтерфейси
4. **Repositories реалізують domain інтерфейси**: Але залежать від infrastructure
5. **Delivery залежить від use cases**: Але не знає про репозиторії напряму

## Приклад потоку даних

### Створення бізнесу:

```
HTTP Request
    ↓
internal/delivery/http/business_handler.go
    ↓ (викликає use case)
internal/usecase/business_usecase.go
    ↓ (використовує інтерфейс)
internal/domain/repository/business_repository.go (interface)
    ↑ (реалізує)
internal/repository/postgres/business_repository.go
    ↓ (використовує)
internal/infrastructure/database/postgres.go
    ↓
PostgreSQL Database
```

## Переваги такої архітектури

1. **Легке тестування**: Use cases можна тестувати з mock репозиторіями
2. **Гнучкість**: Можна замінити БД, HTTP framework без змін в бізнес-логіці
3. **Розуміння**: Чітке розділення відповідальностей
4. **Масштабованість**: Легко додавати нові use cases та handlers
5. **Підтримка**: Зміни в одному шарі не впливають на інші

## Приклади заміни компонентів

### Заміна БД:
- Створити `internal/repository/mongodb/`
- Реалізувати інтерфейси з `internal/domain/repository/`
- Змінити тільки `cmd/server/main.go` для використання нової реалізації
- **Use cases залишаються незмінними**

### Заміна HTTP на gRPC:
- Створити `internal/delivery/grpc/`
- Реалізувати gRPC handlers, які викликають use cases
- Змінити тільки `cmd/server/main.go`
- **Use cases залишаються незмінними**

### Додавання GraphQL:
- Створити `internal/delivery/graphql/`
- Реалізувати GraphQL resolvers, які викликають use cases
- **Use cases залишаються незмінними**
