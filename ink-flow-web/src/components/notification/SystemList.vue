<script setup lang="ts">
import type { Notification } from '@/types/notification.ts'
import SystemItem from '@/components/notification/SystemItem.vue'
import NoData from '@/components/empty/NoData.vue'
import InkLoading from '@/components/loading/InkLoading.vue'
import { useActive } from '@/hook/active.ts'
defineProps({
  notifications: {
    type: Array as () => Notification<never, string>[],
    default: () => [],
  },
  loadMore: {
    type: Function,
    default: () => {},
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const active = useActive()
</script>

<template>
  <div v-infinite-scroll="loadMore" :infinite-scroll-disabled="!active">
    <SystemItem
      class="mb-2"
      v-for="(no, idx) in notifications"
      :key="idx"
      :notification="no"
    ></SystemItem>
    <InkLoading v-show="loading"></InkLoading>
    <NoData v-show="!loading && notifications.length == 0"></NoData>
  </div>
</template>

<style scoped lang="scss"></style>
