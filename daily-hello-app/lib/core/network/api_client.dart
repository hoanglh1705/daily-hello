import 'dart:async';
import 'package:dio/dio.dart';
import '../constants/app_constants.dart';
import 'api_response.dart';
import '../storage/secure_storage.dart';

class ApiClient {
  late final Dio dio;

  final void Function()? onUnauthorized;

  bool _isRefreshing = false;
  final List<Completer<bool>> _pendingQueue = [];

  ApiClient(SecureStorage storage, {this.onUnauthorized}) {
    dio = Dio(BaseOptions(
      baseUrl: AppConstants.baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {'Content-Type': 'application/json'},
    ));

    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        if (_isPublicEndpoint(options.path)) {
          return handler.next(options);
        }

        final isExpired = await storage.isAccessTokenExpired();
        if (isExpired) {
          final refreshed = await _handleRefresh(storage);
          if (!refreshed) {
            onUnauthorized?.call();
            return handler.reject(
              DioException(
                requestOptions: options,
                type: DioExceptionType.cancel,
                message: 'Session expired',
              ),
            );
          }
        }

        final token = await storage.getAccessToken();
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        return handler.next(options);
      },

      onError: (DioException error, handler) async {
        if (error.response?.statusCode == 401 &&
            !_isPublicEndpoint(error.requestOptions.path)) {
          final refreshed = await _handleRefresh(storage);
          if (refreshed) {
            final token = await storage.getAccessToken();
            final opts = error.requestOptions;
            opts.headers['Authorization'] = 'Bearer $token';
            try {
              final response = await dio.fetch(opts);
              return handler.resolve(response);
            } catch (_) {
              return handler.next(error);
            }
          } else {
            onUnauthorized?.call();
          }
        }
        return handler.next(error);
      },
    ));
  }

  bool _isPublicEndpoint(String path) {
    return path.contains('/auth/login') ||
        path.contains('/auth/refresh-token');
  }

  Future<bool> _handleRefresh(SecureStorage storage) async {
    if (_isRefreshing) {
      final completer = Completer<bool>();
      _pendingQueue.add(completer);
      return completer.future;
    }

    _isRefreshing = true;
    try {
      final refreshToken = await storage.getRefreshToken();
      if (refreshToken == null) {
        _resolvePending(false);
        return false;
      }

      final res = await dio.post(
        '/v1/auth/refresh-token',
        data: {'refresh_token': refreshToken},
      );

      if (res.statusCode == 200) {
        final data = unwrapApiData(res.data) as Map<String, dynamic>;
        await storage.saveAuthTokens(
          accessToken: data['access_token'] as String,
          refreshToken: data['refresh_token'] as String,
          tokenType: data['token_type'] as String? ?? 'Bearer',
          expiresIn: data['expires_in'] as int,
        );
        _resolvePending(true);
        return true;
      }

      _resolvePending(false);
      return false;
    } catch (_) {
      _resolvePending(false);
      return false;
    } finally {
      _isRefreshing = false;
    }
  }

  void _resolvePending(bool success) {
    for (final completer in _pendingQueue) {
      completer.complete(success);
    }
    _pendingQueue.clear();
  }
}
