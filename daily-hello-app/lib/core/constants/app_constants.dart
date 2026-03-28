import 'package:flutter/foundation.dart';

class AppConstants {
  static String get baseUrl {
    if (kIsWeb) {
      return 'http://localhost:8282/api';
    }

    return switch (defaultTargetPlatform) {
      TargetPlatform.android => 'http://10.0.2.2:8282/api',
      _ => 'http://localhost:8282/api',
    };
  }

  // Storage keys
  static const String tokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String tokenTypeKey = 'token_type';
  static const String expiresAtKey = 'expires_at'; // ISO8601 string
  static const String userKey = 'current_user';
}
