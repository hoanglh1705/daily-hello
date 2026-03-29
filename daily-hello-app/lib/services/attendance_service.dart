import 'dart:io';

import 'package:device_info_plus/device_info_plus.dart';
import 'package:dio/dio.dart';
import '../core/network/api_response.dart';
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
    final deviceId = await _getDeviceId();
    final res = await dio.post('/v1/attendance/check-in', data: {
      'lat': lat,
      'lng': lng,
      'wifi_ssid': wifiSsid,
      'wifi_bssid': wifiBssid,
      'device_id': deviceId,
    });
    return Attendance.fromJson(
      unwrapApiData(res.data) as Map<String, dynamic>,
    );
  }

  Future<Attendance> checkOut({
    required double lat,
    required double lng,
    String? wifiSsid,
    String? wifiBssid,
  }) async {
    final deviceId = await _getDeviceId();
    final res = await dio.post('/v1/attendance/check-out', data: {
      'lat': lat,
      'lng': lng,
      'wifi_ssid': wifiSsid,
      'wifi_bssid': wifiBssid,
      'device_id': deviceId,
    });
    return Attendance.fromJson(
      unwrapApiData(res.data) as Map<String, dynamic>,
    );
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

  Future<List<Attendance>> getHistory({
    String? from,
    String? to,
    int page = 1,
    int limit = 20,
  }) async {
    final res = await dio.get('/v1/attendance', queryParameters: {
      'from': from,
      'to': to,
      'page': page,
      'limit': limit,
    });
    final list = unwrapApiData(res.data) as List<dynamic>;
    return list.map((e) => Attendance.fromJson(e as Map<String, dynamic>)).toList();
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
