import 'dart:convert';
import 'dart:io';

import 'package:device_info_plus/device_info_plus.dart';
import 'package:dio/dio.dart';
import '../core/network/api_response.dart';
import '../core/utils/hmac_signer.dart';
import '../models/attendance.dart';

class AttendanceService {
  final Dio dio;
  final DeviceInfoPlugin _deviceInfo = DeviceInfoPlugin();

  AttendanceService(this.dio);

  Future<Attendance> checkIn({
    required double lat,
    required double lng,
    String? wifiSsid,
    String? wifiBssid,
  }) async {
    final payload = await _buildPayload(lat: lat, lng: lng, wifiSsid: wifiSsid, wifiBssid: wifiBssid);
    final body = jsonEncode(payload);
    final hmacHeaders = HmacSigner.sign(body);
    final res = await dio.post(
      '/v1/attendance/check-in',
      data: body,
      options: Options(headers: hmacHeaders),
    );
    return Attendance.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
  }

  Future<Attendance> checkOut({
    required double lat,
    required double lng,
    String? wifiSsid,
    String? wifiBssid,
  }) async {
    final payload = await _buildPayload(lat: lat, lng: lng, wifiSsid: wifiSsid, wifiBssid: wifiBssid);
    final body = jsonEncode(payload);
    final hmacHeaders = HmacSigner.sign(body);
    final res = await dio.post(
      '/v1/attendance/check-out',
      data: body,
      options: Options(headers: hmacHeaders),
    );
    if (res.data == null || (res.data is String && (res.data as String).trim().isEmpty)) {
      final attendance = await getTodayAttendance();
      if (attendance != null) return attendance;
      throw DioException(
        requestOptions: res.requestOptions,
        response: res,
        error: 'Check-out thành công nhưng không lấy được dữ liệu chấm công.',
      );
    }
    return Attendance.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
  }

  Future<Attendance> checkInGps({
    required double lat,
    required double lng,
    required String imageBase64,
  }) async {
    final deviceId = await _getDeviceId();
    final body = jsonEncode({
      'lat': lat,
      'lng': lng,
      'device_id': deviceId,
      'image': imageBase64,
    });
    final hmacHeaders = HmacSigner.sign(body);
    final res = await dio.post(
      '/v1/attendance/check-in-gps',
      data: body,
      options: Options(headers: hmacHeaders),
    );
    if (res.data == null || (res.data is String && (res.data as String).trim().isEmpty)) {
      final attendance = await getTodayAttendance();
      if (attendance != null) return attendance;
      throw DioException(
        requestOptions: res.requestOptions,
        response: res,
        error: 'Check-out thành công nhưng không lấy được dữ liệu chấm công.',
      );
    }
    return Attendance.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
  }

  Future<Attendance> checkOutGps({
    required double lat,
    required double lng,
    required String imageBase64,
  }) async {
    final deviceId = await _getDeviceId();
    final body = jsonEncode({
      'lat': lat,
      'lng': lng,
      'device_id': deviceId,
      'image': imageBase64,
    });
    final hmacHeaders = HmacSigner.sign(body);
    final res = await dio.post(
      '/v1/attendance/check-out-gps',
      data: body,
      options: Options(headers: hmacHeaders),
    );
    return Attendance.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
  }

  Future<Map<String, dynamic>> _buildPayload({
    required double lat,
    required double lng,
    String? wifiSsid,
    String? wifiBssid,
  }) async {
    final deviceId = await _getDeviceId();
    return {
      'lat': lat,
      'lng': lng,
      'wifi_ssid': wifiSsid,
      'wifi_bssid': wifiBssid,
      'device_id': deviceId,
    };
  }

  Future<String?> _getDeviceId() async {
    try {
      if (Platform.isAndroid) {
        final info = await _deviceInfo.androidInfo;
        return info.id;
      }
      if (Platform.isIOS) {
        final info = await _deviceInfo.iosInfo;
        return info.identifierForVendor;
      }
    } catch (_) {
      return null;
    }
    return null;
  }

  Future<({List<Attendance> items, int total})> getHistory({
    String? from,
    String? to,
    int page = 1,
    int limit = 20,
  }) async {
    final res = await dio.get('/v1/attendance/my-history', queryParameters: {
      // 'from': from,
      // 'to': to,
      'page': page,
      'limit': limit,
    });
    final data = unwrapApiData(res.data) as Map<String, dynamic>;
    final list = data['items'] as List<dynamic>? ?? [];
    final meta = data['meta'] as Map<String, dynamic>? ?? {};
    final total = (meta['total'] as num?)?.toInt() ?? 0;
    return (
      items: list.map((e) => Attendance.fromJson(e as Map<String, dynamic>)).toList(),
      total: total,
    );
  }

  Future<Attendance?> getTodayAttendance() async {
    try {
      final res = await dio.get('/v1/attendance/today');
      return Attendance.fromJson(
        unwrapApiData(res.data) as Map<String, dynamic>,
      );
    } catch (_) {
      return null;
    }
  }
}
