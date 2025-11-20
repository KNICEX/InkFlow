<script setup lang="ts">
import UserAvatar from '@/components/UserAvatar.vue'
import type { User } from '@/types/user.ts'
import FollowButton from '@/components/FollowButton.vue'
defineProps({
  user: {
    type: Object as () => User,
    required: true,
  },
})
</script>

<template>
  <div class="flex justify-between">
    <div class="flex">
      <router-link :to="`/user/${user.account}`">
        <UserAvatar :src="user.avatar" :size="80"></UserAvatar>
      </router-link>
      <router-link :to="`/user/${user.account}`" class="ml-6 cursor-pointer">
        <div class="text-lg">{{ user.username }}</div>
        <div class="nav-text">
          <el-link>@{{ user.account }}</el-link>
        </div>
        <div class="line-clamp-3">{{ user.aboutMe }}</div>
      </router-link>
    </div>
    <div>
      <FollowButton size="large" :uid="user.id" :followed="user.followed"></FollowButton>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>
