import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiService {
  static const String _baseUrl = 'http://localhost:8081/api';
  static const String _nutritionUrl =
      'http://localhost:8083/api'; // или твой IP
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

  Future<bool> createMeal({
    required String mealTime,
    required String description,
    required int calories,
  }) async {
    final token = await _getAccessToken();
    if (token == null) return false;

    final response = await http.post(
      Uri.parse('$_nutritionUrl/meals'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode({
        'meal_time': mealTime,
        'description': description,
        'calories': calories,
      }),
    );

    return response.statusCode == 201;
  }

  Future<bool> updateMeal({
    required int id,
    String? mealTime,
    String? description,
    int? calories,
  }) async {
    final token = await _getAccessToken();
    if (token == null) return false;

    final body = <String, dynamic>{};
    if (mealTime != null) body['meal_time'] = mealTime;
    if (description != null) body['description'] = description;
    if (calories != null) body['calories'] = calories;

    final response = await http.put(
      Uri.parse('$_nutritionUrl/meals?id=$id'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode(body),
    );

    return response.statusCode == 200;
  }

  Future<bool> deleteMeal(int id) async {
    final token = await _getAccessToken();
    if (token == null) return false;

    final response = await http.delete(
      Uri.parse('$_nutritionUrl/meals?id=$id'),
      headers: {'Authorization': 'Bearer $token'},
    );

    return response.statusCode == 204;
  }

  Future<List<dynamic>?> getMeals() async {
    final token = await _getAccessToken();
    if (token == null) return null;

    final response = await http.get(
      Uri.parse('$_nutritionUrl/meals'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
    );

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    }
    return null;
  }

  Future<bool> createWaterLog({
    required int volumeML,
    required String loggedAt,
  }) async {
    final token = await _getAccessToken();
    if (token == null) return false;

    final response = await http.post(
      Uri.parse('$_nutritionUrl/water'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode({'volume_ml': volumeML, 'logged_at': loggedAt}),
    );

    return response.statusCode == 201;
  }

  Future<List<dynamic>?> getWaterLogs() async {
    final token = await _getAccessToken();
    if (token == null) return null;

    final response = await http.get(
      Uri.parse('$_nutritionUrl/water'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
    );

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    }
    return null;
  }

  Future<bool> deleteWaterLog(int id) async {
    final token = await _getAccessToken();
    if (token == null) return false;

    final response = await http.delete(
      Uri.parse('$_nutritionUrl/water?id=$id'),
      headers: {'Authorization': 'Bearer $token'},
    );

    return response.statusCode == 204;
  }

  Future<bool> createFood({
    required String name,
    required double caloriesPer100g,
    required double proteins,
    required double fats,
    required double carbs,
  }) async {
    final token = await _getAccessToken();
    if (token == null) return false;

    final body = {
      'name': name,
      'callories_per_100g': caloriesPer100g,
      'proteins': proteins,
      'fats': fats,
      'carbs': carbs,
    };

    print('API createFood - sending body: $body');

    final response = await http.post(
      Uri.parse('$_nutritionUrl/foods'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode(body),
    );

    print('API createFood - response status: ${response.statusCode}');
    print('API createFood - response body: ${response.body}');

    return response.statusCode == 201;
  }

  Future<List<dynamic>?> getFoods() async {
    final token = await _getAccessToken();
    if (token == null) return null;

    final response = await http.get(
      Uri.parse('$_nutritionUrl/foods'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
    );

    print('API getFoods - response status: ${response.statusCode}');
    print('API getFoods - response body: ${response.body}');

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    }
    return null;
  }

  Future<bool> updateFood({
    required int id,
    String? name,
    double? caloriesPer100g,
    double? proteins,
    double? fats,
    double? carbs,
  }) async {
    final token = await _getAccessToken();
    if (token == null) return false;

    final body = <String, dynamic>{};
    if (name != null) body['name'] = name;
    if (caloriesPer100g != null) body['callories_per_100g'] = caloriesPer100g;
    if (proteins != null) body['proteins'] = proteins;
    if (fats != null) body['fats'] = fats;
    if (carbs != null) body['carbs'] = carbs;

    final response = await http.put(
      Uri.parse('$_nutritionUrl/foods?id=$id'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode(body),
    );

    return response.statusCode == 200;
  }

  Future<bool> deleteFood(int id) async {
    final token = await _getAccessToken();
    if (token == null) return false;

    final response = await http.delete(
      Uri.parse('$_nutritionUrl/foods?id=$id'),
      headers: {'Authorization': 'Bearer $token'},
    );

    return response.statusCode == 204;
  }
}
