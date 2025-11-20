<script setup lang="ts">
import UserAvatar from '@/components/UserAvatar.vue'
import FollowButton from '@/components/FollowButton.vue'
import { Plus } from '@element-plus/icons-vue'
import type { User } from '@/types/user.ts'
import { computed } from 'vue'
import clsx from 'clsx'
import { useUserStore } from '@/stores/user.ts'
import { daysFrom } from '@/utils/date.ts'
const userStore = useUserStore()
const props = defineProps({
  user: {
    type: Object as () => User,
    required: true,
  },
  autoPadding: {
    type: Boolean,
    default: true,
  },
  showFollow: {
    type: Boolean,
    default: true,
  },
})
const wrapClass = computed(() => {
  return clsx({
    'w-70': true,
    'p-3': props.autoPadding,
  })
})
const joinDays = computed(() => {
  return daysFrom(userStore.getActiveUser()?.user.createdAt || new Date())
})
const isSelf = computed(() => {
  return userStore.getActiveUser()?.user.id == props.user.id
})
</script>
<template>
  <div :class="wrapClass">
    <div class="flex items-center justify-between">
      <router-link :to="`/user/${user.account}`"
        ><UserAvatar :size="60" :src="user.avatar"></UserAvatar
      ></router-link>
      <FollowButton v-show="!isSelf" :uid="user.id" size="large" :followed="user.followed">
        <template #before>
          <el-icon class="mr-1">
            <Plus></Plus>
          </el-icon>
        </template>
      </FollowButton>
    </div>
    <router-link :to="`/user/${user.account}`" class="flex flex-col justify-center mt-3">
      <span class="text-base">{{ user.username }}</span>
      <span class="text-gray-500 text-sm">
        <el-link>@{{ user.account }}</el-link>
      </span>
    </router-link>
    <div class="mt-3">
      <span class="text-base">{{ user.aboutMe }}</span>
    </div>
    <div v-if="showFollow" class="flex items-center mt-4">
      <router-link :to="`/user/${user.account}/followers`"
        ><el-link>{{ user.followers }}关注者</el-link></router-link
      >
      <router-link :to="`/user/${user.account}/following`"
        ><el-link class="ml-7">{{ user.following }}正在关注</el-link></router-link
      >
    </div>
    <div class="flex mt-4">
      <span class="text-base text-gray-600 dark:text-gray-300">{{ joinDays }}天前加入</span>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>
