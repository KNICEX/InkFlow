import { get, post } from './axios.ts'
import type {
  MergedLike,
  Notification,
  ReplyContent,
  UnreadCountMap,
} from '@/types/notification.ts'

interface MaxIdPagedReq {
  maxId?: string
  limit: number
}

interface OffsetPagedReq {
  offset?: number
  limit: number
}

export const unreadCount = async () => {
  return await get<UnreadCountMap>('/notification/count')
}

export const replyNotification = async (req: MaxIdPagedReq) => {
  return await get<Notification<never, ReplyContent>[]>('/notification/reply', {
    params: req,
  })
}

export const mergedLikeNotification = async (req: OffsetPagedReq) => {
  return await get<MergedLike<never>[]>('/notification/like', {
    params: req,
  })
}

export const followNotification = async (req: MaxIdPagedReq) => {
  return await get<Notification<never, never>[]>('/notification/follow', {
    params: req,
  })
}

export const systemNotification = async (req: MaxIdPagedReq) => {
  return await get<Notification<never, never>[]>('/notification/system', {
    params: req,
  })
}
