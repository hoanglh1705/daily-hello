class Attendance {
  final String id;
  final String userId;
  final DateTime checkIn;
  final DateTime? checkOut;
  final double? lat;
  final double? lng;
  final String? wifiSsid;
  final String? wifiBssid;
  final String status;

  Attendance({
    required this.id,
    required this.userId,
    required this.checkIn,
    this.checkOut,
    this.lat,
    this.lng,
    this.wifiSsid,
    this.wifiBssid,
    this.status = 'present',
  });

  factory Attendance.fromJson(Map<String, dynamic> json) {
    return Attendance(
      id: json['id']?.toString() ?? '',
      userId: json['user_id']?.toString() ?? '',
      checkIn: DateTime.parse(json['check_in_time']),
      checkOut: json['check_out_time'] != null
          ? DateTime.parse(json['check_out_time'])
          : null,
      lat: (json['check_in_lat'] as num?)?.toDouble(),
      lng: (json['check_in_lng'] as num?)?.toDouble(),
      wifiSsid: json['wifi_ssid'],
      wifiBssid: json['wifi_bssid'],
      status: json['status'] ?? 'present',
    );
  }

  bool get isCheckedOut => checkOut != null;
}
