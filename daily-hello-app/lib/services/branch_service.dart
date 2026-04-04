import 'package:dio/dio.dart';
import '../core/network/api_response.dart';
import '../models/branch.dart';
import '../models/branch_wifi.dart';

class BranchService {
  final Dio dio;

  BranchService(this.dio);

  Future<List<Branch>> getBranches() async {
    final res = await dio.get('/v1/branches');
    final list = unwrapApiData(res.data) as List<dynamic>;
    return list.map((e) => Branch.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<Branch> getBranch(String id) async {
    final res = await dio.get('/v1/branches/$id');
    return Branch.fromJson(unwrapApiData(res.data) as Map<String, dynamic>);
  }

  Future<List<BranchWifi>> getBranchWifiList(String branchId) async {
    final res = await dio.get('/v1/branch-wifi/branch/$branchId');
    final data = unwrapApiData(res.data);
    if (data is List) {
      return data
          .map((e) => BranchWifi.fromJson(e as Map<String, dynamic>))
          .toList();
    }
    return [];
  }
}
