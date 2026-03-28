import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../constants/app_constants.dart';

class SecureStorage {
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  // ── Save ──────────────────────────────────────────────────────────────────

  Future<void> saveAuthTokens({
    required String accessToken,
    required String refreshToken,
    required String tokenType,
    required int expiresIn,
  }) async {
    final expiresAt = DateTime.now()
        .add(Duration(seconds: expiresIn))
        .toIso8601String();

    await Future.wait([
      _storage.write(key: AppConstants.tokenKey, value: accessToken),
      _storage.write(key: AppConstants.refreshTokenKey, value: refreshToken),
      _storage.write(key: AppConstants.tokenTypeKey, value: tokenType),
      _storage.write(key: AppConstants.expiresAtKey, value: expiresAt),
    ]);
  }

  // ── Read ──────────────────────────────────────────────────────────────────

  Future<String?> getAccessToken() =>
      _storage.read(key: AppConstants.tokenKey);

  Future<String?> getRefreshToken() =>
      _storage.read(key: AppConstants.refreshTokenKey);

  Future<String?> getTokenType() =>
      _storage.read(key: AppConstants.tokenTypeKey);

  Future<bool> isAccessTokenExpired() async {
    final raw = await _storage.read(key: AppConstants.expiresAtKey);
    if (raw == null) return true;
    final expiresAt = DateTime.tryParse(raw);
    if (expiresAt == null) return true;
    // Trừ 10 giây buffer để tránh race condition
    return DateTime.now().isAfter(expiresAt.subtract(const Duration(seconds: 10)));
  }

  Future<bool> hasToken() async {
    final token = await getAccessToken();
    return token != null && token.isNotEmpty;
  }

  // ── Delete ────────────────────────────────────────────────────────────────

  Future<void> clearAuth() async {
    await Future.wait([
      _storage.delete(key: AppConstants.tokenKey),
      _storage.delete(key: AppConstants.refreshTokenKey),
      _storage.delete(key: AppConstants.tokenTypeKey),
      _storage.delete(key: AppConstants.expiresAtKey),
    ]);
  }
}
