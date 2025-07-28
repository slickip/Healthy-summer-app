# Healthy Summer - Health & Fitness Application

A full-featured mobile application for tracking physical activity, nutrition, and social interaction, built on a microservices architecture.

## Technology Stack

### Frontend
- **Flutter** - Cross-platform mobile app development
- **Dart** - Programming language
- **HTTP** - For API requests
- **Flutter Secure Storage** - Secure token storage
- **Intl** - Internationalization

### Backend (Go Microservices)
- **Go** - Backend programming language
- **PostgreSQL** - Relational database
- **JWT** - Authentication and authorization
- **Docker** - Containerization
- **Docker Compose** - Service orchestration

### Architecture
- **Microservices Architecture** - Separation into independent services
- **REST API** - Service communication
- **CORS** - Cross-domain request configuration

## Project Architecture

```
healthy_summer/
├── frontend/                 # Flutter application
│   ├── lib/
│   │   ├── screens/         # Application screens
│   │   ├── services/        # API services
│   │   ├── models/          # Data models
│   │   └── widgets/         # Reusable widgets
│   └── pubspec.yaml         # Flutter dependencies
├── backend/                  # Microservices
│   ├── user-service/        # User management service
│   ├── activity-service/    # Activity tracking service
│   ├── nutrition-service/   # Nutrition management service
│   └── social-service/      # Social features service
├── proto/                   # Protocol Buffers schemas
└── docker-compose.yml       # Docker configuration
```

## Microservices

### 1. User Service (Port: 8081)
- User registration and authentication
- User profile management
- JWT tokens (access/refresh)
- Authentication middleware

### 2. Activity Service (Port: 8082)
- Physical activity tracking
- Exercise type management
- Activity statistics
- Step counter integration

### 3. Nutrition Service (Port: 8083)
- Meal management
- Water consumption tracking
- Food database
- Calorie and macronutrient tracking

### 4. Social Service (Port: 8084)
- Social features and challenges
- Activity feed
- Messaging system
- Friend management

## Application Features

### Main Screens:
- **Welcome Screen** - Welcome page
- **Login/Register** - Authentication and registration
- **Home Screen** - Main dashboard
- **Activities** - Physical activity management
- **Nutrition** - Nutrition and water management
- **Social** - Social features

### Key Capabilities:
- ✅ User registration and authentication
- ✅ Physical activity tracking
- ✅ Nutrition and water consumption management
- ✅ Social features and challenges
- ✅ Secure data storage
- ✅ Cross-platform compatibility (iOS/Android)

## Installation and Setup

### Prerequisites
- Docker and Docker Compose
- Flutter SDK (for frontend development)
- Go (for backend development)

### Project Setup

1. **Clone Repository**
```bash
git clone <repository-url>
cd healthy_summer
```

2. **Start All Services**
```bash
docker-compose up -d
```

3. **Verify Service Status**
```bash
# User Service
curl http://localhost:8081/api/ping

# Activity Service
curl http://localhost:8082/api/ping

# Nutrition Service
curl http://localhost:8083/api/ping

# Social Service
curl http://localhost:8084/api/ping
```

4. **Frontend Development**
```bash
cd frontend
flutter pub get
flutter run
```

## Database

- **PostgreSQL 15** - Main database
- **Port**: 5433 (external), 5432 (internal)
- **User**: healthyuser
- **Database**: healthydb

## Security

- JWT tokens for authentication
- Middleware for endpoint protection
- CORS settings for secure requests
- Secure token storage in Flutter

## API Endpoints

### User Service
- `POST /api/users/register` - User registration
- `POST /api/users/login` - User authentication
- `POST /api/users/refresh` - Token refresh
- `GET /api/users/profile` - User profile

### Activity Service
- `GET /api/activities` - List activities
- `POST /api/activities` - Create activity
- `PUT /api/activities/:id` - Update activity
- `DELETE /api/activities/:id` - Delete activity

### Nutrition Service
- `GET /api/meals` - Meal management
- `POST /api/meals` - Create meal
- `PUT /api/meals/:id` - Update meal
- `DELETE /api/meals/:id` - Delete meal
- `GET /api/foods` - Food database
- `POST /api/foods` - Add food
- `GET /api/water` - Water consumption
- `POST /api/water` - Add water intake

### Social Service
- `GET /api/challenges` - List challenges
- `POST /api/challenges` - Create challenge
- `PUT /api/challenges/:id` - Update challenge
- `DELETE /api/challenges/:id` - Delete challenge
- `POST /api/challenges/:id/join` - Join challenge
- `GET /api/challenges/my` - User's challenges
- `GET /api/challenges/:id/leaderboard` - Challenge leaderboard
- `GET /api/social/feed/friends` - Friends activity feed
- `POST /api/messages` - Send message
- `GET /api/messages` - Get conversation history
- `POST /api/friend-requests` - Send friend request
- `PUT /api/friend-requests/:id` - Accept/reject friend request
- `GET /api/friend-requests` - List friend requests

## Docker

All services are containerized using Docker:
- Separate containers for each microservice
- Shared PostgreSQL database
- Frontend accessible on port 8085
- Automatic service restart

## Deployment

Project is ready for production deployment:
- Docker Compose for local deployment
- Kubernetes deployment capability
- Environment variable configuration
- Logging and monitoring

## Future Enhancements

- Fitness tracker integration
- Machine learning for personalization
- Push notifications
- Gamification and achievements
- Social media integration
- Analytics and reporting

## Development

Project demonstrates:
- Modern development practices
- Microservices architecture
- Security and authentication
- Cross-platform development
- Containerization and DevOps

---
