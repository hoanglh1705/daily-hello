# 📱 1. Tổng quan mobile

Mobile cần các chức năng chính:

* Check-in / Check-out (GPS / WiFi)
* Xem lịch sử chấm công
* Dashboard cá nhân
* Phân quyền (Employee / Manager / Admin – tùy scope mobile)

👉 Vì vậy **Flutter app nên thiên về client + gọi API**, không xử lý business logic phức tạp.

---

# 🧱 2. Kiến trúc đơn giản (khuyên dùng)

Không cần Clean Architecture phức tạp, bạn có thể dùng:

👉 **Feature-based + MVVM nhẹ (hoặc Provider/BLoC đơn giản)**

```
lib/
│
├── main.dart
├── app.dart
│
├── core/                  # dùng chung
│   ├── constants/
│   ├── config/
│   │   ├── app_config.dart
│   │   ├── env.dart
│   │   └── constants.dart
│   ├── network/          # dio client, interceptor
│   ├── storage/          # local storage (token)
│   └── utils/
│
├── models/               # DTO từ backend
│   ├── user.dart
│   ├── attendance.dart
│   └── branch.dart
│
├── services/             # call API
│   ├── auth_service.dart
│   ├── attendance_service.dart
│   └── branch_service.dart
│
├── features/
│   ├── auth/
│   │   ├── login_page.dart
│   │   └── auth_controller.dart
│   │
│   ├── attendance/
│   │   ├── checkin_page.dart
│   │   ├── history_page.dart
│   │   └── attendance_controller.dart
│   │
│   ├── dashboard/
│   │   ├── dashboard_page.dart
│   │   └── dashboard_controller.dart
│   │
│   └── profile/
│       ├── profile_page.dart
│       └── profile_controller.dart
│
└── widgets/              # reusable UI
    ├── app_button.dart
    └── app_card.dart
```

---

# ⚙️ 3. Giải thích từng layer

## 🔹 core/

* Chứa thứ dùng chung toàn app
* Ví dụ:

```dart
class ApiClient {
  final Dio dio;

  ApiClient() : dio = Dio(BaseOptions(
    baseUrl: "http://localhost:8282/api",
  ));
}
```

---

## 🔹 models/

Map với backend response

```dart
class Attendance {
  final String id;
  final DateTime checkIn;
  final DateTime? checkOut;

  Attendance({
    required this.id,
    required this.checkIn,
    this.checkOut,
  });

  factory Attendance.fromJson(Map<String, dynamic> json) {
    return Attendance(
      id: json['id'],
      checkIn: DateTime.parse(json['check_in']),
      checkOut: json['check_out'] != null
          ? DateTime.parse(json['check_out'])
          : null,
    );
  }
}
```

---

## 🔹 services/

👉 Chỉ gọi API, không logic

```dart
class AttendanceService {
  final Dio dio;

  AttendanceService(this.dio);

  Future<void> checkIn(double lat, double lng) async {
    await dio.post('/attendance/check-in', data: {
      "lat": lat,
      "lng": lng,
    });
  }

  Future<List<Attendance>> getHistory() async {
    final res = await dio.get('/attendance');
    return (res.data as List)
        .map((e) => Attendance.fromJson(e))
        .toList();
  }
}
```

---

## 🔹 features/

👉 Mỗi feature = 1 module

### Ví dụ: attendance_controller.dart

```dart
class AttendanceController extends ChangeNotifier {
  final AttendanceService service;

  bool isLoading = false;

  AttendanceController(this.service);

  Future<void> checkIn() async {
    isLoading = true;
    notifyListeners();

    final position = await Geolocator.getCurrentPosition();

    await service.checkIn(position.latitude, position.longitude);

    isLoading = false;
    notifyListeners();
  }
}
```

---

## 🔹 UI (page)

```dart
class CheckInPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final controller = context.watch<AttendanceController>();

    return Scaffold(
      appBar: AppBar(title: Text("Check-in")),
      body: Center(
        child: ElevatedButton(
          onPressed: controller.checkIn,
          child: Text("Check In"),
        ),
      ),
    );
  }
}
```

---

# 📡 4. Các package cần dùng

```yaml
dependencies:
  flutter:
  dio: ^5.x
  provider: ^6.x
  geolocator: ^10.x
  wifi_info_flutter: ^2.x
  flutter_secure_storage: ^9.x
```

---

# 📍 5. Flow chính của app

## 🔐 Login

* Gọi `/auth/login`
* Lưu token vào secure storage

---

## 📍 Check-in

1. Lấy GPS (Geolocator)
2. Lấy WiFi SSID/BSSID
3. Gửi lên backend:

```json
{
  "lat": 10.123,
  "lng": 106.123,
  "wifi_ssid": "ABC",
  "wifi_bssid": "xx:xx"
}
```

---

## 📊 History

* GET `/attendance?from=...&to=...`
* Support pagination (important theo đề)

---

# 🚀 6. Gợi ý cải tiến (để ăn điểm sáng tạo)

Dựa theo đề :

### ⭐ Smart features:

* Auto check-in khi vào geofence
* Nhắc check-in bằng notification
* Offline mode (queue request)
* Face ID (optional AI bonus)

---

# 🧠 7. Tóm lại

👉 Cấu trúc bạn nên dùng:

* Feature-based (rất hợp project thực tế)
* Service = gọi API
* Controller = xử lý state
* UI = đơn giản, reactive

👉 Ưu điểm:

* Dễ code nhanh (5 ngày)
* Dễ scale thêm feature
* Không over-engineering

---

Nếu bạn muốn, mình có thể:

* Generate luôn **boilerplate Flutter project đầy đủ**
* Hoặc viết sẵn **login + check-in flow chạy được ngay**
