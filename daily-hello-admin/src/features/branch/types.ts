export type BranchWifi = {
  id: number
  code: string
  name: string
  branch_id: number
  ssid: string
  bssid: string
  created_at: string
}

export type Branch = {
  id: number
  branch_code: string
  parent_branch_code: string
  name: string
  address: string
  lat: number
  lng: number
  radius: number
  wifi_list: BranchWifi[]
  status: string
  created_at: string
  updated_at: string
}
