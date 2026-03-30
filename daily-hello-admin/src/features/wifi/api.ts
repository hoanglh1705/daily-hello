import axios from '@/services/axios'
import type { PaginatedListResponse } from '@/shared/types/api'
import type { Wifi } from './types'

type ListResponse = PaginatedListResponse<Wifi>

export const getWifiList = (params: { page: number; limit: number; branch_id?: number }) => {
  return axios.get<unknown, ListResponse>('/v1/branch-wifi', { params })
}

export const createWifi = (data: { ssid: string; bssid: string; branch_id: number }) => {
  return axios.post('/v1/branch-wifi', data)
}

export const deleteWifi = (id: number) => {
  return axios.delete(`/v1/branch-wifi/${id}`)
}
