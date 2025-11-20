<script setup lang="ts">
import { useUnreadNotificationStore } from '@/stores/notification.ts'
import { NotificationType } from '@/types/notification.ts'
import { computed } from 'vue'
import { useUserStore } from '@/stores/user.ts'

const unreadStore = useUnreadNotificationStore()

const overdueMap = computed(() => {
  return {
    [NotificationType.Reply]:
      unreadStore.unreadMap[NotificationType.Reply] > 100
        ? '99+'
        : unreadStore.unreadMap[NotificationType.Reply].toString(),
    [NotificationType.Follow]:
      unreadStore.unreadMap[NotificationType.Follow] > 100
        ? '99+'
        : unreadStore.unreadMap[NotificationType.Follow].toString(),
    [NotificationType.Like]:
      unreadStore.unreadMap[NotificationType.Like] > 100
        ? '99+'
        : unreadStore.unreadMap[NotificationType.Like].toString(),
    [NotificationType.System]:
      unreadStore.unreadMap[NotificationType.System] > 100
        ? '99+'
        : unreadStore.unreadMap[NotificationType.System].toString(),
  }
})
const userStore = useUserStore()
const toNotificationUrl = (type: string) => {
  return `/dashboard/${userStore.getActiveUser()?.user.account}/${type}`
}

const numClass =
  'absolute rounded-full bg-[var(--primary-color)] text-sm font-semibold text-white right-1 -top-1 flex h-5 px-2 justify-center items-center'
</script>

<template>
  <div class="my-4 text-base">
    <router-link class="popover-button relative" :to="toNotificationUrl('reply')">
      <span>回复我的</span>
      <div v-if="unreadStore.unreadMap[NotificationType.Reply] > 0" :class="numClass">
        {{ overdueMap[NotificationType.Reply] }}
      </div>
    </router-link>
    <router-link class="popover-button relative" :to="toNotificationUrl('follow')">
      <span>关注我的</span>
      <span v-if="unreadStore.unreadMap[NotificationType.Follow] > 0" :class="numClass">{{
        overdueMap[NotificationType.Follow]
      }}</span>
    </router-link>
    <router-link class="popover-button relative" :to="toNotificationUrl('like')">
      <span>收到的赞</span>
      <span v-if="unreadStore.unreadMap[NotificationType.Like] > 0" :class="numClass">{{
        overdueMap[NotificationType.Like]
      }}</span>
    </router-link>
    <router-link class="popover-button relative" :to="toNotificationUrl('system')">
      <span>系统消息</span>
      <span v-if="unreadStore.unreadMap[NotificationType.System] > 0" :class="numClass">{{
        overdueMap[NotificationType.System]
      }}</span>
    </router-link>
  </div>
</template>

<style scoped lang="scss"></style>
