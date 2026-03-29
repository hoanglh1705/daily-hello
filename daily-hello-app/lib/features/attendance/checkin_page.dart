import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';
import 'package:provider/provider.dart';
import 'package:wifi_signal_strength_indicator/wifi_signal_strength_indicator.dart';
import 'attendance_controller.dart';
import '../../widgets/app_button.dart';
import '../../widgets/app_card.dart';
import '../../core/utils/date_format_utils.dart';

class CheckInPage extends StatefulWidget {
  const CheckInPage({super.key});

  @override
  State<CheckInPage> createState() => _CheckInPageState();
}

class _CheckInPageState extends State<CheckInPage> {
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
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final controller = context.read<AttendanceController>();
      controller.loadTodayAttendance();
      controller.loadWifiInfo();
    });
  }

  @override
  Widget build(BuildContext context) {
    final controller = context.watch<AttendanceController>();
    final today = controller.todayAttendance;
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Chấm công'),
        centerTitle: true,
      ),
      body: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            AppCard(
              child: Column(
                children: [
                  Text(
                    DateFormatUtils.formatDate(DateTime.now()),
                    style: theme.textTheme.titleMedium
                        ?.copyWith(color: Colors.grey[600]),
                  ),
                  const SizedBox(height: 8),
                  _StatusBadge(today: today),
                  const SizedBox(height: 16),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      _TimeInfo(
                        label: 'Check-in',
                        time: today != null
                            ? DateFormatUtils.formatTime(today.checkIn)
                            : '--:--',
                        color: Colors.green,
                      ),
                      Container(
                          width: 1, height: 40, color: Colors.grey[300]),
                      _TimeInfo(
                        label: 'Check-out',
                        time: today?.checkOut != null
                            ? DateFormatUtils.formatTime(today!.checkOut!)
                            : '--:--',
                        color: Colors.red,
                      ),
                      Container(
                          width: 1, height: 40, color: Colors.grey[300]),
                      _TimeInfo(
                        label: 'Tổng giờ',
                        time: DateFormatUtils.formatDuration(
                          today?.checkIn ?? DateTime.now(),
                          today?.checkOut,
                        ),
                        color: theme.colorScheme.primary,
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),
            _WifiInfoCard(
              ssid: controller.wifiSsid,
              signalStrength: controller.wifiSignalStrength,
              signalLevel: controller.wifiSignalLevel,
              onRefresh: () => controller.loadWifiInfo(),
            ),
            const SizedBox(height: 24),
            if (today == null || !today.isCheckedOut) ...[
              if (today == null)
                AppButton(
                  label: 'Check In',
                  icon: Icons.login,
                  isLoading: controller.isLoading,
                  onPressed: controller.isLoading ? null : _handleCheckIn,
                ),
              if (today != null && !today.isCheckedOut)
                AppButton(
                  label: 'Check Out',
                  icon: Icons.logout,
                  color: Colors.red,
                  isLoading: controller.isLoading,
                  onPressed: controller.isLoading ? null : _handleCheckOut,
                ),
            ] else
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.green[50],
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: Colors.green),
                ),
                child: const Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(Icons.check_circle, color: Colors.green),
                    SizedBox(width: 8),
                    Text('Đã hoàn thành chấm công hôm nay',
                        style: TextStyle(color: Colors.green)),
                  ],
                ),
              ),
          ],
        ),
      ),
    );
  }
}

class _WifiInfoCard extends StatelessWidget {
  final String? ssid;
  final int? signalStrength;
  final int? signalLevel;
  final VoidCallback onRefresh;

  const _WifiInfoCard({
    required this.ssid,
    required this.signalStrength,
    required this.signalLevel,
    required this.onRefresh,
  });

  String _signalLabel(int? level) {
    switch (level) {
      case 4:
        return 'Rất mạnh';
      case 3:
        return 'Mạnh';
      case 2:
        return 'Trung bình';
      case 1:
        return 'Yếu';
      default:
        return 'Không xác định';
    }
  }

  Color _signalColor(int? level) {
    switch (level) {
      case 4:
        return Colors.green;
      case 3:
        return Colors.lightGreen;
      case 2:
        return Colors.orange;
      case 1:
        return Colors.deepOrange;
      default:
        return Colors.grey;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isConnected = ssid != null && ssid!.isNotEmpty;

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(
                isConnected ? Icons.wifi : Icons.wifi_off,
                color: isConnected ? _signalColor(signalLevel) : Colors.grey,
                size: 20,
              ),
              const SizedBox(width: 8),
              Text(
                'Thông tin WiFi',
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const Spacer(),
              InkWell(
                onTap: onRefresh,
                borderRadius: BorderRadius.circular(20),
                child: const Padding(
                  padding: EdgeInsets.all(4),
                  child: Icon(Icons.refresh, size: 18, color: Colors.grey),
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          if (!isConnected)
            Text(
              'Không kết nối WiFi',
              style: theme.textTheme.bodyMedium?.copyWith(color: Colors.grey),
            )
          else ...[
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Tên WiFi',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: Colors.grey[600],
                        ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        ssid ?? '--',
                        style: theme.textTheme.bodyMedium?.copyWith(
                          fontWeight: FontWeight.w600,
                        ),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      'Cường độ',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: Colors.grey[600],
                      ),
                    ),
                    const SizedBox(height: 2),
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        if (signalStrength != null)
                          WifiSignalStrengthIndicator(
                            rssi: signalStrength,
                            style: WifiSignalStyle.bars,
                            size: 18,
                          ),
                        const SizedBox(width: 6),
                        Text(
                          _signalLabel(signalLevel),
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: _signalColor(signalLevel),
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ],
            ),
            if (signalStrength != null) ...[
              const SizedBox(height: 8),
              Text(
                '$signalStrength dBm',
                style: theme.textTheme.bodySmall?.copyWith(
                  color: Colors.grey[500],
                ),
              ),
            ],
          ],
        ],
      ),
    );
  }
}

class _StatusBadge extends StatelessWidget {
  final dynamic today;
  const _StatusBadge({this.today});

  @override
  Widget build(BuildContext context) {
    if (today == null) {
      return _badge('Chưa check-in', Colors.grey);
    } else if (!today.isCheckedOut) {
      return _badge('Đang làm việc', Colors.green);
    } else {
      return _badge('Đã check-out', Colors.blue);
    }
  }

  Widget _badge(String label, Color color) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
      decoration: BoxDecoration(
        color: color.withAlpha(30),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: color),
      ),
      child: Text(label, style: TextStyle(color: color, fontWeight: FontWeight.w600)),
    );
  }
}

class _TimeInfo extends StatelessWidget {
  final String label;
  final String time;
  final Color color;

  const _TimeInfo({
    required this.label,
    required this.time,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text(time,
            style: TextStyle(
                fontSize: 20, fontWeight: FontWeight.bold, color: color)),
        Text(label,
            style: const TextStyle(fontSize: 12, color: Colors.grey)),
      ],
    );
  }
}
