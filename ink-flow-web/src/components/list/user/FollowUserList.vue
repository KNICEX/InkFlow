<script setup lang="ts">
import FollowUser from '@/components/list/user/FollowUser.vue'
import type { UserFollow } from '@/types/user.ts'
import NoData from '@/components/empty/NoData.vue'
import InkLoading from '@/components/loading/InkLoading.vue'
import { useActive } from '@/hook/active.ts'
defineProps({
  users: {
    type: Array as () => UserFollow[],
    required: true,
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
  <div class="w-full" v-infinite-scroll="loadMore" :infinite-scroll-disabled="!active">
    <FollowUser v-for="user in users" :key="user.id" :user="user"></FollowUser>
    <InkLoading v-show="loading"></InkLoading>
    <NoData v-show="!loading && users.length == 0" description="还没有关注噢"></NoData>
  </div>
</template>

<style scoped lang="scss"></style>
