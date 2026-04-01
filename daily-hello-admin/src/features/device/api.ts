import axios from '@/services/axios'
import type { PaginatedListResponse } from '@/shared/types/api'
import type { Device, DeviceStatus } from './types'

type ListResponse = PaginatedListResponse<Device>

export const getDevices = (params: {
  page: number
  limit: number
  status: DeviceStatus
  branch_id?: number
}) => {
  return axios.get<unknown, ListResponse>('/v1/admin/devices', { params })
}

export const approveDevice = (id: number) => {
  return axios.put(`/v1/admin/devices/${id}/approve`)
}

export const rejectDevice = (id: number) => {
  return axios.put(`/v1/admin/devices/${id}/reject`)
}
