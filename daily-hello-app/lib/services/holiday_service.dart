import 'package:dio/dio.dart';
import '../core/network/api_response.dart';
import '../models/holiday.dart';

class HolidayService {
  final Dio dio;

  HolidayService(this.dio);

  Future<List<Holiday>> getHolidaysByMonth({
    required int year,
    required int month,
  }) async {
    final res = await dio.get('/v1/holidays', queryParameters: {
      'year': year,
      'month': month,
    });
    final data = unwrapApiData(res.data);
    if (data is List) {
      return data
          .map((e) => Holiday.fromJson(e as Map<String, dynamic>))
          .toList();
    }
    return [];
  }
}
