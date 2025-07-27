# JWT Token Setup Guide

## Обзор

В проекте реализована система JWT аутентификации с поддержкой access и refresh токенов для безопасной работы с API.

## Архитектура

### Backend (Go)

#### Структура токенов:
- **Access Token**: Короткоживущий токен (15 минут по умолчанию) для доступа к защищенным ресурсам
- **Refresh Token**: Долгоживущий токен (7 дней по умолчанию) для обновления access токена

#### Основные компоненты:

1. **Middleware** (`internal/middleware/jwt.go`):
   - `JWTAuth()` - middleware для проверки access токенов
   - `GenerateAccessToken()` - генерация access токена
   - `GenerateRefreshToken()` - генерация refresh токена
   - `ParseToken()` - парсинг и валидация токенов

2. **Handlers**:
   - `LoginHandler` - аутентификация и выдача токенов
   - `RegisterHandler` - регистрация и выдача токенов
   - `RefreshTokenHandler` - обновление access токена
   - `ProfileHandler` - работа с профилем (защищенный роут)

3. **Configuration** (`internal/config/jwt.go`):
   - Загрузка настроек из переменных окружения
   - Настройка времени жизни токенов
   - Настройка секретного ключа

### Frontend (Flutter)

#### ApiService (`lib/services/api_service.dart`):
- Автоматическое управление токенами
- Автоматическое обновление access токена при истечении
- Безопасное хранение токенов в `flutter_secure_storage`

## Настройка

### Переменные окружения

Создайте файл `.env` в корне backend/user-service:

```env
# JWT Configuration
JWT_SECRET_KEY=your-super-secret-key-here
JWT_ACCESS_EXPIRY_MINUTES=15
JWT_REFRESH_EXPIRY_DAYS=7
```

### Значения по умолчанию:
- `JWT_SECRET_KEY`: "OMGMYKEY" (⚠️ измените в продакшене!)
- `JWT_ACCESS_EXPIRY_MINUTES`: 15
- `JWT_REFRESH_EXPIRY_DAYS`: 7

## API Endpoints

### Аутентификация

#### POST /api/users/register
Регистрация нового пользователя.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "display_name": "John Doe"
}
```

**Response:**
```json
{
  "user_id": 1,
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

#### POST /api/users/login
Вход в систему.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

#### POST /api/users/refresh
Обновление access токена.

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

### Защищенные ресурсы

#### GET /api/users/profile
Получение профиля пользователя.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "display_name": "John Doe"
}
```

#### PUT /api/users/profile
Обновление профиля пользователя.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request:**
```json
{
  "display_name": "New Name"
}
```

## Использование в Flutter

### Автоматическое управление токенами

ApiService автоматически:
1. Сохраняет токены при входе/регистрации
2. Добавляет access токен к запросам
3. Обновляет access токен при получении 401 ошибки
4. Очищает токены при logout

### Пример использования:

```dart
final apiService = ApiService();

// Вход в систему
final loginResult = await apiService.login('user@example.com', 'password123');
if (loginResult != null) {
  print('Успешный вход!');
}

// Получение профиля (токен добавляется автоматически)
final profile = await apiService.getProfile();
if (profile != null) {
  print('Профиль: ${profile['display_name']}');
}

// Выход из системы
await apiService.logout();
```

## Безопасность

### Рекомендации для продакшена:

1. **Секретный ключ**: Используйте криптографически стойкий секретный ключ
2. **HTTPS**: Всегда используйте HTTPS в продакшене
3. **Время жизни токенов**: Настройте подходящее время жизни для вашего приложения
4. **Refresh Token Rotation**: Рассмотрите возможность ротации refresh токенов
5. **Blacklisting**: Реализуйте blacklist для отозванных токенов

### Текущие настройки безопасности:

- ✅ Access токены имеют короткое время жизни (15 минут)
- ✅ Refresh токены имеют ограниченное время жизни (7 дней)
- ✅ Токены содержат минимально необходимую информацию
- ✅ Используется HMAC-SHA256 для подписи
- ✅ Валидация типа токена (access/refresh)

## Тестирование

### Тестирование с curl:

```bash
# Регистрация
curl -X POST http://localhost:8081/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","display_name":"Test User"}'

# Вход
curl -X POST http://localhost:8081/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Получение профиля (замените <access_token> на полученный токен)
curl -X GET http://localhost:8081/api/users/profile \
  -H "Authorization: Bearer <access_token>"
```

## Troubleshooting

### Частые проблемы:

1. **401 Unauthorized**: Проверьте правильность токена и его срок действия
2. **Invalid token type**: Убедитесь, что используете access токен для защищенных ресурсов
3. **Token expired**: Используйте refresh токен для получения нового access токена

### Логирование:

Включите debug логирование в Go для диагностики проблем с токенами. 