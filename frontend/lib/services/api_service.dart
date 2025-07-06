import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiService {
  static const String _baseUrl = 'http://localhost:8081/api'; // или твой IP
  final _storage = const FlutterSecureStorage();

  Future<String?> login(String email, String password) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/users/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'email': email, 'password': password}),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final token = data['token'];
      await _storage.write(key: 'jwt_token', value: token);
      return token;
    } else {
      return null;
    }
  }

  Future<bool> register(
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
      final token = data['token'];
      await _storage.write(key: 'jwt_token', value: token);
      return true;
    } else {
      return false;
    }
  }

  Future<Map<String, dynamic>?> getProfile() async {
    final token = await _storage.read(key: 'jwt_token');
    print("Token being sent: $token");
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
    } else {
      return null;
    }
  }

  Future<bool> createActivity({
    required String type,
    required int duration,
    required String intensity,
    required String location,
  }) async {
    final token = await _storage.read(key: 'jwt_token');
    print("Token being sent: $token");
    if (token == null) return false;

    final response = await http.post(
      Uri.parse('http://localhost:8082/api/activities'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode({
        'type': type,
        'duration': duration,
        'intensity': intensity,
        'location': location,
      }),
    );

    return response.statusCode == 201;
  }

  Future<List<dynamic>?> getActivities() async {
    final token = await _storage.read(key: 'jwt_token');
    print("Token being sent: $token");
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
    } else {
      return null;
    }
  }

  Future<void> logout() async {
    await _storage.delete(key: 'jwt_token');
  }
}
