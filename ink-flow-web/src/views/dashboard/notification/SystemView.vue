<script setup lang="ts">
import { onMounted, ref } from 'vue'
import DashboardContent from '@/views/dashboard/DashboardContent.vue'
import SystemList from '@/components/notification/SystemList.vue'
import type { Notification } from '@/types/notification.ts'
import { systemNotification } from '@/service/notification.ts'
import { wrapMaxIdPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { useUnreadNotificationStore } from '@/stores/notification.ts'

const nos = ref<Notification<never, string>[]>([])
const noStore = useUnreadNotificationStore()
const limit = 15

const { loadMore, loading } = wrapMaxIdPagedFunc(async (maxId: string) => {
  const res = await systemNotification({
    maxId,
    limit,
  })
  if (res.length == 0) {
    return '0'
  }
  nos.value = [...nos.value, ...res]
  if (res.length < limit) {
    return '0'
  }
  return nos.value[nos.value.length - 1].id
})

onMounted(() => {
  loadMore().then(() => {
    noStore.unreadMap['system'] = 0
  })
})
</script>

<template>
  <DashboardContent title="Notification" class="overflow-hidden">
    <SystemList
      class="overflow-y-auto"
      :notifications="nos"
      :load-more="loadMore"
      :loading="loading"
    ></SystemList>
  </DashboardContent>
</template>

<style scoped lang="scss"></style>
