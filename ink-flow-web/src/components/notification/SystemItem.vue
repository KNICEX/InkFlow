<script setup lang="ts">
import NotificationItem from '@/components/notification/NotificationItem.vue'
import type { Notification } from '@/types/notification.ts'

defineProps({
  notification: {
    type: Object as () => Notification<never, string>,
    required: true,
  },
})
</script>

<template>
  <NotificationItem :system="true" :users="[]" :created-at="notification.createdAt">
    <template #title>
      <span>系统通知</span>
    </template>
    <div class="text-gray-600 mt-2 line-clamp-3">
      {{ notification.content }}
    </div>
  </NotificationItem>
</template>

<style scoped lang="scss"></style>
