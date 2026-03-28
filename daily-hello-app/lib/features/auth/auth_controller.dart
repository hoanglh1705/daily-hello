import 'package:flutter/material.dart';
import '../../core/network/api_response.dart';
import '../../models/user.dart';
import '../../services/auth_service.dart';
import '../../core/storage/secure_storage.dart';

class AuthController extends ChangeNotifier {
  final AuthService _authService;
  final SecureStorage _storage;

  User? currentUser;
  bool isLoading = false;
  String? errorMessage;

  AuthController(this._authService, this._storage);

  Future<bool> login(String email, String password) async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      final tokens = await _authService.login(email, password);
      await _storage.saveAuthTokens(
        accessToken: tokens.accessToken,
        refreshToken: tokens.refreshToken,
        tokenType: tokens.tokenType,
        expiresIn: tokens.expiresIn,
      );
      return true;
    } catch (error, stackTrace) {
      debugPrint('Login failed: $error\n$stackTrace');
      errorMessage = getApiErrorMessage(error) ??
          'Đăng nhập thất bại. Vui lòng kiểm tra lại.';
      return false;
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<void> loadProfile() async {
    try {
      currentUser = await _authService.getProfile();
      notifyListeners();
    } catch (_) {}
  }

  Future<void> logout() async {
    try {
      final refreshToken = await _storage.getRefreshToken();
      if (refreshToken != null) {
        await _authService.logout(refreshToken);
      }
    } catch (_) {}
    await _storage.clearAuth();
    currentUser = null;
    notifyListeners();
  }

  Future<bool> isLoggedIn() => _storage.hasToken();
}
