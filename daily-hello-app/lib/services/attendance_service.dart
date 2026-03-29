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
    final payload = await _buildPayload(lat: lat, lng: lng, wifiSsid: wifiSsid, wifiBssid: wifiBssid);
    final res = await dio.post('/v1/attendance/check-in', data: payload);
    return Attendance.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
  }

  Future<Attendance> checkOut({
    required double lat,
    required double lng,
    String? wifiSsid,
    String? wifiBssid,
  }) async {
    final payload = await _buildPayload(lat: lat, lng: lng, wifiSsid: wifiSsid, wifiBssid: wifiBssid);
    final res = await dio.post('/v1/attendance/check-out', data: payload);
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
    final res = await dio.get('/v1/attendance/history', queryParameters: {
      'from': from,
      'to': to,
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
