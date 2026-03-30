````md
# 📘 FRONTEND ARCHITECTURE GUIDE
## Daily Hello - Admin Web

---

## 🎯 Mục tiêu

Tài liệu này định nghĩa cách tổ chức code frontend (React) cho hệ thống Smart Attendance.

Yêu cầu:
- Dễ đọc, dễ maintain
- Phù hợp team nhỏ / thời gian ngắn
- Scale tốt cho 100 chi nhánh / 5000 user
- AI tools có thể hiểu và generate code đúng structure

---

# 🧱 1. Kiến trúc tổng thể

Sử dụng:

> ✅ Feature-based architecture + Shared modules

---

## 📁 Folder Structure

```bash
src/
├── app/                 # App config (router, providers)
│   ├── App.tsx
│   └── routes.tsx
│
├── layouts/             # Layouts (Admin layout, sidebar,...)
│   └── AdminLayout.tsx
│
├── shared/              # Reusable across features
│   ├── components/
│   │   ├── Table.tsx
│   │   ├── Modal.tsx
│   │   ├── Form.tsx
│   │   └── Pagination.tsx
│   │
│   ├── hooks/
│   │   └── useDebounce.ts
│   │
│   └── utils/
│       ├── formatDate.ts
│       └── constants.ts
│
├── features/            # 👈 CORE BUSINESS LOGIC
│   ├── branch/
│   ├── wifi/
│   ├── attendance/
│   └── dashboard/
│
├── services/            # API config (axios instance)
│   └── axios.ts
│
├── styles/
└── main.tsx
````

---

# 🧩 2. Feature Structure

Mỗi feature phải tự chứa toàn bộ logic của nó.

### 📁 Example: `branch`

```bash
features/branch/
├── api.ts               # API calls
├── types.ts             # Types/interfaces
├── BranchPage.tsx       # Main page
├── components/
│   ├── BranchTable.tsx
│   └── BranchForm.tsx
```

---

# 🧠 3. Nguyên tắc quan trọng

---

## 🔑 3.1. Feature Ownership

* Mỗi feature **tự quản lý code của nó**
* Không phụ thuộc trực tiếp vào feature khác

```ts
// ❌ Không được làm
features/branch → import từ features/wifi

// ✅ Đúng
shared → dùng chung
```

---

## 🔑 3.2. Shared Rule

Chỉ đưa vào `shared/` nếu:

> ✔ Được dùng ở >= 2 features

Ví dụ:

| Component   | Vị trí  |
| ----------- | ------- |
| Table       | shared  |
| BranchTable | feature |
| Modal       | shared  |
| BranchForm  | feature |

---

## 🔑 3.3. API nằm trong Feature

```ts
// ✅ đúng
features/branch/api.ts

// ❌ sai
services/branchApi.ts
```

---

## 🔑 3.4. Page = entry point

Mỗi feature chỉ có **1 Page chính**

```ts
BranchPage.tsx → render toàn bộ UI branch
```

---

# 🎨 4. Coding Convention

---

## 📛 4.1. Naming

### File

| Type      | Convention      | Example           |
| --------- | --------------- | ----------------- |
| Page      | PascalCase      | `BranchPage.tsx`  |
| Component | PascalCase      | `BranchTable.tsx` |
| Hook      | camelCase + use | `useDebounce.ts`  |
| API       | lowercase       | `api.ts`          |
| Types     | lowercase       | `types.ts`        |

---

## 🧾 4.2. Component Structure

```tsx
import { useEffect } from "react"

type Props = {
  data: any[]
}

export default function BranchTable({ data }: Props) {
  // hooks
  useEffect(() => {}, [])

  // render
  return (
    <div>
      {/* UI */}
    </div>
  )
}
```

---

## 🔌 4.3. API Convention

```ts
import axios from "@/services/axios"

export const getBranches = (params: any) => {
  return axios.get("/v1/branches", { params })
}

export const createBranch = (data: any) => {
  return axios.post("/v1/branches", data)
}
```

---

## 📦 4.4. API Response Format

Backend phải trả:

```json
{
  "data": [],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

---

## 🧮 4.5. State Management

* Ưu tiên:

  * React hooks (`useState`, `useEffect`)
* Không cần Redux (trừ khi project lớn hơn)

---

## 🔍 4.6. Data Fetching Pattern

```tsx
const [data, setData] = useState([])
const [loading, setLoading] = useState(false)

useEffect(() => {
  fetchData()
}, [])

const fetchData = async () => {
  setLoading(true)
  const res = await getBranches()
  setData(res.data)
  setLoading(false)
}
```

---

# 🧱 5. UI Pattern

---

## 📊 5.1. Table Page Pattern

Mỗi page CRUD nên có:

* Filter bar
* Table
* Pagination
* Create/Edit modal

---

## 📐 5.2. Layout

```tsx
<AdminLayout>
  <Page />
</AdminLayout>
```

---

# 🔐 6. Auth & Role (optional)

* Admin → full access
* Manager → branch scoped

---

# ⚡ 7. Performance & Scale

---

## 📌 Pagination (BẮT BUỘC)

```ts
?page=1&limit=20
```

---

## 📌 Không load toàn bộ data

```ts
// ❌ sai
getAllBranches()

// ✅ đúng
getBranches({ page, limit })
```

---

## 📌 Debounce filter

```ts
search input → debounce 300ms
```

---

# 🚫 8. Anti-pattern (CẤM)

---

## ❌ Import chéo feature

```ts
branch → wifi
```

---

## ❌ Logic API nằm trong component

```ts
// ❌
useEffect(() => {
  axios.get(...)
})
```

---

## ❌ Component quá lớn (>300 dòng)

👉 Tách nhỏ

---

# 🤖 9. Rule cho AI Tools

---

## Khi generate code:

AI phải tuân thủ:

1. Đặt code đúng folder:

   * Feature → `features/{name}`
   * Shared → `shared/`

2. Không tạo file ngoài structure

3. Không duplicate component nếu đã có shared

4. Luôn:

   * Có loading state
   * Có error handling (basic)

5. API phải:

   * Tách riêng file `api.ts`
   * Không hardcode URL

---

# 🚀 10. Workflow

---

1. Tạo feature folder
2. Tạo:

   * Page
   * API
   * Components
3. Connect API
4. Add UI (table + modal)
5. Test

---

# 🏁 Kết luận

Architecture này giúp:

* Code sạch
* Dễ mở rộng
* Dễ cho AI generate đúng
* Phù hợp deadline ngắn (5 ngày)