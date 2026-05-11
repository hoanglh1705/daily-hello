import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'attendance_controller.dart';
import '../../core/utils/date_format_utils.dart';
import '../../models/attendance.dart';
import '../../models/holiday.dart';

class HistoryPage extends StatefulWidget {
  const HistoryPage({super.key});

  @override
  State<HistoryPage> createState() => _HistoryPageState();
}

class _HistoryPageState extends State<HistoryPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<AttendanceController>().loadMonth();
    });
  }

  @override
  Widget build(BuildContext context) {
    final controller = context.watch<AttendanceController>();
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Lịch chấm công'),
        centerTitle: true,
      ),
      body: Column(
        children: [
          // Month navigation header
          _buildMonthHeader(controller, theme),
          // Weekday labels
          _buildWeekdayLabels(theme),
          // Calendar grid
          Expanded(
            child: controller.isLoadingMonth
                ? const Center(child: CircularProgressIndicator())
                : _buildCalendarGrid(controller, theme),
          ),
          // Legend
          _buildLegend(theme),
        ],
      ),
    );
  }

  Widget _buildMonthHeader(AttendanceController controller, ThemeData theme) {
    final month = controller.selectedMonth;
    final monthName = DateFormatUtils.formatMonthYear(month);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: theme.colorScheme.primaryContainer.withOpacity(0.3),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          IconButton(
            icon: const Icon(Icons.chevron_left),
            onPressed: controller.goToPreviousMonth,
          ),
          Text(
            monthName,
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          IconButton(
            icon: const Icon(Icons.chevron_right),
            onPressed: controller.goToNextMonth,
          ),
        ],
      ),
    );
  }

  Widget _buildWeekdayLabels(ThemeData theme) {
    const weekdays = ['T2', 'T3', 'T4', 'T5', 'T6', 'T7', 'CN'];
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: weekdays.map((day) {
          final isWeekend = day == 'T7' || day == 'CN';
          return Expanded(
            child: Center(
              child: Text(
                day,
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: isWeekend
                      ? theme.colorScheme.error.withOpacity(0.6)
                      : theme.colorScheme.onSurface.withOpacity(0.6),
                ),
              ),
            ),
          );
        }).toList(),
      ),
    );
  }

  Widget _buildCalendarGrid(AttendanceController controller, ThemeData theme) {
    final month = controller.selectedMonth;
    final firstDay = DateTime(month.year, month.month, 1);
    final lastDay = DateTime(month.year, month.month + 1, 0);
    final today = DateTime.now();

    // Monday = 1, Sunday = 7
    int startWeekday = firstDay.weekday; // 1 = Monday
    int leadingEmptyDays = startWeekday - 1; // days before the 1st

    final totalCells = leadingEmptyDays + lastDay.day;
    final rows = (totalCells / 7).ceil();

    return ListView.builder(
      padding: const EdgeInsets.symmetric(horizontal: 8),
      itemCount: rows,
      itemBuilder: (context, rowIndex) {
        return Row(
          children: List.generate(7, (colIndex) {
            final cellIndex = rowIndex * 7 + colIndex;
            final dayNumber = cellIndex - leadingEmptyDays + 1;

            if (dayNumber < 1 || dayNumber > lastDay.day) {
              return const Expanded(child: SizedBox(height: 72));
            }

            final date = DateTime(month.year, month.month, dayNumber);
            final isToday = date.year == today.year &&
                date.month == today.month &&
                date.day == today.day;
            final isWeekend = date.weekday == 6 || date.weekday == 7;
            final isFuture = date.isAfter(today);

            final holiday = controller.getHolidayForDate(date);
            final attendance = controller.getAttendanceForDate(date);

            return Expanded(
              child: GestureDetector(
                onTap: () => _showDayDetail(context, date, attendance, holiday),
                child: _buildDayCell(
                  theme: theme,
                  dayNumber: dayNumber,
                  isToday: isToday,
                  isWeekend: isWeekend,
                  isFuture: isFuture,
                  holiday: holiday,
                  attendance: attendance,
                ),
              ),
            );
          }),
        );
      },
    );
  }

  Widget _buildDayCell({
    required ThemeData theme,
    required int dayNumber,
    required bool isToday,
    required bool isWeekend,
    required bool isFuture,
    Holiday? holiday,
    Attendance? attendance,
  }) {
    Color bgColor = Colors.transparent;
    Color textColor = theme.colorScheme.onSurface;
    Color? borderColor;
    Widget? statusIcon;

    if (isToday) {
      borderColor = theme.colorScheme.primary;
    }

    if (holiday != null) {
      // Holiday - red/pink background
      bgColor = Colors.red.shade50;
      textColor = Colors.red.shade700;
      statusIcon = Icon(Icons.celebration, size: 14, color: Colors.red.shade400);
    } else if (isFuture) {
      // Future date - dimmed
      textColor = theme.colorScheme.onSurface.withOpacity(0.3);
    } else if (isWeekend) {
      // Weekend - slightly dimmed
      textColor = theme.colorScheme.error.withOpacity(0.5);
    } else if (attendance != null) {
      if (attendance.isCheckedOut) {
        // Full attendance (checkin + checkout) - green
        bgColor = Colors.green.shade50;
        statusIcon = Icon(Icons.check_circle, size: 14, color: Colors.green.shade600);
      } else {
        // Missing checkout - orange/warning
        bgColor = Colors.orange.shade50;
        statusIcon = Icon(Icons.warning_amber, size: 14, color: Colors.orange.shade600);
      }
    } else if (!isFuture && !isWeekend) {
      // Working day without attendance - missing checkin - red/warning
      bgColor = Colors.red.shade50.withOpacity(0.5);
      statusIcon = Icon(Icons.close, size: 14, color: Colors.red.shade400);
    }

    return Container(
      height: 72,
      margin: const EdgeInsets.all(2),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(8),
        border: borderColor != null
            ? Border.all(color: borderColor, width: 2)
            : Border.all(color: theme.dividerColor.withOpacity(0.2)),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            '$dayNumber',
            style: TextStyle(
              fontSize: 16,
              fontWeight: isToday ? FontWeight.bold : FontWeight.w500,
              color: textColor,
            ),
          ),
          if (statusIcon != null) ...[
            const SizedBox(height: 2),
            statusIcon,
          ],
          if (holiday != null) ...[
            const SizedBox(height: 1),
            Text(
              'Nghỉ',
              style: TextStyle(fontSize: 8, color: Colors.red.shade600),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildLegend(ThemeData theme) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(top: BorderSide(color: theme.dividerColor)),
      ),
      child: Wrap(
        spacing: 16,
        runSpacing: 8,
        children: [
          _legendItem(Icons.check_circle, Colors.green.shade600, 'Đầy đủ'),
          _legendItem(Icons.warning_amber, Colors.orange.shade600, 'Thiếu checkout'),
          _legendItem(Icons.close, Colors.red.shade400, 'Vắng mặt'),
          _legendItem(Icons.celebration, Colors.red.shade400, 'Ngày nghỉ'),
        ],
      ),
    );
  }

  Widget _legendItem(IconData icon, Color color, String label) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 14, color: color),
        const SizedBox(width: 4),
        Text(label, style: TextStyle(fontSize: 11, color: color)),
      ],
    );
  }

  void _showDayDetail(
    BuildContext context,
    DateTime date,
    Attendance? attendance,
    Holiday? holiday,
  ) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) {
        return Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                DateFormatUtils.formatVietnameseDate(date),
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 16),
              if (holiday != null) ...[
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.red.shade50,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Row(
                    children: [
                      Icon(Icons.celebration, color: Colors.red.shade600),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              holiday.name,
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                color: Colors.red.shade700,
                              ),
                            ),
                            if (holiday.description != null &&
                                holiday.description!.isNotEmpty)
                              Text(
                                holiday.description!,
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.red.shade600,
                                ),
                              ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 12),
              ],
              if (attendance != null) ...[
                _detailRow(
                  icon: Icons.login,
                  iconColor: Colors.green,
                  label: 'Check-in',
                  value: DateFormatUtils.formatTime(attendance.checkIn),
                ),
                const SizedBox(height: 8),
                _detailRow(
                  icon: Icons.logout,
                  iconColor: Colors.red,
                  label: 'Check-out',
                  value: attendance.checkOut != null
                      ? DateFormatUtils.formatTime(attendance.checkOut!)
                      : 'Chưa checkout',
                ),
                const SizedBox(height: 8),
                _detailRow(
                  icon: Icons.timer,
                  iconColor: Colors.blue,
                  label: 'Thời gian',
                  value: DateFormatUtils.formatDuration(
                      attendance.checkIn, attendance.checkOut),
                ),
              ] else if (holiday == null) ...[
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.grey.shade100,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: const Row(
                    children: [
                      Icon(Icons.info_outline, color: Colors.grey),
                      SizedBox(width: 8),
                      Text('Không có dữ liệu chấm công'),
                    ],
                  ),
                ),
              ],
              const SizedBox(height: 16),
            ],
          ),
        );
      },
    );
  }

  Widget _detailRow({
    required IconData icon,
    required Color iconColor,
    required String label,
    required String value,
  }) {
    return Row(
      children: [
        Icon(icon, size: 18, color: iconColor),
        const SizedBox(width: 8),
        Text('$label: ', style: const TextStyle(fontWeight: FontWeight.w500)),
        Text(value),
      ],
    );
  }
}
