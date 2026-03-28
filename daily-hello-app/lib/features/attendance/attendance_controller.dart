import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';
import 'package:network_info_plus/network_info_plus.dart';
import '../../core/network/api_response.dart';
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

  AttendanceController(this._service);

  Future<void> loadTodayAttendance() async {
    todayAttendance = await _service.getTodayAttendance();
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
          getApiErrorMessage(error) ?? 'Check-in thất bại: ${error.toString()}';
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
          'Check-out thất bại: ${error.toString()}';
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
    if (!hasMore) return;

    isLoadingHistory = true;
    notifyListeners();

    try {
      final result = await _service.getHistory(page: currentPage);
      if (result.isEmpty) {
        hasMore = false;
      } else {
        history.addAll(result);
        currentPage++;
      }
    } catch (error) {
      errorMessage = getApiErrorMessage(error) ??
          'Lỗi tải lịch sử: ${error.toString()}';
    } finally {
      isLoadingHistory = false;
      notifyListeners();
    }
  }

  Future<Position> _getPosition() async {
    bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
    if (!serviceEnabled) throw Exception('GPS chưa được bật');

    LocationPermission permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
      if (permission == LocationPermission.denied) {
        throw Exception('Quyền GPS bị từ chối');
      }
    }
    if (permission == LocationPermission.deniedForever) {
      throw Exception('Quyền GPS bị từ chối vĩnh viễn');
    }

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
