import { get } from './axios.ts'
import type { UserFollow } from '@/types/user.ts'
interface MaxIdPagedReq {
  maxId?: string
  limit: number
}
export const following = async (uid: string, page: MaxIdPagedReq) => {
  return await get<UserFollow[]>(`/user/${uid}/following`, {
    params: page,
  })
}

export const follower = async (uid: string, page: MaxIdPagedReq) => {
  return await get<UserFollow[]>(`/user/${uid}/follower`, {
    params: page,
  })
}
