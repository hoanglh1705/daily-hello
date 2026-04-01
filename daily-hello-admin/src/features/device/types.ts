export type DeviceStatus = 'pending' | 'approved' | 'rejected'

type DeviceUser = {
  id: number
  name: string
  email: string
  code: string
  branch_id?: number | null
}

export type Device = {
  id: number
  user_id: number | null
  user?: DeviceUser | null
  device_id: string
  device_name: string
  platform: string
  model: string
  status: DeviceStatus
  approved_by: number | null
  approved_at: string | null
  created_at: string
}
