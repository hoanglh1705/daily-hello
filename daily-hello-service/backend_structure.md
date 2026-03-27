OK, chọn **MVC + Echo** là hợp lý cho bài này 👍 (nhanh, dễ code, vẫn đủ scale nếu tổ chức tốt).

Mình sẽ thiết kế cho bạn một **backend structure đơn giản nhưng vẫn “pro”**, gồm:

* MVC rõ ràng
* Response chuẩn hóa
* Error handling sạch
* Dễ mở rộng sau này

---

# 🧱 1. Kiến trúc MVC với Echo

## 📂 Folder structure

```go
/cmd/server/main.go

/internal
  /handlers   // nhận request (Echo handler)
  /services      // business logic
  /repositories  // DB (Postgres)
  /models        // struct DB + DTO
  /middlewares
  /pkg
    /response
    /errors
    /utils

/config
```

---

## 🧠 Mapping MVC

| Layer      | Vai trò                         |
| ---------- | ------------------------------- |
| Handler    | Parse request + return response |
| Service    | Business logic                  |
| Repository | Query DB                        |
| Model      | Struct dữ liệu                  |

---

# 🔄 2. Flow xử lý

```text
Request → Handler → Service → Repository → DB
                            ↓
                       Response JSON
```

---

# 📥 3. Request Design (Echo)

## Ví dụ: Check-in

```go
type CheckInRequest struct {
    Lat       float64 `json:"lat" validate:"required"`
    Lng       float64 `json:"lng" validate:"required"`
    WifiBSSID string  `json:"wifi_bssid"`
    DeviceID  string  `json:"device_id" validate:"required"`
}
```

---

## Handler (Echo)

```go
func (c *AttendanceHandler) CheckIn(ctx echo.Context) error {
    var req models.CheckInRequest

    if err := ctx.Bind(&req); err != nil {
        return response.Error(ctx, errors.ErrInvalidInput)
    }

    if err := ctx.Validate(req); err != nil {
        return response.Error(ctx, errors.ErrInvalidInput)
    }

    result, err := c.service.CheckIn(ctx.Request().Context(), req)
    if err != nil {
        return response.HandleError(ctx, err)
    }

    return response.Success(ctx, result)
}
```

---

# 📤 4. Response chuẩn hóa (QUAN TRỌNG)

## 📦 Format chung

```json
{
  "success": true,
  "data": {},
  "error": null
}
```

---

## ✅ Success

```json
{
  "success": true,
  "data": {
    "id": 1,
    "check_in_time": "2026-03-27T08:00:00Z"
  }
}
```

---

## ❌ Error

```json
{
  "success": false,
  "error": {
    "code": "INVALID_LOCATION",
    "message": "You are not in allowed area"
  }
}
```

---

## 🧰 Implement response helper

```go
package response

import "github.com/labstack/echo/v4"

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *AppError   `json:"error,omitempty"`
}

type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

---

### Success

```go
func Success(c echo.Context, data interface{}) error {
    return c.JSON(200, Response{
        Success: true,
        Data:    data,
    })
}
```

---

### Error

```go
func Error(c echo.Context, err AppError) error {
    return c.JSON(mapStatus(err.Code), Response{
        Success: false,
        Error:   &err,
    })
}
```

---

# ⚠️ 5. Error Handling Design

## 🧩 5.1 Define error chuẩn

```go
package errors

type AppError struct {
    Code    string
    Message string
}
```

---

## 📚 5.2 Danh sách error

```go
var (
    ErrInvalidInput = AppError{"INVALID_INPUT", "Invalid request data"}
    ErrUnauthorized = AppError{"UNAUTHORIZED", "Unauthorized"}
    ErrNotFound     = AppError{"NOT_FOUND", "Resource not found"}
    ErrInternal     = AppError{"INTERNAL_ERROR", "Internal server error"}

    // Attendance
    ErrInvalidLocation  = AppError{"INVALID_LOCATION", "Not in allowed area"}
    ErrFakeGPS          = AppError{"FAKE_GPS", "Fake GPS detected"}
    ErrAlreadyCheckedIn = AppError{"ALREADY_CHECKED_IN", "Already checked in"}
)
```

---

## 🔁 5.3 Mapping HTTP status

```go
func mapStatus(code string) int {
    switch code {
    case "INVALID_INPUT":
        return 400
    case "UNAUTHORIZED":
        return 401
    case "NOT_FOUND":
        return 404
    default:
        return 500
    }
}
```

---

## ⚙️ 5.4 HandleError (quan trọng)

```go
func HandleError(c echo.Context, err error) error {
    if appErr, ok := err.(errors.AppError); ok {
        return Error(c, AppError{
            Code:    appErr.Code,
            Message: appErr.Message,
        })
    }

    return Error(c, AppError{
        Code:    "INTERNAL_ERROR",
        Message: "Something went wrong",
    })
}
```

---

# 🧠 6. Service Layer (business logic)

```go
func (s *AttendanceService) CheckIn(ctx context.Context, req models.CheckInRequest) (*models.Attendance, error) {

    if !s.locationService.IsValid(req.Lat, req.Lng) {
        return nil, errors.ErrInvalidLocation
    }

    // TODO: check already check-in

    att := &models.Attendance{
        UserID:  1,
        Lat:     req.Lat,
        Lng:     req.Lng,
    }

    if err := s.repo.Create(ctx, att); err != nil {
        return nil, errors.ErrInternal
    }

    return att, nil
}
```

---

# 🗄️ 7. Repository Layer

```go
func (r *AttendanceRepository) Create(ctx context.Context, att *models.Attendance) error {
    return r.db.WithContext(ctx).Create(att).Error
}
```

---

# 📊 8. Pagination chuẩn

## Request

```http
GET /attendance?page=1&limit=20
```

## Response

```json
{
  "success": true,
  "data": {
    "items": [],
    "meta": {
      "page": 1,
      "limit": 20,
      "total": 100
    }
  }
}
```

---

# 🔐 9. Middleware cơ bản

## Auth middleware

```go
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Request().Header.Get("Authorization")

        if token == "" {
            return response.Error(c, errors.ErrUnauthorized)
        }

        // parse JWT...

        return next(c)
    }
}
```

---

# 🎯 10. Best practices (MVC version)

* Controller **không chứa logic**
* Service xử lý business
* Error luôn có code
* Response luôn 1 format
* Không return raw error

---

# ✅ Tổng kết

👉 Với Echo + MVC, bạn nên:

* Tách rõ: controller / service / repo
* Chuẩn hóa response + error
* Dùng struct request rõ ràng
* Handle error tập trung

---