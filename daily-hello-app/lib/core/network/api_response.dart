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
      final apiErrorMessage = ApiException.fromResponse(responseData).errorMessage;
      if (apiErrorMessage != null && apiErrorMessage.isNotEmpty) {
        return apiErrorMessage;
      }
    }

    final statusCode = error.response?.statusCode;
    if (statusCode != null) {
      switch (statusCode) {
        case 400:
          return 'Yeu cau khong hop le.';
        case 401:
          return 'Phien dang nhap da het han. Vui long dang nhap lai.';
        case 403:
          return 'Ban khong co quyen thuc hien thao tac nay.';
        case 404:
          return 'Khong tim thay du lieu hoac API yeu cau.';
        case 408:
          return 'Ket noi bi timeout. Vui long thu lai.';
      }

      if (statusCode >= 500) {
        return 'He thong dang ban. Vui long thu lai sau.';
      }
    }

    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return 'Ket noi bi timeout. Vui long thu lai.';
      case DioExceptionType.connectionError:
        return 'Khong the ket noi may chu. Vui long kiem tra mang.';
      case DioExceptionType.cancel:
        return 'Yeu cau da bi huy.';
      case DioExceptionType.badCertificate:
        return 'Chung chi bao mat khong hop le.';
      case DioExceptionType.badResponse:
        return 'Khong the xu ly yeu cau luc nay.';
      case DioExceptionType.unknown:
        return 'Da xay ra loi khong xac dinh. Vui long thu lai.';
    }
  }

  return null;
}
