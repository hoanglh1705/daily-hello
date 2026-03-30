import 'dart:io';

import 'package:device_info_plus/device_info_plus.dart';
import 'package:dio/dio.dart';
import '../core/network/api_response.dart';
import '../models/device.dart';

class DeviceService {
  final Dio dio;
  final DeviceInfoPlugin _deviceInfo = DeviceInfoPlugin();

  DeviceService(this.dio);

  Future<String?> getDeviceId() async {
    try {
      if (Platform.isAndroid) {
        final info = await _deviceInfo.androidInfo;
        return info.id;
      }
      if (Platform.isIOS) {
        final info = await _deviceInfo.iosInfo;
        return info.identifierForVendor;
      }
    } catch (_) {}
    return null;
  }

  Future<({String platform, String model, String deviceName})> getDeviceInfo() async {
    try {
      if (Platform.isAndroid) {
        final info = await _deviceInfo.androidInfo;
        return (
          platform: 'android',
          model: info.model,
          deviceName: '${info.brand} ${info.model}',
        );
      }
      if (Platform.isIOS) {
        final info = await _deviceInfo.iosInfo;
        return (
          platform: 'ios',
          model: info.utsname.machine,
          deviceName: info.name,
        );
      }
    } catch (_) {}
    return (platform: 'unknown', model: 'unknown', deviceName: 'unknown');
  }

  /// Check device status. Returns null if device not found (404).
  Future<Device?> checkStatus(String deviceId) async {
    try {
      final res = await dio.get(
        '/v1/devices/status',
        queryParameters: {'device_id': deviceId},
      );
      return Device.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
    } on DioException catch (e) {
      if (e.response?.statusCode == 404) return null;
      rethrow;
    }
  }

  /// Register device. Returns the created/existing device.
  Future<Device> register({
    required String deviceId,
    required String deviceName,
    required String platform,
    required String model,
  }) async {
    final res = await dio.post('/v1/devices/register', data: {
      'device_id': deviceId,
      'device_name': deviceName,
      'platform': platform,
      'model': model,
    });
    return Device.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
  }
}
