import 'package:dio/dio.dart';
import '../core/network/api_response.dart';
import '../models/auth_tokens.dart';
import '../models/user.dart';

class AuthService {
  final Dio dio;

  AuthService(this.dio);

  Future<AuthTokens> login(String email, String password) async {
    final res = await dio.post('/v1/auth/login', data: {
      'email': email,
      'password': password,
    });
    return AuthTokens.fromJson(
      unwrapApiData(res.data) as Map<String, dynamic>,
    );
  }

  Future<AuthTokens> refreshToken(String refreshToken) async {
    final res = await dio.post('/v1/auth/refresh-token', data: {
      'refresh_token': refreshToken,
    });
    return AuthTokens.fromJson(
      unwrapApiData(res.data) as Map<String, dynamic>,
    );
  }

  Future<void> logout(String refreshToken) async {
    await dio.post('/v1/auth/logout', data: {
      'refresh_token': refreshToken,
    });
  }

  Future<User> getProfile() async {
    final res = await dio.get('/v1/users/me');
    return User.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
  }
}
