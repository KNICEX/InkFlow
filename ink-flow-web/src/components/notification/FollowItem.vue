<script setup lang="ts">
import NotificationItem from '@/components/notification/NotificationItem.vue'
import type { Notification } from '@/types/notification.ts'

defineProps({
  notification: {
    type: Object as () => Notification<never, never>,
    required: true,
  },
})

// const users = demoUsers(Math.ceil(Math.random() * 5))
</script>

<template>
  <NotificationItem :users="[notification.user!]" :created-at="notification.createdAt">
    <template #title>
      <!--      <span v-if="users.length == 1">关注了你</span>-->
      <!--      <span v-else>等{{ users.length }}人关注了你</span>-->
      <span>关注了你</span>
    </template>
  </NotificationItem>
</template>

<style scoped lang="scss"></style>
