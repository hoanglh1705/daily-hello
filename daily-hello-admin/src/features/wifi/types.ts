import type { Branch } from "../branch/types"

export type Wifi = {
  id: number
  name: string
  code: string
  ssid: string
  bssid: string
  branch_id: number
  branch: Branch
  branch_name: string
  created_at: string
}
