import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';
import 'package:network_info_plus/network_info_plus.dart';
import 'package:wifi_signal_strength_indicator/wifi_signal_strength.dart';
import '../../core/network/api_response.dart';
import '../../core/utils/location_permission_utils.dart';
import '../../models/attendance.dart';
import '../../services/attendance_service.dart';

class AttendanceController extends ChangeNotifier {
  final AttendanceService _service;
  final NetworkInfo _networkInfo = NetworkInfo();

  Attendance? todayAttendance;
  List<Attendance> history = [];
  bool isLoading = false;
  bool isLoadingHistory = false;
  String? errorMessage;
  int currentPage = 1;
  bool hasMore = true;

  // WiFi info state
  String? wifiSsid;
  int? wifiSignalStrength; // dBm
  int? wifiSignalLevel; // 0-4

  AttendanceController(this._service);

  String _formatError(Object error) {
    return error.toString().replaceFirst('Exception: ', '');
  }

  Future<void> loadTodayAttendance() async {
    todayAttendance = await _service.getTodayAttendance();
    notifyListeners();
  }

  Future<void> loadWifiInfo() async {
    try {
      wifiSsid = await _networkInfo.getWifiName();
      // Remove surrounding quotes if present (Android quirk)
      if (wifiSsid != null && wifiSsid!.startsWith('"') && wifiSsid!.endsWith('"')) {
        wifiSsid = wifiSsid!.substring(1, wifiSsid!.length - 1);
      }

      wifiSignalStrength = await WifiSignalStrength.getSignalStrength();
      wifiSignalLevel = await WifiSignalStrength.getSignalLevel();
    } catch (_) {
      wifiSsid = null;
      wifiSignalStrength = null;
      wifiSignalLevel = null;
    }
    notifyListeners();
  }

  Future<bool> checkIn() async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      final position = await _getPosition();
      final wifiInfo = await _getWifiInfo();
      todayAttendance = await _service.checkIn(
        lat: position.latitude,
        lng: position.longitude,
        wifiSsid: wifiInfo['ssid'],
        wifiBssid: wifiInfo['bssid'],
      );
      return true;
    } catch (error) {
      errorMessage =
          getApiErrorMessage(error) ?? 'Check-in thất bại: ${_formatError(error)}';
      return false;
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> checkOut() async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      final position = await _getPosition();
      final wifiInfo = await _getWifiInfo();
      todayAttendance = await _service.checkOut(
        lat: position.latitude,
        lng: position.longitude,
        wifiSsid: wifiInfo['ssid'],
        wifiBssid: wifiInfo['bssid'],
      );
      return true;
    } catch (error) {
      errorMessage = getApiErrorMessage(error) ??
          'Check-out thất bại: ${_formatError(error)}';
      return false;
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<void> loadHistory({bool refresh = false}) async {
    if (refresh) {
      currentPage = 1;
      hasMore = true;
      history = [];
    }
    if (!hasMore || isLoadingHistory) return;

    isLoadingHistory = true;
    notifyListeners();

    try {
      final result = await _service.getHistory(page: currentPage);
      history.addAll(result.items);
      currentPage++;
      hasMore = history.length < result.total;
    } catch (error) {
      errorMessage = getApiErrorMessage(error) ??
          'Lỗi tải lịch sử: ${_formatError(error)}';
    } finally {
      isLoadingHistory = false;
      notifyListeners();
    }
  }

  Future<Position> _getPosition() async {
    await LocationPermissionUtils.ensureLocationPermission();

    bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
    if (!serviceEnabled) throw Exception('GPS chưa được bật');

    return await Geolocator.getCurrentPosition(
      locationSettings: const LocationSettings(
        accuracy: LocationAccuracy.high,
      ),
    );
  }

  Future<Map<String, String?>> _getWifiInfo() async {
    try {
      final ssid = await _networkInfo.getWifiName();
      final bssid = await _networkInfo.getWifiBSSID();
      return {'ssid': ssid, 'bssid': bssid};
    } catch (_) {
      return {'ssid': null, 'bssid': null};
    }
  }
}
