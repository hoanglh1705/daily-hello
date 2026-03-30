import type { Branch } from '../branch/types'

export type AttendanceUser = {
  id: number
  branch_id: number
  code: string
  email: string
  name: string
  phone: string
  role: string
  status: string
  created_at: string
  updated_at: string
}

export type Attendance = {
  id: number
  user_id: number
  branch_id: number
  user: AttendanceUser
  branch: Branch
  check_in_time: string
  check_out_time: string | null
  check_in_lat: number
  check_in_lng: number
  check_out_lat: number | null
  check_out_lng: number | null
  wifi_bssid: string
  device_id: string
  status: string
  created_at: string
}
