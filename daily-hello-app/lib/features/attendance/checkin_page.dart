import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';
import 'package:image_picker/image_picker.dart';
import 'package:provider/provider.dart';
import 'attendance_controller.dart';
import '../auth/auth_controller.dart';
import '../../widgets/app_button.dart';
import '../../core/utils/date_format_utils.dart';

class CheckInPage extends StatefulWidget {
  const CheckInPage({super.key});

  @override
  State<CheckInPage> createState() => _CheckInPageState();
}

class _CheckInPageState extends State<CheckInPage> {
  late Timer _clockTimer;
  DateTime _now = DateTime.now();

  bool _isLocationError(String message) {
    final normalized = message.toLowerCase();
    return normalized.contains('vi tri') ||
        normalized.contains('gps') ||
        normalized.contains('location');
  }

  Future<void> _showLocationErrorSheet(
    String message, {
    required Future<void> Function() onRetry,
  }) async {
    final isPermissionDeniedForever = message.toLowerCase().contains(
      'vinh vien',
    );
    final isGpsDisabled = message.toLowerCase().contains('gps');

    await showModalBottomSheet<void>(
      context: context,
      showDragHandle: true,
      builder: (sheetContext) {
        final theme = Theme.of(sheetContext);
        return SafeArea(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(24, 8, 24, 24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(
                      Icons.location_off_outlined,
                      color: theme.colorScheme.error,
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Text(
                        'Can quyen vi tri',
                        style: theme.textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Text(message, style: theme.textTheme.bodyMedium),
                const SizedBox(height: 20),
                SizedBox(
                  width: double.infinity,
                  child: AppButton(
                    label: isPermissionDeniedForever
                        ? 'Mo cai dat'
                        : isGpsDisabled
                        ? 'Mo GPS'
                        : 'Thu lai',
                    icon: isPermissionDeniedForever
                        ? Icons.settings_outlined
                        : isGpsDisabled
                        ? Icons.my_location
                        : Icons.refresh,
                    onPressed: () async {
                      Navigator.pop(sheetContext);
                      if (isPermissionDeniedForever) {
                        await Geolocator.openAppSettings();
                        return;
                      }
                      if (isGpsDisabled) {
                        await Geolocator.openLocationSettings();
                        return;
                      }
                      await onRetry();
                    },
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  void _showActionErrorSnackBar({
    required String message,
    required Future<void> Function() onRetry,
  }) {
    ScaffoldMessenger.of(context)
      ..hideCurrentSnackBar()
      ..showSnackBar(
        SnackBar(
          content: Text(message),
          behavior: SnackBarBehavior.floating,
          action: SnackBarAction(
            label: 'Thu lai',
            onPressed: () {
              onRetry();
            },
          ),
        ),
      );
  }

  Future<void> _handleCheckIn() async {
    final controller = context.read<AttendanceController>();
    final ok = await controller.checkIn();
    if (!mounted) return;

    if (ok) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Check-in thanh cong!'),
          backgroundColor: Colors.green,
        ),
      );
      return;
    }

    final message = controller.errorMessage ?? 'Check-in that bai.';
    if (_isLocationError(message)) {
      await _showLocationErrorSheet(message, onRetry: _handleCheckIn);
      return;
    }
    _showActionErrorSnackBar(message: message, onRetry: _handleCheckIn);
  }

  Future<void> _handleCheckOut() async {
    final controller = context.read<AttendanceController>();
    final ok = await controller.checkOut();
    if (!mounted) return;

    if (ok) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Check-out thanh cong!'),
          backgroundColor: Colors.orange,
        ),
      );
      return;
    }

    final message = controller.errorMessage ?? 'Check-out that bai.';
    if (_isLocationError(message)) {
      await _showLocationErrorSheet(message, onRetry: _handleCheckOut);
      return;
    }
    _showActionErrorSnackBar(message: message, onRetry: _handleCheckOut);
  }

  Future<void> _handleCheckInGps() async {
    final controller = context.read<AttendanceController>();

    final picker = ImagePicker();
    final file = await picker.pickImage(
      source: ImageSource.camera,
      maxWidth: 800,
      imageQuality: 75,
    );
    if (file == null) return;

    final bytes = await file.readAsBytes();
    final base64Image = 'data:image/jpeg;base64,${base64Encode(bytes)}';

    final ok = await controller.checkInGps(base64Image);
    if (!mounted) return;

    if (ok) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Check-in GPS chờ duyệt!'),
          backgroundColor: Colors.green,
        ),
      );
      return;
    }

    final message = controller.errorMessage ?? 'Check-in GPS thất bại.';
    if (_isLocationError(message)) {
      await _showLocationErrorSheet(message, onRetry: _handleCheckInGps);
      return;
    }
    _showActionErrorSnackBar(message: message, onRetry: _handleCheckInGps);
  }

  Future<void> _handleCheckOutGps() async {
    final controller = context.read<AttendanceController>();

    final picker = ImagePicker();
    final file = await picker.pickImage(
      source: ImageSource.camera,
      maxWidth: 800,
      imageQuality: 75,
    );
    if (file == null) return;

    final bytes = await file.readAsBytes();
    final base64Image = 'data:image/jpeg;base64,${base64Encode(bytes)}';

    final ok = await controller.checkOutGps(base64Image);
    if (!mounted) return;

    if (ok) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Check-out GPS chờ duyệt!'),
          backgroundColor: Colors.orange,
        ),
      );
      return;
    }

    final message = controller.errorMessage ?? 'Check-out GPS thất bại.';
    if (_isLocationError(message)) {
      await _showLocationErrorSheet(message, onRetry: _handleCheckOutGps);
      return;
    }
    _showActionErrorSnackBar(message: message, onRetry: _handleCheckOutGps);
  }

  List<Widget> _buildGroupedHistory(
    List<dynamic> items,
    Color primaryColor,
    ThemeData theme,
  ) {
    // Group items by date
    final grouped = <String, List<dynamic>>{};
    for (final item in items) {
      final dateKey = DateFormatUtils.formatVietnameseDateShort(item.checkIn);
      grouped.putIfAbsent(dateKey, () => []).add(item);
    }

    final widgets = <Widget>[];
    for (final entry in grouped.entries) {
      widgets.add(
        Container(
          margin: const EdgeInsets.only(bottom: 12),
          decoration: BoxDecoration(
            color: theme.colorScheme.surface,
            borderRadius: BorderRadius.circular(14),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withAlpha(10),
                blurRadius: 8,
                offset: const Offset(0, 2),
              ),
            ],
            border: Border.all(color: Colors.grey[200]!),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Date header
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 14, 16, 8),
                child: Text(
                  entry.key,
                  style: TextStyle(
                    fontSize: 13,
                    fontWeight: FontWeight.w600,
                    color: Colors.grey[600],
                  ),
                ),
              ),
              // Activity rows for this date
              ...entry.value.expand<Widget>((item) {
                final rows = <Widget>[
                  _CompactActivityRow(
                    icon: Icons.login,
                    iconColor: primaryColor,
                    label: 'Check-in',
                    time: DateFormatUtils.formatTime(item.checkIn),
                    status: _getCheckInStatus(item),
                    statusColor: _getStatusColor(item),
                  ),
                ];
                if (item.checkOut != null) {
                  rows.add(
                    _CompactActivityRow(
                      icon: Icons.logout,
                      iconColor: Colors.green,
                      label: 'Check-out',
                      time: DateFormatUtils.formatTime(item.checkOut!),
                      status: 'HOÀN THÀNH',
                      statusColor: Colors.green[600]!,
                    ),
                  );
                }
                return rows;
              }),
              const SizedBox(height: 8),
            ],
          ),
        ),
      );
    }
    return widgets;
  }

  String _getCheckInStatus(dynamic item) {
    final hour = item.checkIn.hour;
    if (hour < 8) return 'SỚM';
    if (hour == 8 && item.checkIn.minute <= 15) return 'ĐÚNG GIỜ';
    return 'MUỘN';
  }

  Color _getStatusColor(dynamic item) {
    final hour = item.checkIn.hour;
    if (hour < 8) return Colors.blue;
    if (hour == 8 && item.checkIn.minute <= 15) return Colors.green[600]!;
    return Colors.orange[700]!;
  }

  @override
  void initState() {
    super.initState();
    _clockTimer = Timer.periodic(const Duration(seconds: 1), (_) {
      setState(() => _now = DateTime.now());
    });
    WidgetsBinding.instance.addPostFrameCallback((_) async {
      final controller = context.read<AttendanceController>();
      final authController = context.read<AuthController>();
      controller.loadTodayAttendance();
      await controller.loadWifiInfo();
      controller.loadHistory(refresh: true);
      // Load profile then validate branch wifi
      await authController.loadProfile();
      if (mounted) {
        await controller.loadBranchWifi(authController.currentUser?.branchId);
      }
    });
  }

  @override
  void dispose() {
    _clockTimer.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final controller = context.watch<AttendanceController>();
    final today = controller.todayAttendance;
    final theme = Theme.of(context);
    final primaryColor = theme.colorScheme.primary;

    return Scaffold(
      appBar: AppBar(title: const Text('Chấm công'), centerTitle: true),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            // Main attendance card
            Container(
              width: double.infinity,
              decoration: BoxDecoration(
                color: theme.colorScheme.surface,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [
                  BoxShadow(
                    color: primaryColor.withAlpha(25),
                    blurRadius: 20,
                    offset: const Offset(0, 4),
                  ),
                ],
                border: Border.all(color: primaryColor.withAlpha(30)),
              ),
              child: Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 24,
                  vertical: 28,
                ),
                child: Column(
                  children: [
                    // Date
                    Text(
                      DateFormatUtils.formatVietnameseDate(_now),
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: Colors.grey[500],
                        fontWeight: FontWeight.w600,
                        letterSpacing: 1.2,
                      ),
                    ),
                    const SizedBox(height: 12),

                    // Live clock
                    _LiveClock(now: _now, primaryColor: primaryColor),
                    const SizedBox(height: 12),

                    // WiFi location badge
                    if (controller.wifiSsid != null &&
                        controller.wifiSsid!.isNotEmpty)
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 14,
                          vertical: 8,
                        ),
                        decoration: BoxDecoration(
                          color: Colors.grey[100],
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.wifi, size: 16, color: primaryColor),
                            const SizedBox(width: 6),
                            Text(
                              controller.wifiSsid!,
                              style: theme.textTheme.bodySmall?.copyWith(
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ],
                        ),
                      ),
                    const SizedBox(height: 20),

                    // Check-in / Check-out time boxes
                    Row(
                      children: [
                        Expanded(
                          child: _TimeBox(
                            label: 'GIỜ CHECK-IN',
                            time: today != null
                                ? DateFormatUtils.formatTime(today.checkIn)
                                : '--:--',
                            primaryColor: primaryColor,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: _TimeBox(
                            label: 'GIỜ CHECK-OUT',
                            time: today?.checkOut != null
                                ? DateFormatUtils.formatTime(today!.checkOut!)
                                : '--:--',
                            primaryColor: primaryColor,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),

                    // Fraud warning message
                    if (controller.fraudWarning != null) ...[
                      Container(
                        width: double.infinity,
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: Colors.red[50],
                          borderRadius: BorderRadius.circular(10),
                          border: Border.all(color: Colors.red[300]!),
                        ),
                        child: Row(
                          children: [
                            Icon(
                              Icons.gps_off,
                              size: 18,
                              color: Colors.red[700],
                            ),
                            const SizedBox(width: 8),
                            Expanded(
                              child: Text(
                                controller.fraudWarning!,
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.red[700],
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 12),
                    ],

                    // WiFi error message
                    if (controller.isWifiChecked &&
                        !controller.isWifiValid &&
                        controller.wifiErrorMessage != null) ...[
                      Container(
                        width: double.infinity,
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: Colors.red[50],
                          borderRadius: BorderRadius.circular(10),
                          border: Border.all(color: Colors.red[200]!),
                        ),
                        child: Row(
                          children: [
                            Icon(
                              Icons.wifi_off,
                              size: 18,
                              color: Colors.red[600],
                            ),
                            const SizedBox(width: 8),
                            Expanded(
                              child: Text(
                                controller.wifiErrorMessage!,
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.red[700],
                                ),
                              ),
                            ),
                            GestureDetector(
                              onTap: () => controller.refreshWifiValidation(),
                              child: Icon(
                                Icons.refresh,
                                size: 18,
                                color: Colors.red[400],
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 16),
                    ],

                    // Action buttons (always visible)
                    Row(
                      children: [
                        Expanded(
                          child: _ActionButton(
                            label: 'Check-in',
                            icon: Icons.login,
                            isActive:
                                !controller.isWifiChecked ||
                                controller.isWifiValid,
                            isLoading: controller.isLoading && today == null,
                            primaryColor: primaryColor,
                            onPressed:
                                (!controller.isWifiChecked ||
                                        controller.isWifiValid) &&
                                    !controller.isLoading
                                ? _handleCheckIn
                                : null,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: _ActionButton(
                            label: 'Check-out',
                            icon: Icons.logout,
                            isActive:
                                today != null &&
                                (!controller.isWifiChecked ||
                                    controller.isWifiValid),
                            isLoading: controller.isLoading && today != null,
                            primaryColor: primaryColor,
                            onPressed:
                                today != null &&
                                    (!controller.isWifiChecked ||
                                        controller.isWifiValid) &&
                                    !controller.isLoading
                                ? _handleCheckOut
                                : null,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        Expanded(
                          child: _ActionButton(
                            label: 'GPS In',
                            icon: Icons.camera_alt,
                            isActive: true,
                            isLoading: controller.isLoading,
                            primaryColor: Colors.blue,
                            onPressed: !controller.isLoading
                                ? _handleCheckInGps
                                : null,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: _ActionButton(
                            label: 'GPS Out',
                            icon: Icons.camera_enhance,
                            isActive:
                                true, // as backend allows checkout without checkin now
                            isLoading: controller.isLoading,
                            primaryColor: Colors.blue,
                            onPressed: !controller.isLoading
                                ? _handleCheckOutGps
                                : null,
                          ),
                        ),
                      ],
                    ),
                    if (today != null && today.isCheckedOut) ...[
                      const SizedBox(height: 14),
                      _CompletedBadge(primaryColor: primaryColor),
                    ],
                  ],
                ),
              ),
            ),
            const SizedBox(height: 28),

            // Recent activity header
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Hoạt động gần đây',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w700,
                  ),
                ),
                GestureDetector(
                  onTap: () {
                    Navigator.pushNamed(context, '/history');
                  },
                  child: Text(
                    'Xem tất cả',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: primaryColor,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // Recent activity list (grouped by date)
            if (controller.history.isEmpty && !controller.isLoadingHistory)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 24),
                child: Text(
                  'Chưa có hoạt động',
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: Colors.grey,
                  ),
                ),
              )
            else
              ..._buildGroupedHistory(
                controller.history.take(4).toList(),
                primaryColor,
                theme,
              ),
            if (controller.isLoadingHistory)
              const Padding(
                padding: EdgeInsets.symmetric(vertical: 16),
                child: Center(
                  child: SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}

class _LiveClock extends StatelessWidget {
  final DateTime now;
  final Color primaryColor;

  const _LiveClock({required this.now, required this.primaryColor});

  @override
  Widget build(BuildContext context) {
    final hours = now.hour.toString().padLeft(2, '0');
    final minutes = now.minute.toString().padLeft(2, '0');
    final seconds = now.second.toString().padLeft(2, '0');

    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.baseline,
      textBaseline: TextBaseline.alphabetic,
      children: [
        Text(
          '$hours:$minutes',
          style: TextStyle(
            fontSize: 56,
            fontWeight: FontWeight.w700,
            color: primaryColor,
            height: 1,
          ),
        ),
        Text(
          ':$seconds',
          style: TextStyle(
            fontSize: 32,
            fontWeight: FontWeight.w500,
            color: primaryColor.withAlpha(150),
            height: 1,
          ),
        ),
      ],
    );
  }
}

class _TimeBox extends StatelessWidget {
  final String label;
  final String time;
  final Color primaryColor;

  const _TimeBox({
    required this.label,
    required this.time,
    required this.primaryColor,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.grey[300]!),
      ),
      child: Column(
        children: [
          Text(
            label,
            style: TextStyle(
              fontSize: 10,
              fontWeight: FontWeight.w600,
              color: Colors.grey[500],
              letterSpacing: 0.8,
            ),
          ),
          const SizedBox(height: 6),
          Text(
            time,
            style: TextStyle(
              fontSize: 26,
              fontWeight: FontWeight.w700,
              color: primaryColor,
            ),
          ),
        ],
      ),
    );
  }
}

class _ActionButton extends StatelessWidget {
  final String label;
  final IconData icon;
  final bool isActive;
  final bool isLoading;
  final Color primaryColor;
  final VoidCallback? onPressed;

  const _ActionButton({
    required this.label,
    required this.icon,
    required this.isActive,
    required this.isLoading,
    required this.primaryColor,
    this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    final color = isActive ? primaryColor : Colors.grey[400]!;

    return OutlinedButton(
      onPressed: isLoading ? null : onPressed,
      style: OutlinedButton.styleFrom(
        foregroundColor: color,
        side: BorderSide(color: color),
        padding: const EdgeInsets.symmetric(vertical: 12),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
      ),
      child: isLoading
          ? SizedBox(
              width: 18,
              height: 18,
              child: CircularProgressIndicator(strokeWidth: 2, color: color),
            )
          : Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(icon, size: 18),
                const SizedBox(width: 6),
                Text(
                  label,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
    );
  }
}

class _CompletedBadge extends StatelessWidget {
  final Color primaryColor;

  const _CompletedBadge({required this.primaryColor});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Icon(Icons.check_circle, color: Colors.green[600], size: 22),
        const SizedBox(width: 8),
        Text(
          'Đã hoàn thành check-out',
          style: TextStyle(
            color: Colors.green[600],
            fontWeight: FontWeight.w600,
            fontSize: 14,
          ),
        ),
      ],
    );
  }
}

class _CompactActivityRow extends StatelessWidget {
  final IconData icon;
  final Color iconColor;
  final String label;
  final String time;
  final String status;
  final Color statusColor;

  const _CompactActivityRow({
    required this.icon,
    required this.iconColor,
    required this.label,
    required this.time,
    required this.status,
    required this.statusColor,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
      child: Row(
        children: [
          Icon(icon, size: 18, color: iconColor),
          const SizedBox(width: 10),
          Expanded(
            child: Text(
              label,
              style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w500),
            ),
          ),
          Text(
            time,
            style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 15),
          ),
          const SizedBox(width: 10),
          SizedBox(
            width: 72,
            child: Text(
              status,
              textAlign: TextAlign.end,
              style: TextStyle(
                fontSize: 11,
                fontWeight: FontWeight.w600,
                color: statusColor,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
