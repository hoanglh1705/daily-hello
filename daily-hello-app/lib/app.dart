import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'features/auth/login_page.dart';
import 'features/auth/auth_controller.dart';
import 'features/attendance/checkin_page.dart';
import 'features/attendance/history_page.dart';
import 'features/device/device_controller.dart';
import 'features/profile/profile_page.dart';
import 'main.dart';

class DailyHelloApp extends StatelessWidget {
  const DailyHelloApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Daily Hello',
      navigatorKey: navigatorKey,
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color(0xFF1976D2),
          brightness: Brightness.light,
        ),
        useMaterial3: true,
        appBarTheme: const AppBarTheme(
          centerTitle: true,
          elevation: 0,
        ),
        cardTheme: const CardThemeData(
          elevation: 2,
        ),
      ),
      home: const _AuthGate(),
      routes: {
        '/login': (_) => const LoginPage(),
        '/home': (_) => const MainShell(),
        '/history': (_) => const HistoryPage(),
      },
    );
  }
}

class _AuthGate extends StatefulWidget {
  const _AuthGate();

  @override
  State<_AuthGate> createState() => _AuthGateState();
}

class _AuthGateState extends State<_AuthGate> {
  @override
  void initState() {
    super.initState();
    _checkAuth();
  }

  Future<void> _checkAuth() async {
    final isLoggedIn =
        await context.read<AuthController>().isLoggedIn();
    if (!mounted) return;
    if (isLoggedIn) {
      Navigator.pushReplacementNamed(context, '/home');
    } else {
      Navigator.pushReplacementNamed(context, '/login');
    }
  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(child: CircularProgressIndicator()),
    );
  }
}

class MainShell extends StatefulWidget {
  const MainShell({super.key});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _selectedIndex = 0;
  bool _deviceChecked = false;

  static const _pages = [
    CheckInPage(),
    ProfilePage(),
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _checkDevice());
  }

  Future<void> _checkDevice() async {
    if (_deviceChecked) return;
    _deviceChecked = true;

    final controller = context.read<DeviceController>();
    final device = await controller.checkDevice();

    if (!mounted) return;

    if (device == null && controller.error == null) {
      // Device not registered — ask user
      _showRegisterDialog();
    } else if (device != null && device.isPending) {
      _showStatusSnackBar('Thiết bị đang chờ admin phê duyệt.');
    } else if (device != null && device.isRejected) {
      _showStatusSnackBar('Thiết bị đã bị từ chối. Vui lòng liên hệ admin.');
    }
  }

  void _showStatusSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        duration: const Duration(seconds: 1),
        showCloseIcon: true,
      ),
    );
  }

  Future<void> _showRegisterDialog() async {
    final confirmed = await showDialog<bool>(
      context: context,
      barrierDismissible: false,
      builder: (ctx) => AlertDialog(
        title: const Text('Đăng ký thiết bị'),
        content: const Text(
          'Thiết bị này chưa được đăng ký.\nBạn có muốn đăng ký thiết bị để sử dụng chấm công không?',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Để sau'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('Đăng ký'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      final controller = context.read<DeviceController>();
      final device = await controller.registerDevice();
      if (!mounted) return;

      if (device != null) {
        _showStatusSnackBar(
          'Đăng ký thành công! Thiết bị đang chờ admin phê duyệt.',
        );
      } else {
        _showStatusSnackBar(controller.error ?? 'Đăng ký thất bại.');
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _pages[_selectedIndex],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _selectedIndex,
        onDestinationSelected: (i) => setState(() => _selectedIndex = i),
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.fingerprint_outlined),
            selectedIcon: Icon(Icons.fingerprint),
            label: 'Chấm công',
          ),
          NavigationDestination(
            icon: Icon(Icons.person_outline),
            selectedIcon: Icon(Icons.person),
            label: 'Hồ sơ',
          ),
        ],
      ),
    );
  }
}
