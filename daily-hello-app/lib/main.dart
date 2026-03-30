import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'app.dart';
import 'core/network/api_client.dart';
import 'core/storage/secure_storage.dart';
import 'features/auth/auth_controller.dart';
import 'features/attendance/attendance_controller.dart';
import 'features/device/device_controller.dart';
import 'features/profile/profile_controller.dart';
import 'services/auth_service.dart';
import 'services/attendance_service.dart';
import 'services/branch_service.dart';
import 'services/device_service.dart';

final navigatorKey = GlobalKey<NavigatorState>();

void main() {
  runApp(const AppProviders());
}

class AppProviders extends StatelessWidget {
  const AppProviders({super.key});

  @override
  Widget build(BuildContext context) {
    final secureStorage = SecureStorage();
    final apiClient = ApiClient(
      secureStorage,
      onUnauthorized: () {
        navigatorKey.currentState
            ?.pushNamedAndRemoveUntil('/login', (_) => false);
      },
    );
    final authService = AuthService(apiClient.dio);
    final attendanceService = AttendanceService(apiClient.dio);
    final branchService = BranchService(apiClient.dio);
    final deviceService = DeviceService(apiClient.dio);

    return MultiProvider(
      providers: [
        Provider<SecureStorage>.value(value: secureStorage),
        Provider<ApiClient>.value(value: apiClient),
        Provider<AuthService>.value(value: authService),
        Provider<AttendanceService>.value(value: attendanceService),
        Provider<BranchService>.value(value: branchService),
        Provider<DeviceService>.value(value: deviceService),
        ChangeNotifierProvider(
          create: (_) => AuthController(authService, secureStorage),
        ),
        ChangeNotifierProvider(
          create: (_) => AttendanceController(attendanceService, branchService),
        ),
        ChangeNotifierProvider(
          create: (_) => ProfileController(authService, secureStorage),
        ),
        ChangeNotifierProvider(
          create: (_) => DeviceController(deviceService),
        ),
      ],
      child: const DailyHelloApp(),
    );
  }
}
