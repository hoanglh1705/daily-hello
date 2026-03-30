export interface BranchData {
  address: string
  branch_code: string
  created_at: string
  id: number
  lat: number
  lng: number
  name: string
  parent_branch_code: string
  radius: number
  status: string
  updated_at: string
  wifi_list?: any[]
}

export interface User {
  id: number
  branch_id?: number
  branch?: BranchData
  code: string
  created_at: string
  email: string
  name: string
  phone: string
  role: string
  status: string
  updated_at: string
}

export interface UserFilter {
  page: number
  limit: number
  keyword?: string
  branch_id?: number | null
}
