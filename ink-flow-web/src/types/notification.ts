import { type Comment, type CommentPayload, emptyComment } from '@/types/comment.ts'
import type { User } from '@/types/user.ts'
import { emptyInk, type Ink } from '@/types/ink.ts'

export interface Notification<S, C> {
  id: string
  recipientId: string
  user: User | null // null if system notification
  subjectType: SubjectType
  subject: S | null // null if system notification
  notificationType: NotificationType
  content: C
  read: boolean
  createdAt: Date
}

export interface MergedLike<S> {
  users: User[]
  subjectType: string
  subject: S
  total: number
  read: boolean
  updatedAt: Date
}

export const emptyMergedLike = (ink?: boolean) => {
  return {
    users: [],
    subjectType: ink ? SubjectType.Ink : SubjectType.Comment,
    subject: ink ? emptyInk() : emptyComment(),
    total: 0,
    read: false,
    updatedAt: new Date(),
  } as MergedLike<Ink | Comment>
}

export interface ReplyContent {
  commentId: number
  source: CommentPayload
  target: CommentPayload | null
}

export const emptyNotification = (type: NotificationType) => {
  return {
    id: '0',
    recipientId: '0',
    user: null,
    subjectType: 'ink',
    subject: null,
    notificationType: type,
    content: null,
    read: false,
    createdAt: new Date(),
  } as Notification<Ink | User | Comment, ReplyContent | null>
}

export enum NotificationType {
  Like = 'like',
  Reply = 'reply',
  Follow = 'follow',
  System = 'system',
  Mention = 'mention',
  Subscribe = 'subscribe',
}

export interface UnreadCountMap {
  [key: string]: number
}

export enum SubjectType {
  Ink = 'ink',
  Comment = 'comment',
  User = 'user',
  System = 'system',
}
