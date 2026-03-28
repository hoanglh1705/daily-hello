import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
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
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<AttendanceController>().loadTodayAttendance();
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
            const SizedBox(height: 24),
            if (controller.errorMessage != null)
              Padding(
                padding: const EdgeInsets.only(bottom: 12),
                child: Text(
                  controller.errorMessage!,
                  style: TextStyle(color: theme.colorScheme.error),
                  textAlign: TextAlign.center,
                ),
              ),
            if (today == null || !today.isCheckedOut) ...[
              if (today == null)
                AppButton(
                  label: 'Check In',
                  icon: Icons.login,
                  isLoading: controller.isLoading,
                  onPressed: controller.isLoading
                      ? null
                      : () async {
                          final ok =
                              await context.read<AttendanceController>().checkIn();
                          if (ok && context.mounted) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(
                                  content: Text('Check-in thành công!'),
                                  backgroundColor: Colors.green),
                            );
                          }
                        },
                ),
              if (today != null && !today.isCheckedOut)
                AppButton(
                  label: 'Check Out',
                  icon: Icons.logout,
                  color: Colors.red,
                  isLoading: controller.isLoading,
                  onPressed: controller.isLoading
                      ? null
                      : () async {
                          final ok =
                              await context.read<AttendanceController>().checkOut();
                          if (ok && context.mounted) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(
                                  content: Text('Check-out thành công!'),
                                  backgroundColor: Colors.orange),
                            );
                          }
                        },
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
