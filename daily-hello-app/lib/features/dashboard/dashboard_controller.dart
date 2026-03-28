import 'package:flutter/material.dart';
import '../../models/attendance.dart';
import '../../services/attendance_service.dart';

class DashboardController extends ChangeNotifier {
  final AttendanceService _service;

  Attendance? todayAttendance;
  List<Attendance> recentHistory = [];
  bool isLoading = false;

  DashboardController(this._service);

  Future<void> loadDashboard() async {
    isLoading = true;
    notifyListeners();

    try {
      final results = await Future.wait([
        _service.getTodayAttendance(),
        _service.getHistory(page: 1, limit: 5),
      ]);
      todayAttendance = results[0] as Attendance?;
      recentHistory = results[1] as List<Attendance>;
    } catch (_) {}

    isLoading = false;
    notifyListeners();
  }

  int get presentDays => recentHistory.where((a) => a.isCheckedOut).length;
}
