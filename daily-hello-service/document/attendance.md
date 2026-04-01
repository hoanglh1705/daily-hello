# Tài liệu Đặc tả Flow Chấm Công (Check-in & Check-out)

Tài liệu này mô tả chi tiết quy trình, nguyên tắc nghiệp vụ và hành vi thiết kế trong hệ thống `daily-hello-service` dành cho tính năng điểm danh (**Attendance**) của nhân viên.

---

## 1. Điểm danh vào (Check-in)

### API Endpoint
* **URL**: `/api/v1/attendance/check-in`
* **Method**: `POST`
* **Headers**: 
  - `Authorization: Bearer <token>`
  - `X-Signature: <hmac_signature>` (Bảo vệ thao tác chấm công khống)

**Body Request (`application/json`)**:
```json
{
  "lat": 10.762622,
  "lng": 106.660172,
  "wifi_bssid": "00:14:22:01:23:45",
  "wifi_ssid": "Daily_Hello_Network",
  "device_id": "uuid_of_registered_device",
  "branch_id": 1
}
```

### API Endpoint: Check-in bằng GPS (hỗ trợ hình ảnh)
* **URL**: `/api/v1/attendance/check-in-gps`
* **Method**: `POST`
* **Headers**: `Authorization: Bearer <token>`, `X-Signature: <hmac_signature>`
* **Mô tả**: Dùng cho nhân viên chấm công bằng kết nối mạng ngoài. Máy chủ không thực hiện đối chiếu khoảng cách (bỏ qua check khoảng cách với chi nhánh) mà chỉ lưu vết tọa độ (`lat`/`lng`) hiện tại gửi kèm vào data để phục vụ công tác tra soát. Bắt buộc chụp hình kèm theo. Hình ảnh Base64 tự động được nén dung lượng (< 5MB). Trạng thái điểm danh mặc định là `waiting_approve` (cần admin duyệt).

**Body Request (`application/json`)**:
```json
{
  "lat": 10.762622,
  "lng": 106.660172,
  "device_id": "uuid_of_registered_device",
  "branch_id": 1,
  "image": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQ..." 
}
```

### Luồng xử lý (Flow)
1. **Tiếp nhận Request**: Thiết bị di động của nhân viên gửi yêu cầu check-in tới hệ thống, được xác thực thông qua HMAC Signature và JWT Token.
2. **Kiểm tra thiết bị (Device Validation)**: Hệ thống phải xác nhận thiết bị gửi request đang ở trạng thái đã được Approval (`status = 'approved'`).
3. **Kiểm tra vị trí mạng/trạm (Location/Wifi)**: 
   - Kiểm tra IP/Wifi BSSID hoặc GPS xem nhân viên có nằm trong vòng bán kính cho phép của chi nhánh (Branch) đó không.
4. **Kiểm tra lặp lại (Duplication Check)**: 
   - Kiểm tra trong ngày xem `user_id` hiện tại đã có dữ liệu `check_in_time` chưa để đảo bảm mỗi ngày chỉ tạo 1 record attendance duy nhất (áp dụng Unique Index trên DB theo Date).
5. **Ghi nhận trạng thái đúng giờ (On-time / Late)**:
   - Dựa vào mốc thời gian chuẩn là `08:00:00`.
   - Nếu `check_in_time <= 08:00:00`: Hệ thống ghi nhận trạng thái đi làm đúng giờ (`check_in_status = 'on_time'`).
   - Nếu `check_in_time > 08:00:00`: Hệ thống ghi nhận trạng thái ĐI TRỄ (`check_in_status = 'late'`).
6. **Lưu database**: Lưu bản ghi tạo mới vào bảng `attendances` và trả về kết nối thành công.

---

## 2. Điểm danh ra (Check-out)

### API Endpoint
* **URL**: `/api/v1/attendance/check-out`
* **Method**: `POST`
* **Headers**: Giống với Check-in (`Authorization` và `X-Signature`).

**Body Request (`application/json`)**:
```json
{
  "lat": 10.762622,
  "lng": 106.660172,
  "wifi_bssid": "00:14:22:01:23:45",
  "wifi_ssid": "Daily_Hello_Network",
  "device_id": "uuid_of_registered_device",
  "branch_id": 1
}
```

### API Endpoint: Check-out bằng GPS (hỗ trợ hình ảnh)
* **URL**: `/api/v1/attendance/check-out-gps`
* **Method**: `POST`
* **Headers**: Giống với Check-in (`Authorization` và `X-Signature`).
* **Mô tả**: Tương tự như Check-in GPS, API không kiểm tra xác thực vị trí nhân viên có nằm trong chi nhánh hay không mà chỉ lưu lại tọa độ. Bắt buộc có hình. Hệ thống từ chối ảnh > 5MB và tự resize nén, trạng thái lúc ra về là `waiting_approve`.

**Body Request (`application/json`)**:
```json
{
  "lat": 10.762622,
  "lng": 106.660172,
  "device_id": "uuid_of_registered_device",
  "branch_id": 1,
  "image": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQ..."
}
```

### Luồng xử lý (Flow)
1. **Tiếp nhận Request**: Cùng cơ chế bảo mật (HMAC/JWT) tương tự như luồng Check-in.
2. **Tìm kiếm Record ban sáng (Find Record)**:
   - Hệ thống truy vấn bản ghi attendance tương ứng trong **cùng ngày hôm nay** của `user_id` hiện tại.
   - **Thay đổi quan trọng**: Nếu nhân viên vô tình quên chấm công vào buổi sáng (chưa có dữ liệu `check_in`), hệ thống vẫn cho phép nhân viên Check-out và tự động khởi tạo 1 record điểm danh mới dành riêng cho mục đích khai báo lúc ra về (Lúc này `check_in_time` = Null).
3. **Kiểm tra vị trí (Location Check)**:
   - Đối chiếu vị trí giống như bước Check-in.
4. **Cập nhật dữ liệu giờ ra (Update Check-out Time)**:
   - Tính toán trạng thái thời gian về (dựa vào chuẩn `17:00:00` chiều).
   - Nếu `check_out_time <= 17:00:00`: Đang về sớm hơn giờ quy định.
   - Nếu `check_out_time > 17:00:00`: Có thể được đánh giá là `Depart soon` hoặc Overtime. Hiện tại theo logic Dashboard thiết lập cho hệ thống: `check_out_time > '17:00:00'` được gắn nhãn Action type là `"Depart soon"` (theo thiết lập cụ thể của Admin).
   - *Lưu ý: Bất kỳ cập nhật nào cho Check-out đều tăng thời điểm `updated_at` trong DB.*
5. **Lưu database**: Update field `check_out_time`, vị trí `check_out_lat/lng/wifi_bssid` vào record điểm danh hôm nay của User và trả kết quả thành công.

---

## 3. Các điểm lưu ý đối với API & Database

* **Index Database**: Các dữ liệu Query rất nhiều vào cột `check_in_time`, trường thiết kế phải luôn có INDEX trên cột này để Admin Dashboard không bị timeout khi truy xuất khối lượng record (ví dụ thống kê `GetRecentActivities`).
* **Trường Action Type trên Dashboard**: 
   * Frontend / Admin Dashboard không chứa field này trong DB mà được tính toán Real-time bằng lệnh SQL `CASE WHEN` trực tiếp. Cụ thể: 
      * Có Check-out Event -> Xem xét giờ Check-out để lấy type (vd `> 17:00` -> `Depart soon`, hoặc mặc định là `Check-out`).
      * Chưa Check-out nhưng `check_in_time > 08:00` -> `Late arrival`.
      * Mặt định -> `Check-in`.
* **Giới hạn 1 Event / 1 Ngày**: Business flow hiện tại thiết kế phục vụ tối thiểu (MVP) 1 ngày chỉ có 1 Flow Check-in - Check-out. Nếu cần Scale theo dạng (Ca làm - Shift), thì quan hệ bảng Record Check-in phải nối với `shift_id` thay vì độc quyền dựa trên `DATE(check_in_time)`.
