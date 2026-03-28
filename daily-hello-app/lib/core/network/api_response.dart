import 'package:dio/dio.dart';

class ApiException implements Exception {
  final String? errorCode;
  final String? errorMessage;

  const ApiException({this.errorCode, this.errorMessage});

  factory ApiException.fromResponse(Map<String, dynamic> json) {
    return ApiException(
      errorCode: json['error_code'] as String?,
      errorMessage: json['error_message'] as String?,
    );
  }

  @override
  String toString() => errorMessage ?? 'API request failed';
}

dynamic unwrapApiData(dynamic responseData) {
  if (responseData is Map<String, dynamic>) {
    final success = responseData['success'];
    if (success == false) {
      throw ApiException.fromResponse(responseData);
    }

    if (responseData.containsKey('data')) {
      return responseData['data'];
    }
  }

  return responseData;
}

String? getApiErrorMessage(Object error) {
  if (error is ApiException) {
    return error.errorMessage;
  }

  if (error is DioException) {
    final responseData = error.response?.data;
    if (responseData is Map<String, dynamic>) {
      return ApiException.fromResponse(responseData).errorMessage;
    }
  }

  return null;
}
