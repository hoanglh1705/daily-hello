import 'package:dio/dio.dart';
import '../core/network/api_response.dart';
import '../models/branch.dart';

class BranchService {
  final Dio dio;

  BranchService(this.dio);

  Future<List<Branch>> getBranches() async {
    final res = await dio.get('/branches');
    final list = unwrapApiData(res.data) as List<dynamic>;
    return list.map((e) => Branch.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<Branch> getBranch(String id) async {
    final res = await dio.get('/branches/$id');
    return Branch.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
  }
}
