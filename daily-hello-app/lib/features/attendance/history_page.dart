import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'attendance_controller.dart';
import '../../widgets/app_card.dart';
import '../../core/utils/date_format_utils.dart';

class HistoryPage extends StatefulWidget {
  const HistoryPage({super.key});

  @override
  State<HistoryPage> createState() => _HistoryPageState();
}

class _HistoryPageState extends State<HistoryPage> {
  final _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<AttendanceController>().loadHistory(refresh: true);
    });
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent - 200) {
      context.read<AttendanceController>().loadHistory();
    }
  }

  @override
  Widget build(BuildContext context) {
    final controller = context.watch<AttendanceController>();

    return Scaffold(
      appBar: AppBar(
        title: const Text('Lịch sử chấm công'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () =>
                context.read<AttendanceController>().loadHistory(refresh: true),
          ),
        ],
      ),
      body: controller.isLoadingHistory && controller.history.isEmpty
          ? const Center(child: CircularProgressIndicator())
          : controller.history.isEmpty
              ? const Center(
                  child: Text('Chưa có lịch sử chấm công'),
                )
              : RefreshIndicator(
                  onRefresh: () =>
                      context.read<AttendanceController>().loadHistory(refresh: true),
                  child: ListView.builder(
                    controller: _scrollController,
                    padding: const EdgeInsets.all(16),
                    itemCount: controller.history.length +
                        (controller.hasMore ? 1 : 0),
                    itemBuilder: (ctx, i) {
                      if (i == controller.history.length) {
                        return const Padding(
                          padding: EdgeInsets.all(16),
                          child: Center(child: CircularProgressIndicator()),
                        );
                      }
                      final item = controller.history[i];
                      return Padding(
                        padding: const EdgeInsets.only(bottom: 12),
                        child: AppCard(
                          child: Row(
                            children: [
                              Container(
                                width: 4,
                                height: 60,
                                decoration: BoxDecoration(
                                  color: item.isCheckedOut
                                      ? Colors.blue
                                      : Colors.green,
                                  borderRadius: BorderRadius.circular(2),
                                ),
                              ),
                              const SizedBox(width: 16),
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      DateFormatUtils.formatDate(item.checkIn),
                                      style: const TextStyle(
                                          fontWeight: FontWeight.bold,
                                          fontSize: 15),
                                    ),
                                    const SizedBox(height: 4),
                                    Row(
                                      children: [
                                        const Icon(Icons.login,
                                            size: 14, color: Colors.green),
                                        const SizedBox(width: 4),
                                        Text(DateFormatUtils.formatTime(
                                            item.checkIn)),
                                        const SizedBox(width: 16),
                                        const Icon(Icons.logout,
                                            size: 14, color: Colors.red),
                                        const SizedBox(width: 4),
                                        Text(item.checkOut != null
                                            ? DateFormatUtils.formatTime(
                                                item.checkOut!)
                                            : '--:--'),
                                      ],
                                    ),
                                  ],
                                ),
                              ),
                              Column(
                                crossAxisAlignment: CrossAxisAlignment.end,
                                children: [
                                  Text(
                                    DateFormatUtils.formatDuration(
                                        item.checkIn, item.checkOut),
                                    style: const TextStyle(
                                        fontWeight: FontWeight.bold,
                                        fontSize: 16),
                                  ),
                                  const SizedBox(height: 4),
                                  Container(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 8, vertical: 2),
                                    decoration: BoxDecoration(
                                      color: item.isCheckedOut
                                          ? Colors.blue[50]
                                          : Colors.green[50],
                                      borderRadius: BorderRadius.circular(8),
                                    ),
                                    child: Text(
                                      item.isCheckedOut ? 'Đầy đủ' : 'Đang làm',
                                      style: TextStyle(
                                        fontSize: 11,
                                        color: item.isCheckedOut
                                            ? Colors.blue
                                            : Colors.green,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  ),
                ),
    );
  }
}
