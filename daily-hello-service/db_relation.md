OK, mình sẽ thiết kế DB cho Smart Attendance theo hướng:

* ✅ Đủ dùng cho MVP
* ✅ Scale được (5000 users, 100 chi nhánh)
* ✅ Dễ query báo cáo
* ✅ Phù hợp với MVC + Echo

---

# 🧱 1. Tổng quan thiết kế

## 🎯 Nguyên tắc

* Multi-branch (chi nhánh riêng)
* Query nhanh (index hợp lý)
* Tách bảng rõ ràng
* Dễ mở rộng (overtime, shift sau này)

---

# 🗂️ 2. ERD Overview (các bảng chính)

```text
users ───────┐
             ├── attendance
branches ────┘
   │
   ├── branch_wifi
   └── shifts (optional)

users ──── devices (status: pending / approved / rejected)
users ──── refresh_tokens
```

---

# 📊 3. Chi tiết từng bảng

---

## 👤 3.1 users

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL, -- admin, manager, employee
    branch_id BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### 📌 Notes

* `role`: phân quyền
* `branch_id`: user thuộc chi nhánh nào

👉 Index:

```sql
CREATE INDEX idx_users_branch ON users(branch_id);
```

---

## 🏢 3.2 branches

```sql
CREATE TABLE branches (
    id BIGSERIAL PRIMARY KEY,
    branch_code VARCHAR(100) UNIQUE NOT NULL,
    parent_branch_code VARCHAR(100),
    name VARCHAR(100) NOT NULL,
    address TEXT,
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    radius INT, -- mét (geofence)
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 📌 Notes

* `lat/lng/radius`: check GPS

---

## 📶 3.3 branch_wifi

```sql
CREATE TABLE branch_wifi (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    branch_id BIGINT NOT NULL,
    ssid VARCHAR(100),
    bssid VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);
```

👉 Index:

```sql
CREATE INDEX idx_wifi_branch ON branch_wifi(branch_id);
```

---

## ⏱️ 3.4 attendance (CORE TABLE)

```sql
CREATE TABLE attendance (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    branch_id BIGINT NOT NULL,

    check_in_time TIMESTAMP,
    check_out_time TIMESTAMP,

    check_in_lat DOUBLE PRECISION,
    check_in_lng DOUBLE PRECISION,

    check_out_lat DOUBLE PRECISION,
    check_out_lng DOUBLE PRECISION,

    check_in_type VARCHAR(20), -- wifi, gps
    check_out_type VARCHAR(20), -- wifi, gps
    check_in_wifi_bssid VARCHAR(100),
    check_out_wifi_bssid VARCHAR(100),
    check_in_device_id VARCHAR(100),
    check_out_device_id VARCHAR(100),

    check_in_status VARCHAR(20), -- waiting_approve, approved, rejected
    check_out_status VARCHAR(20), -- waiting_approve, approved, rejected

    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## ⚡ Index (RẤT QUAN TRỌNG)

```sql
CREATE INDEX idx_att_user ON attendance(user_id);
CREATE INDEX idx_att_branch ON attendance(branch_id);
CREATE INDEX idx_att_checkin ON attendance(check_in_time);
```

👉 Query nhanh:

* Lịch sử user
* Report theo branch
* Filter theo thời gian

---

# 🧠 4. Optional (ăn điểm nâng cao)

---

## 🕒 4.1 shifts (ca làm)

```sql
CREATE TABLE shifts (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT,
    start_time TIME,
    end_time TIME,
    created_at TIMESTAMP DEFAULT NOW()
);
```

👉 Dùng để:

* Tính trễ / đúng giờ
* Overtime

---

## 📅 4.2 attendance_summary (cache report)

```sql
CREATE TABLE attendance_summary (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    date DATE,
    total_hours FLOAT,
    overtime_hours FLOAT,
    status VARCHAR(20)
);
```

👉 Giúp:

* Dashboard nhanh hơn

---

## 📱 4.3 devices (device registration + approval)

Thiết bị phải được admin/manager phê duyệt trước khi được phép chấm công.

```sql
CREATE TABLE devices (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(100) NOT NULL,
    device_name VARCHAR(200),
    platform VARCHAR(20),     -- android / ios
    model VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, approved, rejected
    approved_by BIGINT REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 📌 Status flow

```
[register] → pending → approved ✅ (chấm công được)
                      → rejected ❌ (không chấm công được)
```

### 👉 Index

```sql
CREATE INDEX idx_devices_user ON devices(user_id);
CREATE INDEX idx_devices_status ON devices(status);
```

### 🔗 Liên quan

* `approved_by` → FK `users.id` (admin hoặc manager thực hiện duyệt)
* Khi user bị xóa → cascade xóa toàn bộ thiết bị của user đó

---

# 🔥 5. Quy tắc business quan trọng

## ✔️ 1 ngày chỉ 1 record attendance

👉 enforce:

```sql
CREATE UNIQUE INDEX unique_user_day
ON attendance(user_id, DATE(check_in_time));
```

---

## ✔️ Check-in flow

* Insert record
* Check-out → update record

---

# 📈 6. Scale Strategy (điểm kiến trúc)

## ⚡ Với 5000 user:

### 1. Partition table

```sql
attendance_2026_03
attendance_2026_04
```

---

### 2. Read-heavy tối ưu

* Index + pagination
* Cache dashboard bằng Redis

---

### 3. Query chuẩn

```sql
SELECT * FROM attendance
WHERE user_id = ?
AND check_in_time BETWEEN ? AND ?
ORDER BY check_in_time DESC
LIMIT 20 OFFSET 0;
```

---

# 🔐 7. Quan hệ (FK - optional)

👉 Nếu muốn strict:

```sql
ALTER TABLE users
ADD CONSTRAINT fk_branch FOREIGN KEY (branch_id) REFERENCES branches(id);
```

---

# 🎯 8. MVP tối thiểu cần có

👉 Nếu thiếu thời gian, chỉ cần:

* users
* branches
* branch_wifi
* attendance

---

# ✅ Tổng kết

Thiết kế này đảm bảo:

* ✔️ Đáp ứng đề bài 
* ✔️ Scale được 5000 users
* ✔️ Query report nhanh
* ✔️ Dễ mở rộng

---

# 👉 Nếu bạn muốn đi tiếp

Mình có thể:

* Vẽ **ERD diagram đẹp để nộp**
* Generate **migration SQL + seed data**
* Viết luôn **GORM model chuẩn**
* Hoặc tối ưu phần **report query (rất hay bị chậm)**

Chỉ cần nói: *"làm tiếp phần X"* 🚀
