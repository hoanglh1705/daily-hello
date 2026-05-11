class Holiday {
  final int id;
  final String name;
  final DateTime date;
  final String? description;

  Holiday({
    required this.id,
    required this.name,
    required this.date,
    this.description,
  });

  factory Holiday.fromJson(Map<String, dynamic> json) {
    return Holiday(
      id: json['id'] as int,
      name: json['name'] as String,
      date: DateTime.parse(json['date']),
      description: json['description'] as String?,
    );
  }
}
