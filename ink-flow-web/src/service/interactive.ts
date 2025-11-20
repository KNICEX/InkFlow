import { post, get, del } from '@/service/axios.ts'
import type { Favorite } from '@/types/interactive.ts'

export interface CreateFavoriteReq {
  biz: string
  name: string
  private: boolean
}
export const createFavorite = async (req: CreateFavoriteReq) => {
  return await post<number, CreateFavoriteReq>('/interactive/favorite', req)
}

export const getFavorites = async (biz: string) => {
  return await get<Favorite[]>(`/interactive/favorite/${biz}`)
}

export const deleteFavorite = async (id: string) => {
  return await del<never>(`/interactive/favorite/${id}`)
}
