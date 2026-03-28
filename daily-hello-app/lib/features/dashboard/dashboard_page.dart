import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'dashboard_controller.dart';
import '../auth/auth_controller.dart';
import '../../widgets/app_card.dart';
import '../../core/utils/date_format_utils.dart';

class DashboardPage extends StatefulWidget {
  const DashboardPage({super.key});

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<DashboardController>().loadDashboard();
      context.read<AuthController>().loadProfile();
    });
  }

  @override
  Widget build(BuildContext context) {
    final controller = context.watch<DashboardController>();
    final authController = context.watch<AuthController>();
    final theme = Theme.of(context);
    final today = controller.todayAttendance;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Daily Hello'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => context.read<DashboardController>().loadDashboard(),
          ),
        ],
      ),
      body: controller.isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: () =>
                  context.read<DashboardController>().loadDashboard(),
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  // Greeting
                  AppCard(
                    color: theme.colorScheme.primary,
                    child: Row(
                      children: [
                        const CircleAvatar(
                          radius: 28,
                          backgroundColor: Colors.white24,
                          child: Icon(Icons.person, size: 32, color: Colors.white),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Xin chào!',
                                style: theme.textTheme.bodyMedium
                                    ?.copyWith(color: Colors.white70),
                              ),
                              Text(
                                authController.currentUser?.fullName ?? '--',
                                style: theme.textTheme.titleMedium?.copyWith(
                                    color: Colors.white,
                                    fontWeight: FontWeight.bold),
                              ),
                              Text(
                                authController.currentUser?.role.toUpperCase() ??
                                    '',
                                style: const TextStyle(
                                    color: Colors.white60, fontSize: 12),
                              ),
                            ],
                          ),
                        ),
                        Text(
                          DateFormatUtils.formatDate(DateTime.now()),
                          style:
                              const TextStyle(color: Colors.white70, fontSize: 12),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Today status
                  Text('Hôm nay',
                      style: theme.textTheme.titleMedium
                          ?.copyWith(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  AppCard(
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: [
                        _StatItem(
                          label: 'Check-in',
                          value: today != null
                              ? DateFormatUtils.formatTime(today.checkIn)
                              : '--:--',
                          icon: Icons.login,
                          color: Colors.green,
                        ),
                        _StatItem(
                          label: 'Check-out',
                          value: today?.checkOut != null
                              ? DateFormatUtils.formatTime(today!.checkOut!)
                              : '--:--',
                          icon: Icons.logout,
                          color: Colors.red,
                        ),
                        _StatItem(
                          label: 'Tổng giờ',
                          value: DateFormatUtils.formatDuration(
                            today?.checkIn ?? DateTime.now(),
                            today?.checkOut,
                          ),
                          icon: Icons.access_time,
                          color: theme.colorScheme.primary,
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Recent history
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text('Gần đây',
                          style: theme.textTheme.titleMedium
                              ?.copyWith(fontWeight: FontWeight.bold)),
                      TextButton(
                        onPressed: () =>
                            Navigator.pushNamed(context, '/history'),
                        child: const Text('Xem tất cả'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  if (controller.recentHistory.isEmpty)
                    const Center(
                        child: Padding(
                      padding: EdgeInsets.all(16),
                      child: Text('Chưa có dữ liệu'),
                    ))
                  else
                    ...controller.recentHistory.map((item) => Padding(
                          padding: const EdgeInsets.only(bottom: 8),
                          child: AppCard(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 16, vertical: 12),
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  DateFormatUtils.formatDate(item.checkIn),
                                  style: const TextStyle(
                                      fontWeight: FontWeight.w500),
                                ),
                                Text(
                                  '${DateFormatUtils.formatTime(item.checkIn)} - '
                                  '${item.checkOut != null ? DateFormatUtils.formatTime(item.checkOut!) : "--:--"}',
                                  style: const TextStyle(color: Colors.grey),
                                ),
                                Text(
                                  DateFormatUtils.formatDuration(
                                      item.checkIn, item.checkOut),
                                  style: TextStyle(
                                      color: theme.colorScheme.primary,
                                      fontWeight: FontWeight.bold),
                                ),
                              ],
                            ),
                          ),
                        )),
                ],
              ),
            ),
    );
  }
}

class _StatItem extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;
  final Color color;

  const _StatItem({
    required this.label,
    required this.value,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Icon(icon, color: color, size: 24),
        const SizedBox(height: 4),
        Text(value,
            style:
                TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: color)),
        Text(label,
            style: const TextStyle(fontSize: 11, color: Colors.grey)),
      ],
    );
  }
}
