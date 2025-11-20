import type { User } from '@/types/user.ts'

export interface Comment {
  id: string
  biz: string
  bizId: string
  commentator: User
  isAuthor: boolean
  payload: CommentPayload
  parent: Comment | null
  root: Comment | null
  children: Comment[] | null

  stats: CommentStats
  createdAt: Date
}

interface CommentStats {
  liked: boolean
  likeCnt: number
  viewCnt: number
  replyCnt: number
}

export interface CommentPayload {
  content: string
  images: string[]
}

export const emptyCommentStats = (): CommentStats => {
  return {
    liked: false,
    likeCnt: 0,
    viewCnt: 0,
    replyCnt: 0,
  }
}

export const emptyComment = (): Comment => {
  return {
    id: '0',
    biz: '',
    bizId: '0',
    commentator: {
      id: '0',
      avatar: '',
      banner: '',
      username: '',
      account: '',
      aboutMe: '',
      email: '',
      phone: '',
      links: [],
      followed: false,
      followers: 0,
      following: 0,
      createdAt: new Date(),
      updatedAt: new Date(),
    },
    isAuthor: false,
    payload: {
      content: '',
      images: [],
    },
    parent: null,
    root: null,
    children: [],
    stats: emptyCommentStats(),
    createdAt: new Date(),
  }
}
