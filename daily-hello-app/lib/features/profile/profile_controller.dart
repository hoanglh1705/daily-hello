import 'package:flutter/material.dart';
import '../../models/user.dart';
import '../../services/auth_service.dart';
import '../../core/storage/secure_storage.dart';

class ProfileController extends ChangeNotifier {
  final AuthService _authService;
  final SecureStorage _storage;

  User? user;
  bool isLoading = false;

  ProfileController(this._authService, this._storage);

  Future<void> loadProfile() async {
    isLoading = true;
    notifyListeners();
    try {
      user = await _authService.getProfile();
    } catch (_) {}
    isLoading = false;
    notifyListeners();
  }

  Future<void> logout(BuildContext context) async {
    try {
      final refreshToken = await _storage.getRefreshToken();
      if (refreshToken != null) {
        await _authService.logout(refreshToken);
      }
    } catch (_) {}
    await _storage.clearAuth();
    if (context.mounted) {
      Navigator.pushNamedAndRemoveUntil(context, '/login', (_) => false);
    }
  }
}
