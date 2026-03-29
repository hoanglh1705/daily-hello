import 'dart:async';

import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';
import 'package:provider/provider.dart';
import 'attendance_controller.dart';
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
    final isPermissionDeniedForever =
        message.toLowerCase().contains('vinh vien');
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
                Text(
                  message,
                  style: theme.textTheme.bodyMedium,
                ),
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
      await _showLocationErrorSheet(
        message,
        onRetry: _handleCheckIn,
      );
      return;
    }
    _showActionErrorSnackBar(
      message: message,
      onRetry: _handleCheckIn,
    );
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
      await _showLocationErrorSheet(
        message,
        onRetry: _handleCheckOut,
      );
      return;
    }
    _showActionErrorSnackBar(
      message: message,
      onRetry: _handleCheckOut,
    );
  }

  @override
  void initState() {
    super.initState();
    _clockTimer = Timer.periodic(const Duration(seconds: 1), (_) {
      setState(() => _now = DateTime.now());
    });
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final controller = context.read<AttendanceController>();
      controller.loadTodayAttendance();
      controller.loadWifiInfo();
      controller.loadHistory(refresh: true);
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
      appBar: AppBar(
        title: const Text('Chấm công'),
        centerTitle: true,
      ),
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
                border: Border.all(
                  color: primaryColor.withAlpha(30),
                ),
              ),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 28),
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
                            horizontal: 14, vertical: 8),
                        decoration: BoxDecoration(
                          color: Colors.grey[100],
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.wifi,
                                size: 16, color: primaryColor),
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

                    // Action buttons (always visible)
                    Row(
                      children: [
                        Expanded(
                          child: _ActionButton(
                            label: 'Check-in',
                            icon: Icons.login,
                            isActive: true,
                            isLoading: controller.isLoading && today == null,
                            primaryColor: primaryColor,
                            onPressed: !controller.isLoading
                                ? _handleCheckIn
                                : null,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: _ActionButton(
                            label: 'Check-out',
                            icon: Icons.logout,
                            isActive: today != null,
                            isLoading: controller.isLoading &&
                                today != null,
                            primaryColor: primaryColor,
                            onPressed: today != null && !controller.isLoading
                                ? _handleCheckOut
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

            // Recent activity list
            if (controller.history.isEmpty &&
                !controller.isLoadingHistory)
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
              ...controller.history.take(4).map(
                    (item) => _RecentActivityItem(
                      item: item,
                      primaryColor: primaryColor,
                    ),
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
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(10),
        ),
      ),
      child: isLoading
          ? SizedBox(
              width: 18,
              height: 18,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                color: color,
              ),
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

class _RecentActivityItem extends StatelessWidget {
  final dynamic item;
  final Color primaryColor;

  const _RecentActivityItem({
    required this.item,
    required this.primaryColor,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final checkInTime = DateFormatUtils.formatTime(item.checkIn);
    final checkOutTime = item.checkOut != null
        ? DateFormatUtils.formatTime(item.checkOut!)
        : null;
    final dateLabel = DateFormatUtils.formatVietnameseDateShort(item.checkIn);

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
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
        children: [
          // Check-in row
          Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: primaryColor.withAlpha(20),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Icon(
                  Icons.login,
                  color: primaryColor,
                  size: 20,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Check-in sáng',
                      style: TextStyle(
                        fontWeight: FontWeight.w600,
                        fontSize: 14,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      dateLabel,
                      style: TextStyle(
                        fontSize: 12,
                        color: Colors.grey[500],
                      ),
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    checkInTime,
                    style: const TextStyle(
                      fontWeight: FontWeight.w700,
                      fontSize: 15,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    _getCheckInStatus(item),
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                      color: _getStatusColor(item),
                    ),
                  ),
                ],
              ),
            ],
          ),
          if (checkOutTime != null) ...[
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 8),
              child: Divider(height: 1, color: Colors.grey[200]),
            ),
            // Check-out row
            Row(
              children: [
                Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: Colors.green.withAlpha(20),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: const Icon(
                    Icons.logout,
                    color: Colors.green,
                    size: 20,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Check-out chiều',
                        style: TextStyle(
                          fontWeight: FontWeight.w600,
                          fontSize: 14,
                        ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        dateLabel,
                        style: TextStyle(
                          fontSize: 12,
                          color: Colors.grey[500],
                        ),
                      ),
                    ],
                  ),
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      checkOutTime,
                      style: const TextStyle(
                        fontWeight: FontWeight.w700,
                        fontSize: 15,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      'HOÀN THÀNH',
                      style: TextStyle(
                        fontSize: 11,
                        fontWeight: FontWeight.w600,
                        color: Colors.green[600],
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ],
        ],
      ),
    );
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
}
