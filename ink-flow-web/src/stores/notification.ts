import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { UnreadCountMap } from '@/types/notification.ts'
import { unreadCount } from '@/service/notification.ts'

export const useUnreadNotificationStore = defineStore('unread-notification-store', () => {
  const unreadMap = ref<UnreadCountMap>({
    reply: 0,
    like: 0,
    follow: 0,
    system: 0,
    mention: 0,
  })

  const autoUpdate = async () => {
    unreadMap.value = await unreadCount()
    // 1分钟更新一次
    setInterval(() => {
      unreadCount().then((res) => {
        unreadMap.value = res
      })
    }, 1000 * 60)
  }
  autoUpdate()

  return {
    unreadMap,
  }
})
