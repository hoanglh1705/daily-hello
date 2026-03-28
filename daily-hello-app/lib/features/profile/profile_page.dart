import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'profile_controller.dart';
import '../../widgets/app_button.dart';
import '../../widgets/app_card.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ProfileController>().loadProfile();
    });
  }

  @override
  Widget build(BuildContext context) {
    final controller = context.watch<ProfileController>();
    final user = controller.user;
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Hồ sơ cá nhân'),
        centerTitle: true,
      ),
      body: controller.isLoading
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              padding: const EdgeInsets.all(20),
              children: [
                // Avatar
                Center(
                  child: CircleAvatar(
                    radius: 48,
                    backgroundColor: theme.colorScheme.primary.withAlpha(30),
                    child: Icon(Icons.person,
                        size: 52, color: theme.colorScheme.primary),
                  ),
                ),
                const SizedBox(height: 12),
                Center(
                  child: Text(
                    user?.fullName ?? '--',
                    style: theme.textTheme.headlineSmall
                        ?.copyWith(fontWeight: FontWeight.bold),
                  ),
                ),
                const SizedBox(height: 4),
                Center(
                  child: Container(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.primary.withAlpha(20),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Text(
                      _roleLabel(user?.role),
                      style: TextStyle(
                          color: theme.colorScheme.primary,
                          fontWeight: FontWeight.w600),
                    ),
                  ),
                ),
                const SizedBox(height: 24),

                // Info
                AppCard(
                  child: Column(
                    children: [
                      _InfoRow(
                          icon: Icons.badge_outlined,
                          label: 'Username',
                          value: user?.username ?? '--'),
                      const Divider(height: 24),
                      _InfoRow(
                          icon: Icons.email_outlined,
                          label: 'Email',
                          value: user?.email ?? '--'),
                      const Divider(height: 24),
                      _InfoRow(
                          icon: Icons.business_outlined,
                          label: 'Chi nhánh',
                          value: user?.branchId ?? 'Chưa phân công'),
                    ],
                  ),
                ),
                const SizedBox(height: 24),

                AppButton(
                  label: 'Đăng xuất',
                  icon: Icons.logout,
                  color: Colors.red,
                  onPressed: () async {
                    final confirmed = await showDialog<bool>(
                      context: context,
                      builder: (ctx) => AlertDialog(
                        title: const Text('Đăng xuất'),
                        content:
                            const Text('Bạn có chắc muốn đăng xuất không?'),
                        actions: [
                          TextButton(
                              onPressed: () => Navigator.pop(ctx, false),
                              child: const Text('Hủy')),
                          TextButton(
                              onPressed: () => Navigator.pop(ctx, true),
                              child: const Text('Đăng xuất',
                                  style: TextStyle(color: Colors.red))),
                        ],
                      ),
                    );
                    if (confirmed == true && context.mounted) {
                      await context
                          .read<ProfileController>()
                          .logout(context);
                    }
                  },
                ),
              ],
            ),
    );
  }

  String _roleLabel(String? role) {
    switch (role) {
      case 'admin':
        return 'Admin';
      case 'manager':
        return 'Quản lý';
      default:
        return 'Nhân viên';
    }
  }
}

class _InfoRow extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;

  const _InfoRow({
    required this.icon,
    required this.label,
    required this.value,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, size: 20, color: Colors.grey),
        const SizedBox(width: 12),
        Text('$label:', style: const TextStyle(color: Colors.grey)),
        const SizedBox(width: 8),
        Expanded(
          child: Text(value,
              style: const TextStyle(fontWeight: FontWeight.w500),
              overflow: TextOverflow.ellipsis),
        ),
      ],
    );
  }
}
