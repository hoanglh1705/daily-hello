import 'package:flutter/material.dart';
import '../../models/device.dart';
import '../../services/device_service.dart';

class DeviceController extends ChangeNotifier {
  final DeviceService _service;

  Device? device;
  bool isChecking = false;
  bool isRegistering = false;
  String? error;

  DeviceController(this._service);

  /// Check if current device is registered for this user.
  /// Returns the device if found, null if not registered.
  Future<Device?> checkDevice() async {
    isChecking = true;
    error = null;
    notifyListeners();

    try {
      final deviceId = await _service.getDeviceId();
      if (deviceId == null) {
        error = 'Không thể lấy thông tin thiết bị.';
        isChecking = false;
        notifyListeners();
        return null;
      }

      device = await _service.checkStatus(deviceId);
      isChecking = false;
      notifyListeners();
      return device;
    } catch (e) {
      error = 'Lỗi kiểm tra thiết bị.';
      isChecking = false;
      notifyListeners();
      return null;
    }
  }

  /// Register the current device.
  Future<Device?> registerDevice() async {
    isRegistering = true;
    error = null;
    notifyListeners();

    try {
      final deviceId = await _service.getDeviceId();
      if (deviceId == null) {
        error = 'Không thể lấy thông tin thiết bị.';
        isRegistering = false;
        notifyListeners();
        return null;
      }

      final info = await _service.getDeviceInfo();
      device = await _service.register(
        deviceId: deviceId,
        deviceName: info.deviceName,
        platform: info.platform,
        model: info.model,
      );
      isRegistering = false;
      notifyListeners();
      return device;
    } catch (e) {
      error = 'Đăng ký thiết bị thất bại.';
      isRegistering = false;
      notifyListeners();
      return null;
    }
  }
}
