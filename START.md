# Як запустити бекенд

## Проблема: Docker daemon не запущений

Якщо бекенд не відповідає, перевірте чи запущений Docker:

### 1. Запустіть Docker/OrbStack

- Відкрийте **OrbStack** або **Docker Desktop**
- Дочекайтесь поки Docker повністю запуститься

### 2. Перевірте чи Docker працює

```bash
docker ps
```

Якщо команда працює без помилок - Docker запущений.

### 3. Запустіть бекенд

```bash
cd backend
make docker-up
```

Або вручну:

```bash
cd backend
docker-compose up -d
```

### 4. Перевірте статус контейнерів

```bash
cd backend
docker-compose ps
```

Має бути два контейнери:
- `stempo_postgres` - база даних
- `stempo_backend` - API сервер

### 5. Перевірте логи якщо щось не працює

```bash
cd backend
docker-compose logs backend
```

### 6. Запустіть міграції та сідери

```bash
cd backend
make docker-migrate
make docker-seed
```

### 7. Перевірте чи API працює

```bash
curl http://192.168.178.111:3000/api/businesses/featured
curl http://192.168.178.111:3000/api/categories
```

Якщо все працює, ви побачите JSON відповіді.

## Альтернатива: Запуск без Docker

Якщо Docker не працює, можна запустити локально:

1. Встановіть PostgreSQL локально
2. Створіть `.env` файл з налаштуваннями БД
3. Запустіть:
   ```bash
   cd backend
   make migrate
   make seed
   make run
   ```
