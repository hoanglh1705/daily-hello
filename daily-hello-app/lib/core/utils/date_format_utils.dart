import 'package:intl/intl.dart';

class DateFormatUtils {
  static final DateFormat _dateFormat = DateFormat('dd/MM/yyyy');
  static final DateFormat _timeFormat = DateFormat('HH:mm');
  static final DateFormat _dateTimeFormat = DateFormat('dd/MM/yyyy HH:mm');

  static const _weekdays = [
    'THỨ HAI',
    'THỨ BA',
    'THỨ TƯ',
    'THỨ NĂM',
    'THỨ SÁU',
    'THỨ BẢY',
    'CHỦ NHẬT',
  ];

  static const _months = [
    'THÁNG 1', 'THÁNG 2', 'THÁNG 3', 'THÁNG 4',
    'THÁNG 5', 'THÁNG 6', 'THÁNG 7', 'THÁNG 8',
    'THÁNG 9', 'THÁNG 10', 'THÁNG 11', 'THÁNG 12',
  ];

  static const _weekdaysShort = [
    'Thứ Hai', 'Thứ Ba', 'Thứ Tư', 'Thứ Năm',
    'Thứ Sáu', 'Thứ Bảy', 'Chủ Nhật',
  ];

  static const _monthsLong = [
    'Tháng 1', 'Tháng 2', 'Tháng 3', 'Tháng 4',
    'Tháng 5', 'Tháng 6', 'Tháng 7', 'Tháng 8',
    'Tháng 9', 'Tháng 10', 'Tháng 11', 'Tháng 12',
  ];

  static String formatDate(DateTime date) => _dateFormat.format(date);
  static String formatTime(DateTime date) => _timeFormat.format(date);
  static String formatDateTime(DateTime date) => _dateTimeFormat.format(date);

  /// e.g. "THỨ BA, 24 THÁNG 5, 2024"
  static String formatVietnameseDate(DateTime date) {
    final weekday = _weekdays[date.weekday - 1];
    final month = _months[date.month - 1];
    return '$weekday, ${date.day} $month, ${date.year}';
  }

  /// e.g. "Thứ Hai, 23 Tháng 5"
  static String formatVietnameseDateShort(DateTime date) {
    final weekday = _weekdaysShort[date.weekday - 1];
    final month = _monthsLong[date.month - 1];
    return '$weekday, ${date.day} $month';
  }

  static String formatDuration(DateTime checkIn, DateTime? checkOut) {
    if (checkOut == null) return '--';
    final diff = checkOut.difference(checkIn);
    final hours = diff.inHours;
    final minutes = diff.inMinutes % 60;
    return '${hours}h ${minutes}m';
  }
}
