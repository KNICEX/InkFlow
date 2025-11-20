<script setup lang="ts">
import ReplyItem from '@/components/notification/ReplyItem.vue'
import type { Notification, ReplyContent } from '@/types/notification'
import NoData from '@/components/empty/NoData.vue'
import InkLoading from '@/components/loading/InkLoading.vue'
import { useActive } from '@/hook/active.ts'
defineProps({
  replies: {
    type: Array as () => Notification<never, ReplyContent>[],
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
    <ReplyItem v-for="reply in replies" :key="reply.id" :reply="reply" class="mb-2"></ReplyItem>
    <InkLoading v-show="loading"></InkLoading>
    <NoData v-show="!loading && replies.length == 0"></NoData>
  </div>
</template>

<style scoped lang="scss"></style>
