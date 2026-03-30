export type Attendance = {
  id: number
  user_name: string
  branch_name: string
  check_in_at: string
  check_out_at: string | null
  status: string
}
