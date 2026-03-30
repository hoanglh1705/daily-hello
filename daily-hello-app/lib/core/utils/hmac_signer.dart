import 'dart:convert';
import 'dart:math';

import 'package:crypto/crypto.dart';

class HmacSigner {
  static const String _secretKey = 'dh-hmac-secret-change-me';

  static Map<String, String> sign(String body) {
    final timestamp =
        (DateTime.now().millisecondsSinceEpoch ~/ 1000).toString();
    final nonce = _generateNonce();
    final message = '$timestamp.$nonce.$body';
    final hmacSha256 = Hmac(sha256, utf8.encode(_secretKey));
    final digest = hmacSha256.convert(utf8.encode(message));

    return {
      'X-Timestamp': timestamp,
      'X-Nonce': nonce,
      'X-Signature': digest.toString(),
    };
  }

  static String _generateNonce() {
    final random = Random.secure();
    final bytes = List<int>.generate(16, (_) => random.nextInt(256));
    return bytes.map((b) => b.toRadixString(16).padLeft(2, '0')).join();
  }
}
