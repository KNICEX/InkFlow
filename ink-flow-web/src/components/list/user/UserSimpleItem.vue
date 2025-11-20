<script setup lang="ts">
import UserAvatar from '@/components/UserAvatar.vue'
import FollowButton from '@/components/FollowButton.vue'
import UserCard from '@/components/UserCard.vue'
import InkPopover from '@/components/popover/InkPopover.vue'
import type { User } from '@/types/user.ts'
defineProps({
  user: {
    type: Object as () => User,
    required: true,
  },
})
</script>

<template>
  <div class="flex items-center justify-between">
    <div class="flex">
      <InkPopover place="bottom">
        <template #content>
          <UserCard :user="user"></UserCard>
        </template>
        <template #reference>
          <div class="flex">
            <UserAvatar :src="user.avatar" class="cursor-pointer" :size="48"></UserAvatar>
          </div>
        </template>
      </InkPopover>
      <div class="ml-2 flex flex-col justify-center items-start">
        <div type="info" class="text-sm cursor-pointer">{{ user.username }}</div>
        <router-link :to="`/user/${user.account}`">
          <el-link class="text-sm text-gray-500">@{{ user.account }}</el-link>
        </router-link>
      </div>
    </div>
    <FollowButton :uid="user.id" :followed="user.followed"></FollowButton>
  </div>
</template>

<style scoped lang="scss"></style>
