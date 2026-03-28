class Branch {
  final String id;
  final String name;
  final String branchCode;
  final String? parentBranchCode;
  final String? address;
  final double? lat;
  final double? lng;
  final List<String>? wifiSsids;

  Branch({
    required this.id,
    required this.name,
    required this.branchCode,
    this.parentBranchCode,
    this.address,
    this.lat,
    this.lng,
    this.wifiSsids,
  });

  factory Branch.fromJson(Map<String, dynamic> json) {
    return Branch(
      id: json['id']?.toString() ?? '',
      name: json['name'] ?? '',
      branchCode: json['branch_code'] ?? '',
      parentBranchCode: json['parent_branch_code'],
      address: json['address'],
      lat: (json['lat'] as num?)?.toDouble(),
      lng: (json['lng'] as num?)?.toDouble(),
      wifiSsids: json['wifi_ssids'] != null
          ? List<String>.from(json['wifi_ssids'])
          : null,
    );
  }
}
