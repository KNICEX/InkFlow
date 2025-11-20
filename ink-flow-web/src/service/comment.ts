import { del, get, post } from './axios.ts'
import type { CommentPayload, Comment } from '@/types/comment.ts'
export interface CommentReq {
  biz: string
  bizId: string
  maxId?: string
  limit: number
}

export interface MaxIdPagedReq {
  maxId?: string
  limit: number
}

export interface ReplyReq {
  biz: string
  bizId: string
  rootId?: string
  parentId?: string
  payload: CommentPayload
}

export const loadComment = async (req: CommentReq) => {
  return await get<Comment[]>('/comment', {
    params: req,
  })
}

export const loadChild = async (rid: string, req: MaxIdPagedReq) => {
  return await get<Comment[]>('/comment/child', {
    params: {
      rid,
      ...req,
    },
  })
}

export const reply = async (req: ReplyReq) => {
  return await post<string, ReplyReq>('/comment/reply', req)
}

export const likeComment = async (cid: string) => {
  return await post<Comment[], null>(`/comment/like/${cid}`, null)
}

export const cancelLikeComment = async (cid: string) => {
  return await del<never>(`/comment/like/${cid}`)
}

export const deleteComment = async (cid: string) => {
  return await del<Comment[]>(`/comment/${cid}`)
}
