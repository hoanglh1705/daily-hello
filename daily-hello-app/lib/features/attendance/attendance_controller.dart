import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';
import 'package:network_info_plus/network_info_plus.dart';
import 'package:wifi_signal_strength_indicator/wifi_signal_strength.dart';
import '../../core/network/api_response.dart';
import '../../core/utils/gps_fraud_detector.dart';
import '../../core/utils/location_permission_utils.dart';
import '../../models/attendance.dart';
import '../../models/branch_wifi.dart';
import '../../services/attendance_service.dart';
import '../../services/branch_service.dart';

class AttendanceController extends ChangeNotifier {
  final AttendanceService _service;
  final BranchService _branchService;
  final NetworkInfo _networkInfo = NetworkInfo();

  Attendance? todayAttendance;
  List<Attendance> history = [];
  bool isLoading = false;
  bool isLoadingHistory = false;
  String? errorMessage;
  String? fraudWarning;
  int currentPage = 1;
  bool hasMore = true;

  // WiFi info state
  String? wifiSsid;
  int? wifiSignalStrength; // dBm
  int? wifiSignalLevel; // 0-4

  // Branch WiFi validation
  List<BranchWifi> _branchWifiList = [];
  bool _wifiValidated = false;
  bool _wifiMatched = false;
  String? wifiErrorMessage;

  bool get isWifiValid => _wifiMatched;
  bool get isWifiChecked => _wifiValidated;
  List<BranchWifi> get branchWifiList => _branchWifiList;

  AttendanceController(this._service, this._branchService);

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

  Future<void> loadBranchWifi(String? branchId) async {
    if (branchId == null || branchId.isEmpty) {
      _wifiValidated = true;
      _wifiMatched = false;
      wifiErrorMessage = 'Tài khoản chưa được gán chi nhánh';
      notifyListeners();
      return;
    }

    try {
      _branchWifiList = await _branchService.getBranchWifiList(branchId);
      _wifiValidated = true;

      if (_branchWifiList.isEmpty) {
        // No wifi configured for branch — allow check-in
        _wifiMatched = true;
        wifiErrorMessage = null;
      } else {
        _validateWifi();
      }
    } catch (error) {
      _wifiValidated = true;
      _wifiMatched = false;
      wifiErrorMessage = 'Không thể tải danh sách WiFi chi nhánh';
    }
    notifyListeners();
  }

  void _validateWifi() {
    if (_branchWifiList.isEmpty) {
      _wifiMatched = true;
      wifiErrorMessage = null;
      return;
    }

    final currentSsid = wifiSsid?.toLowerCase();
    if (currentSsid == null || currentSsid.isEmpty) {
      _wifiMatched = false;
      wifiErrorMessage = 'Không kết nối WiFi. Vui lòng kết nối WiFi của chi nhánh';
      return;
    }

    final matched = _branchWifiList.any((bw) {
      final bwSsid = bw.ssid?.toLowerCase();
      return bwSsid != null && bwSsid.isNotEmpty && bwSsid == currentSsid;
    });

    _wifiMatched = matched;
    if (!matched) {
      final allowedNames = _branchWifiList
          .where((bw) => bw.ssid != null && bw.ssid!.isNotEmpty)
          .map((bw) => bw.ssid!)
          .join(', ');
      wifiErrorMessage =
          'WiFi "$wifiSsid" không thuộc chi nhánh.\nWiFi hợp lệ: $allowedNames';
    } else {
      wifiErrorMessage = null;
    }
  }

  Future<void> refreshWifiValidation() async {
    await loadWifiInfo();
    _validateWifi();
    notifyListeners();
  }

  Future<bool> checkIn() async {
    isLoading = true;
    errorMessage = null;
    fraudWarning = null;
    notifyListeners();

    try {
      final position = await _getPosition();

      // Fraud detection (GPS spoof, VPN, root/jailbreak)
      final fraudResult = await GpsFraudDetector.detect(position);
      if (fraudResult.isFraudulent) {
        fraudWarning = fraudResult.reason;
        errorMessage = 'Phát hiện gian lận. ${fraudResult.reason}';
        return false;
      }

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
    fraudWarning = null;
    notifyListeners();

    try {
      final position = await _getPosition();

      // Fraud detection (GPS spoof, VPN, root/jailbreak)
      final fraudResult = await GpsFraudDetector.detect(position);
      if (fraudResult.isFraudulent) {
        fraudWarning = fraudResult.reason;
        errorMessage = 'Phát hiện gian lận. ${fraudResult.reason}';
        return false;
      }

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

  Future<bool> checkInGps(String imageBase64) async {
    isLoading = true;
    errorMessage = null;
    fraudWarning = null;
    notifyListeners();

    try {
      final position = await _getPosition();

      final fraudResult = await GpsFraudDetector.detect(position);
      if (fraudResult.isFraudulent) {
        fraudWarning = fraudResult.reason;
        errorMessage = 'Phát hiện gian lận. ${fraudResult.reason}';
        return false;
      }

      todayAttendance = await _service.checkInGps(
        lat: position.latitude,
        lng: position.longitude,
        imageBase64: imageBase64,
      );
      return true;
    } catch (error) {
      errorMessage =
          getApiErrorMessage(error) ?? 'Check-in GPS thất bại: ${_formatError(error)}';
      return false;
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> checkOutGps(String imageBase64) async {
    isLoading = true;
    errorMessage = null;
    fraudWarning = null;
    notifyListeners();

    try {
      final position = await _getPosition();

      final fraudResult = await GpsFraudDetector.detect(position);
      if (fraudResult.isFraudulent) {
        fraudWarning = fraudResult.reason;
        errorMessage = 'Phát hiện gian lận. ${fraudResult.reason}';
        return false;
      }

      todayAttendance = await _service.checkOutGps(
        lat: position.latitude,
        lng: position.longitude,
        imageBase64: imageBase64,
      );
      return true;
    } catch (error) {
      errorMessage = getApiErrorMessage(error) ??
          'Check-out GPS thất bại: ${_formatError(error)}';
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
