import { get, post } from './axios.ts'
import type { User } from '@/types/user.ts'
import type { Ink } from '@/types/ink.ts'
import type { Comment } from '@/types/comment.ts'
interface SearchReq {
  keyword: string
  orderBy?: string
  offset: number
  limit: number
}

export const searchUser = (req: SearchReq) => {
  return get<User[]>('/search/user', {
    params: req,
  })
}

export const searchInk = (req: SearchReq) => {
  return get<Ink[]>('/search/ink', {
    params: req,
  })
}

export const searchComment = (req: SearchReq) => {
  return get<Comment[]>('/search/comment', {
    params: req,
  })
}
