import axios from '@/services/axios'
import type { DashboardOverviewResponse, DashboardRecentActivitiesResponse } from './types'

export const getDashboardOverview = (params?: { branch_id?: number }) => {
  return axios.get<unknown, { data: DashboardOverviewResponse }>('/v1/admin/dashboard/overview', { params })
}

export const getRecentActivities = (params?: { branch_id?: number; limit?: number }) => {
  return axios.get<unknown, { data: DashboardRecentActivitiesResponse }>('/v1/admin/dashboard/recent-activities', { params })
}
