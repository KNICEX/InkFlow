<script setup lang="ts">
import UserSimpleItem from '@/components/list/user/UserSimpleItem.vue'
import type { User } from '@/types/user.ts'

defineProps({
  users: {
    type: Array as () => User[],
    default: () => [],
  },
})
</script>

<template>
  <div class="rounded-border">
    <div class="pb-2 px-4 py-4 text-lg font-bold">为你推荐的创作者</div>
    <UserSimpleItem
      class="cursor-pointer p-4 hover-bg"
      v-for="user in users"
      :user="user"
      :key="user.id"
    ></UserSimpleItem>
    <div class="text-center py-2 hover-bg">
      <el-link type="primary">
        <span class="text-base">展示更多</span>
      </el-link>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>
