<script setup lang="ts">
import UserItem from '@/components/list/user/UserItem.vue'
import { demoUsers } from '@/mock/demo_data.ts'
import type { User } from '@/types/user.ts'
import NoData from '@/components/empty/NoData.vue'
import InkLoading from '@/components/loading/InkLoading.vue'
import BackTop from '@/components/BackTop.vue'
import { useActive } from '@/hook/active.ts'
defineProps({
  users: {
    type: Array as () => User[],
    default: () => demoUsers(),
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
  <div
    class="grid gap-x-20 gap-y-10 md:grid-cols-2"
    v-infinite-scroll="loadMore"
    :infinite-scroll-disabled="!active"
  >
    <UserItem v-for="user in users" :user="user" :key="user.id"></UserItem>
  </div>
  <InkLoading v-show="loading"></InkLoading>
  <BackTop></BackTop>
  <NoData v-show="!loading && users.length == 0"></NoData>
</template>

<style scoped lang="scss"></style>
