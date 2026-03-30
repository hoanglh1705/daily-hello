import axios from '@/services/axios'
import type { DashboardData } from './types'

type DashboardResponse = {
  success: boolean
  data: DashboardData
}

export const getDashboard = (params?: { branch_id?: number }) => {
  return axios.get<unknown, DashboardResponse>('/v1/admin/dashboard', { params })
}
