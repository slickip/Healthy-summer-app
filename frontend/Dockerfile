# Используем официальный образ Flutter
FROM dart:stable AS dart-base

# Flutter SDK добавим отдельно
ENV FLUTTER_HOME=/flutter
ENV PATH="$PATH:/flutter/bin"

# Скачиваем Flutter
RUN git clone https://github.com/flutter/flutter.git -b stable /flutter

# Включаем web
RUN flutter config --enable-web

# Кешируем зависимости
WORKDIR /app/frontend
COPY /pubspec.* ./
RUN flutter pub get

# Копируем остальной проект
COPY . ./



# Собираем веб
RUN flutter build web --release

# Второй этап - минимальный nginx для отдачи
FROM nginx:alpine
COPY --from=dart-base /app/frontend/build/web /usr/share/nginx/html

# Порт для локального запуска (по умолчанию 80)
EXPOSE 8085

# Запускаем nginx
CMD ["nginx", "-g", "daemon off;"]
