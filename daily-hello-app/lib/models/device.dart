class Device {
  final int id;
  final int? userId;
  final String deviceId;
  final String deviceName;
  final String platform;
  final String model;
  final String status; // pending, approved, rejected
  final int? approvedBy;
  final DateTime? approvedAt;
  final DateTime createdAt;

  Device({
    required this.id,
    this.userId,
    required this.deviceId,
    required this.deviceName,
    required this.platform,
    required this.model,
    required this.status,
    this.approvedBy,
    this.approvedAt,
    required this.createdAt,
  });

  factory Device.fromJson(Map<String, dynamic> json) {
    return Device(
      id: (json['id'] as num).toInt(),
      userId: (json['user_id'] as num?)?.toInt(),
      deviceId: json['device_id'] ?? '',
      deviceName: json['device_name'] ?? '',
      platform: json['platform'] ?? '',
      model: json['model'] ?? '',
      status: json['status'] ?? 'pending',
      approvedBy: (json['approved_by'] as num?)?.toInt(),
      approvedAt: json['approved_at'] != null
          ? DateTime.parse(json['approved_at'] as String)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  bool get isPending => status == 'pending';
  bool get isApproved => status == 'approved';
  bool get isRejected => status == 'rejected';
}
