Trong app (đặc biệt là mobile như Flutter cho Smart Attendance ), việc **show dialog khi lỗi** không chỉ để báo lỗi mà còn ảnh hưởng lớn tới UX. Không phải lỗi nào cũng nên “popup” lên.

Dưới đây là các **cách show dialog + use case + ví dụ thực tế** 👇

---

# 1. Alert Dialog (Popup cơ bản)

![Image](https://i.sstatic.net/ooFwV.png)

![Image](https://cdn.dribbble.com/userupload/4394480/file/original-af229eec7fe3dbaae5433e24d381a247.png?resize=400x0)

![Image](https://miro.medium.com/v2/resize%3Afit%3A1400/0%2AvlpIQjgJ7mcRX4sT.png)

![Image](https://miro.medium.com/v2/resize%3Afit%3A1400/1%2Ajzm4sqHigu8fhtc0NkvxlQ.png)

### ✅ Use case

* Lỗi quan trọng cần user biết ngay
* Không thể tiếp tục flow
* Cần user xác nhận

### 📌 Ví dụ

* Sai mật khẩu khi login
* Token hết hạn → yêu cầu login lại
* Không có quyền (403)

### 💡 Flutter example

```dart
showDialog(
  context: context,
  builder: (_) => AlertDialog(
    title: Text("Lỗi"),
    content: Text("Sai tài khoản hoặc mật khẩu"),
    actions: [
      TextButton(
        onPressed: () => Navigator.pop(context),
        child: Text("OK"),
      )
    ],
  ),
);
```

---

# 2. Snackbar (Thông báo nhẹ, không chặn)

![Image](https://raw.githubusercontent.com/DNQuyTD/save_personal_images/main/Screenshot_20230617-104101~2.png)

![Image](https://material-design.storage.googleapis.com/publish/material_v_9/0Bzhp5Z4wHba3dEZTUF9idzBHMWc/patterns_errors_userinput19.png)

![Image](https://user-images.githubusercontent.com/3165635/92245839-326a0b80-eec5-11ea-87f5-fbcbc3f808f1.png)

![Image](https://user-images.githubusercontent.com/11846339/59129654-0e024500-898b-11e9-9931-d0a012104f97.png)

### ✅ Use case

* Lỗi nhẹ, không cần chặn user
* Hành động phụ thất bại

### 📌 Ví dụ

* Check-in thất bại do GPS yếu
* Load data lỗi nhưng vẫn retry được

### 💡 Flutter example

```dart
ScaffoldMessenger.of(context).showSnackBar(
  SnackBar(content: Text("Không thể kết nối server")),
);
```

---

# 3. Full-screen Error (Trang lỗi riêng)

![Image](https://cdn.dribbble.com/userupload/23243827/file/still-f6f44296110e558269d678677bc8d13e.gif)

![Image](https://cdn.dribbble.com/users/3821/screenshots/2530692/empty_states.jpg)

![Image](https://i.sstatic.net/y9r1a.png)

![Image](https://cdn.dribbble.com/userupload/23716684/file/still-71da13f2d824d255033a395a4c04d490.gif?resize=400x0)

### ✅ Use case

* Lỗi blocking toàn bộ màn hình
* Không có data để hiển thị

### 📌 Ví dụ

* Không có internet
* API fail toàn bộ dashboard
* Lần đầu load app bị lỗi

### 💡 Flutter idea

```dart
if (state == Error) {
  return ErrorScreen(
    message: "Không có kết nối",
    onRetry: fetchData,
  );
}
```

---

# 4. Bottom Sheet Error (Thân thiện UX hơn dialog)

![Image](https://i.sstatic.net/1VH4H.png)

![Image](https://cdn.dribbble.com/userupload/12993630/file/original-665bf7410458f8881aa756ff15cc9e23.png?resize=400x0)

![Image](https://miro.medium.com/1%2AOGijyonwKq1Idn0foht0yA.gif)

![Image](https://cdn.dribbble.com/userupload/30873883/file/still-640c3daab6bf5f21ce82613278936da1.png?resize=400x0)

### ✅ Use case

* Cảnh báo nhưng không quá nghiêm trọng
* Muốn UX mềm mại hơn dialog

### 📌 Ví dụ

* Check-in sai vị trí (ngoài vùng GPS)
* Thiếu quyền location

---

# 5. Inline Error (Hiển thị ngay trong UI)

![Image](https://user-images.githubusercontent.com/51165119/86409585-2a081100-bc87-11ea-880a-abe0cf070b70.png)

![Image](https://i.sstatic.net/cxuHY.png)

![Image](https://media.nngroup.com/media/editor/2018/11/23/best-buy-error.png)

![Image](https://i.sstatic.net/uyUdd.jpg)

### ✅ Use case

* Form validation
* Lỗi cụ thể từng field

### 📌 Ví dụ

* Email không hợp lệ
* Password quá ngắn

### 💡 Flutter example

```dart
TextField(
  decoration: InputDecoration(
    labelText: "Email",
    errorText: isValid ? null : "Email không hợp lệ",
  ),
);
```

---

# 6. Toast (Nhanh – nhẹ – tự biến mất)

### ✅ Use case

* Thông báo cực nhẹ
* Không quan trọng

### 📌 Ví dụ

* “Đã lưu thành công”
* “Copy thành công”

---

# 7. Dialog + Retry (Best practice cho API)

### ✅ Use case

* Lỗi có thể retry
* UX tốt hơn chỉ báo lỗi

### 📌 Ví dụ

* Check-in thất bại do timeout
* Load dashboard fail

```dart
AlertDialog(
  title: Text("Lỗi kết nối"),
  content: Text("Thử lại?"),
  actions: [
    TextButton(onPressed: retry, child: Text("Retry")),
    TextButton(onPressed: close, child: Text("Huỷ")),
  ],
);
```

---

# 🔥 Best Practice (rất quan trọng)

### ❌ Đừng lạm dụng dialog

* Popup liên tục → user khó chịu
* Chỉ dùng khi thật sự cần attention

### ✅ Rule đơn giản

| Loại lỗi         | Cách hiển thị  |
| ---------------- | -------------- |
| Validation       | Inline         |
| Lỗi nhẹ          | Snackbar       |
| Lỗi có thể retry | Dialog + Retry |
| Lỗi blocking     | Full screen    |
| Cảnh báo nhẹ     | Bottom sheet   |

---

# 💡 Áp dụng cho Smart Attendance

Một số case thực tế từ đề bài:

* Fake GPS → **Dialog (nghiêm trọng)**
* Check-in ngoài vùng → **Bottom sheet**
* Mất mạng → **Full screen**
* Sai password → **Inline + dialog**
* Check-in fail → **Snackbar + retry**
