import 'package:geolocator/geolocator.dart';

class LocationPermissionUtils {
  static Future<void> ensureLocationPermission() async {
    LocationPermission permission = await Geolocator.checkPermission();

    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
    }

    if (permission == LocationPermission.denied) {
      throw Exception('Quyền truy cập vị trí bị từ chối');
    }

    if (permission == LocationPermission.deniedForever) {
      throw Exception('Quyền truy cập vị trí bị từ chối vĩnh viễn');
    }
  }
}
