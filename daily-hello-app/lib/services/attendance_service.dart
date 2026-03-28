import 'package:dio/dio.dart';
import '../core/network/api_response.dart';
import '../models/attendance.dart';

class AttendanceService {
  final Dio dio;

  AttendanceService(this.dio);

  Future<Attendance> checkIn({
    required double lat,
    required double lng,
    String? wifiSsid,
    String? wifiBssid,
  }) async {
    final res = await dio.post('/attendance/check-in', data: {
      'lat': lat,
      'lng': lng,
      'wifi_ssid': wifiSsid,
      'wifi_bssid': wifiBssid,
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
    final res = await dio.post('/attendance/check-out', data: {
      'lat': lat,
      'lng': lng,
      'wifi_ssid': wifiSsid,
      'wifi_bssid': wifiBssid,
    });
    return Attendance.fromJson(
      unwrapApiData(res.data) as Map<String, dynamic>,
    );
  }

  Future<List<Attendance>> getHistory({
    String? from,
    String? to,
    int page = 1,
    int limit = 20,
  }) async {
    final res = await dio.get('/attendance', queryParameters: {
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
      final res = await dio.get('/attendance/today');
      return Attendance.fromJson(
        unwrapApiData(res.data) as Map<String, dynamic>,
      );
    } catch (_) {
      return null;
    }
  }
}
