class User {
  final String id;
  final String username;
  final String fullName;
  final String email;
  final String role;
  final String? branchId;

  User({
    required this.id,
    required this.username,
    required this.fullName,
    required this.email,
    required this.role,
    this.branchId,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id']?.toString() ?? '',
      username: json['username'] ?? '',
      fullName: json['full_name'] ?? '',
      email: json['email'] ?? '',
      role: json['role'] ?? 'employee',
      branchId: json['branch_id']?.toString(),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'username': username,
        'full_name': fullName,
        'email': email,
        'role': role,
        'branch_id': branchId,
      };
}
