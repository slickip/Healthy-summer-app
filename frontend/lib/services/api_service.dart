import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiService {
  static const String _baseUrl = 'http://localhost:8081/api'; // или твой IP
  final _storage = const FlutterSecureStorage();

  // Ключи для хранения токенов
  static const String _accessTokenKey = 'access_token';
  static const String _refreshTokenKey = 'refresh_token';

  Future<Map<String, dynamic>?> login(String email, String password) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/users/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'email': email, 'password': password}),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final accessToken = data['access_token'];
      final refreshToken = data['refresh_token'];

      // Сохраняем оба токена
      await _storage.write(key: _accessTokenKey, value: accessToken);
      await _storage.write(key: _refreshTokenKey, value: refreshToken);

      return data;
    } else {
      return null;
    }
  }

  Future<Map<String, dynamic>?> register(
    String email,
    String password,
    String displayName,
  ) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/users/register'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'email': email,
        'password': password,
        'display_name': displayName,
      }),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final accessToken = data['access_token'];
      final refreshToken = data['refresh_token'];

      // Сохраняем оба токена
      await _storage.write(key: _accessTokenKey, value: accessToken);
      await _storage.write(key: _refreshTokenKey, value: refreshToken);

      return data;
    } else {
      return null;
    }
  }

  Future<String?> _getAccessToken() async {
    return await _storage.read(key: _accessTokenKey);
  }

  Future<String?> _getRefreshToken() async {
    return await _storage.read(key: _refreshTokenKey);
  }

  Future<bool> _refreshAccessToken() async {
    final refreshToken = await _getRefreshToken();
    if (refreshToken == null) return false;

    final response = await http.post(
      Uri.parse('$_baseUrl/users/refresh'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'refresh_token': refreshToken}),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final newAccessToken = data['access_token'];
      await _storage.write(key: _accessTokenKey, value: newAccessToken);
      return true;
    } else {
      // Если refresh токен тоже истек, очищаем все токены
      await logout();
      return false;
    }
  }

  Future<Map<String, dynamic>?> getProfile() async {
    final token = await _getAccessToken();
    if (token == null) return null;

    final response = await http.get(
      Uri.parse('$_baseUrl/users/profile'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
    );

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    } else if (response.statusCode == 401) {
      // Токен истек, пробуем обновить
      if (await _refreshAccessToken()) {
        // Повторяем запрос с новым токеном
        final newToken = await _getAccessToken();
        final retryResponse = await http.get(
          Uri.parse('$_baseUrl/users/profile'),
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer $newToken',
          },
        );

        if (retryResponse.statusCode == 200) {
          return jsonDecode(retryResponse.body);
        }
      }
    }

    return null;
  }

  Future<bool> createActivity({
    required int activityTypeId,
    required int duration,
    required String intensity,
  }) async {
    final token = await _getAccessToken();
    if (token == null) return false;

    final response = await http.post(
      Uri.parse('http://localhost:8082/api/activities'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode({
        'activity_type_id': activityTypeId,
        'duration': duration,
        'intensity': intensity,
      }),
    );

    if (response.statusCode == 201) {
      return true;
    } else if (response.statusCode == 401) {
      if (await _refreshAccessToken()) {
        final newToken = await _getAccessToken();
        final retryResponse = await http.post(
          Uri.parse('http://localhost:8082/api/activities'),
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer $newToken',
          },
          body: jsonEncode({
            'activity_type_id': activityTypeId,
            'duration': duration,
            'intensity': intensity,
          }),
        );
        return retryResponse.statusCode == 201;
      }
    }
    return false;
  }

  Future<List<dynamic>?> getActivities() async {
    final token = await _getAccessToken();
    if (token == null) return null;

    final response = await http.get(
      Uri.parse('http://localhost:8082/api/activities'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
    );

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    } else if (response.statusCode == 401) {
      // Токен истек, пробуем обновить
      if (await _refreshAccessToken()) {
        final newToken = await _getAccessToken();
        final retryResponse = await http.get(
          Uri.parse('http://localhost:8082/api/activities'),
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer $newToken',
          },
        );

        if (retryResponse.statusCode == 200) {
          return jsonDecode(retryResponse.body);
        }
      }
    }

    return null;
  }

  Future<void> logout() async {
    await _storage.delete(key: _accessTokenKey);
    await _storage.delete(key: _refreshTokenKey);
  }

  Future<bool> isLoggedIn() async {
    final accessToken = await _getAccessToken();
    return accessToken != null;
  }
}
