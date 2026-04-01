# Thiết kế API - Admin Dashboard

Dựa trên UI design cung cấp và thông tin database trong `db_relation.md`, đưới đây là thiết kế API phục vụ hiển thị Dashboard. 
Để tối ưu hiệu suất, dữ liệu được chia làm 2 API: 1 API cho biểu đồ/số liệu tổng quan và 1 API cho luồng sự kiện (list) để tránh việc query list làm chậm query statistic.

---

## 1. API: Thống kê tổng quan (Dashboard Overview)

**Endpoint:** `GET /api/v1/admin/dashboard/overview`
**Description:** Lấy các số liệu thống kê (tổng user, % đúng giờ, số lượng trễ), dữ liệu biểu đồ 7 ngày và số liệu thống kê nhanh (Quick Stats).

### Query Parameters

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `branch_id` | `int` | No | ID chi nhánh để filter. Bỏ trống hoặc `0` để lấy "All Branches". |
| `date` | `string` | No | Ngày xem báo cáo (Format `YYYY-MM-DD`). Mặc định hiện tại. |

### Response payload

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "summary": {
      "total_employee": 1284,
      "on_time": {
        "percentage": 94.0,
        "trend": 2.4
      },
      "late_arrival": {
        "count": 12,
        "trend": -3
      }
    },
    "attendance_trends": [
      {
        "day": "Mon",
        "date": "2024-03-25",
        "present_count": 1150
      },
      {
        "day": "Tue",
        "date": "2024-03-26",
        "present_count": 1180
      }
    ],
    "quick_stats": {
      "checked_in_today": 1207,
      "pending_approval": 12,
      "active_branches": 3
    }
  }
}
```

### Chi tiết mapping với Database:

*   **`summary.total_employee`**: `COUNT(id)` từ bảng `users` (áp dụng filter `branch_id` nếu có).
*   **`summary.on_time`**:
    *   `percentage`: `(Số lượt đúng giờ / Tổng số lượt check-in hôm nay) * 100`. Lấy dựa vào cột `check_in_status = 'on_time'` ở bảng `attendance`.
    *   `trend`: Mức chênh lệch % so với ngày hôm trước (ví dụ: hôm nay 94%, hôm qua 91.6% -> +2.4%).
*   **`summary.late_arrival`**:
    *   `count`: `COUNT(id)` từ bảng `attendance` có `DATE(check_in_time) = today` và `check_in_status = 'late'`.
    *   `trend`: Mức chêch lệch độ trễ so với ngày hôm trước.
*   **`attendance_trends`**:
    *   Sử dụng `GROUP BY DATE(check_in_time)` trên bảng `attendance` trong khoảng 7 ngày gần nhất để lấy tổng user đi làm mỗi ngày.
*   **`quick_stats.checked_in_today`**: Tổng số lượng record có `check_in_time` khác null trong bảng `attendance` của ngày hôm nay.
*   **`quick_stats.pending_approval`**: `COUNT(id)` trên bảng `devices` trạng thái chờ duyệt đăng ký thiết bị (`status = 'pending'`). (Phục vụ cho flow device registration trong database).
*   **`quick_stats.active_branches`**: `COUNT(id)` trên bảng `branches`.

---

## 2. API: Lịch sử hoạt động gần nhất (Recent Activity)

**Endpoint:** `GET /api/v1/admin/dashboard/recent-activities`
**Description:** Lấy danh sách các hoạt động check-in/check-out/đi trễ mới nhất của nhân viên trong ngày.

### Query Parameters

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `branch_id` | `int` | No | ID chi nhánh để filter. |
| `limit` | `int` | No | Số lượng bản ghi trả về, mặc định `10`. |

### Response payload

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1023,
        "user_name": "Nguyen Van A",
        "avatar_text": "N",
        "action_type": "Check-in",
        "time": "08:01",
        "timestamp": "2024-04-01T08:01:00Z"
      },
      {
        "id": 1024,
        "user_name": "Tran Thi B",
        "avatar_text": "T",
        "action_type": "Check-in",
        "time": "08:15",
        "timestamp": "2024-04-01T08:15:00Z"
      },
      {
        "id": 1025,
        "user_name": "Le Van C",
        "avatar_text": "L",
        "action_type": "Late arrival",
        "time": "09:30",
        "timestamp": "2024-04-01T09:30:00Z"
      }
    ]
  }
}
```

### Chi tiết mapping với Database:

*   **Logic Query**: Cần `JOIN` bảng `attendance` với `users` (qua user_id) để lấy được `users.name`. 
*   **Order**: Sắp xếp ưu tiên độ mới của thời điểm thao tác (`ORDER BY GREATEST(check_out_time, check_in_time) DESC`).
*   **`avatar_text`**: Ký tự đầu tiên trong tên của User (VD: "Nguyen" lấy "N", "Tran" lấy "T"), có thể trả từ Back-end để FE tiện render luôn icon.
*   **`action_type`**: Đánh giá dựa trên 2 cột `check_in_time`, `check_out_time` và thời gian điểm danh:
    *   Nếu sự kiện gần nhất là `check_out` -> Type là `"Check-out"`.
    *   Nếu đánh dấu checkin với `check_in_time > 8:00` -> Type là `"Late arrival"`.
    *   Nếu đánh dấu checkout với `check_out_time > '17:00'` -> Type là `"Depart soon"`.
*   **`time`**: Thời gian thao tác trực quan theo format `HH:mm`.
