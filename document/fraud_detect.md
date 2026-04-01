Trong bài toán **Smart Attendance**, chống gian lận là phần *cực kỳ quan trọng* vì user luôn tìm cách “lách luật”. Bạn nên thiết kế theo kiểu **multi-layer (nhiều lớp bảo vệ)** — không có cách nào 100% an toàn nếu chỉ dùng 1 kỹ thuật.

Dưới đây là các **nhóm giải pháp chính + cách implement thực tế** 👇

---

# 1. Xác thực vị trí (Location validation)

## ✅ GPS + Geofencing

* Chỉ cho check-in khi:

  * Nằm trong bán kính (vd: 50–100m)
* Backend kiểm tra:

  * Lat/Lng có hợp lệ không
  * Không cho client tự quyết định

👉 Nâng cao:

* Check tốc độ di chuyển (không thể nhảy 5km trong 1s)
* Lưu lịch sử vị trí gần nhất

---

## ✅ WiFi (SSID + BSSID)

* BSSID (MAC router) **khó fake hơn GPS**
* Mỗi chi nhánh:

  * Lưu whitelist WiFi

👉 Rule:

```text
Check-in hợp lệ nếu:
- GPS hợp lệ OR
- WiFi match BSSID
```

📌 Theo đề bài cũng yêu cầu rõ:

* “Xác định vị trí qua WiFi SSID/BSSID hoặc GPS geofencing” 

---

# 2. Chống Fake GPS

## 🚫 Detect mock location (Android)

* Check:

  * `isFromMockProvider`
  * Developer mode bật không

## 🚫 Detect app giả GPS

* Scan installed apps:

  * Fake GPS, Location Spoofer...

## 🚫 Check sensor inconsistency

* GPS nói đứng yên
* Nhưng:

  * Accelerometer / gyro → đang di chuyển

👉 Nếu mismatch → flag gian lận

---

# 3. Chống VPN / Proxy

## 🚫 Detect VPN

* Check:

  * IP thuộc datacenter (AWS, GCP...)
  * Network interface (tun0, ppp0...)

👉 Rule:

```text
VPN bật → không cho check-in
```

---

# 4. Device Integrity (RẤT QUAN TRỌNG)

## 🚫 Root / Jailbreak detection

* Android:

  * SafetyNet / Play Integrity API
* iOS:

  * Jailbreak detection

👉 Rule:

```text
Thiết bị không an toàn → block
```

---

## 🚫 Device binding

* Mỗi user chỉ dùng:

  * 1–2 device

👉 Lưu:

```json
user_devices {
  user_id,
  device_id,
  last_used
}
```

---

# 5. Backend Anti-cheat (đừng tin client)

## ⚠️ Nguyên tắc:

👉 **Client chỉ gửi data — backend quyết định**

### Validate:

* Timestamp:

  * Không cho sửa giờ local
* Khoảng cách:

  * Tính bằng backend
* Check-in liên tục:

  * Không spam

---

## 🚫 Rule bất thường (Anomaly detection)

Ví dụ:

* Check-in:

* HCM → 1 phút sau Hà Nội
* Check-in ngoài giờ
* Check-in quá nhiều lần

👉 Flag:

```text
status = suspicious
```

---

# 6. Camera / Face ID (level nâng cao)

## 📸 Selfie khi check-in

* So sánh:

  * Face recognition
* Detect:

  * Ảnh chụp lại / deepfake

👉 Có thể dùng:

* AI model
* Hoặc dịch vụ bên thứ 3

---

# 7. QR Code nội bộ (creative)

👉 Mỗi chi nhánh:

* Có QR thay đổi theo thời gian

Check-in:

* Scan QR + GPS/WiFi

👉 Khó fake từ xa

---

# 8. Logging & Audit

## 🧾 Log EVERYTHING

* IP
* Device
* Location
* WiFi

👉 Sau này:

* Audit
* Machine learning detect gian lận

---

# 9. Strategy tổng thể (BEST PRACTICE)

👉 Không dùng 1 cách — dùng combo:

```text
Check-in hợp lệ nếu:
- GPS đúng
- WiFi đúng (nếu có)
- Không dùng VPN
- Không fake GPS
- Device hợp lệ
```

👉 Thêm:

* Risk score:

```text
score = 0 → 100

> 80 → OK
50–80 → warning
< 50 → block
```

---

# 10. Gợi ý kiến trúc

### Mobile

* Check sơ bộ:

  * GPS
  * WiFi
  * Root/VPN

### Backend (QUAN TRỌNG NHẤT)

* Validate lại toàn bộ
* Tính toán:

  * distance
  * rule
  * fraud detection

---

# ✅ Kết luận

Muốn hệ thống “xịn”:
👉 Không phải chống tuyệt đối
👉 Mà là **làm gian lận trở nên khó và tốn công**
