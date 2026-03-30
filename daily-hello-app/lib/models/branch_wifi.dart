class BranchWifi {
  final int id;
  final String code;
  final String name;
  final int branchId;
  final String? ssid;
  final String? bssid;

  BranchWifi({
    required this.id,
    required this.code,
    required this.name,
    required this.branchId,
    this.ssid,
    this.bssid,
  });

  factory BranchWifi.fromJson(Map<String, dynamic> json) {
    return BranchWifi(
      id: json['id'] as int,
      code: json['code'] ?? '',
      name: json['name'] ?? '',
      branchId: json['branch_id'] as int,
      ssid: json['ssid'] as String?,
      bssid: json['bssid'] as String?,
    );
  }
}
