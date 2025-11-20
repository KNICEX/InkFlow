<script setup lang="ts">
import NotificationItem from '@/components/notification/NotificationItem.vue'
import { demoUsers } from '@/mock/demo_data.ts'
import InkCover from '@/components/image/InkCover.vue'
import { defaultBanner } from '@/consts/default.ts'

const subjectType = Math.random() > 0.5 ? 'comment' : 'ink'
const users = demoUsers(1)
</script>

<template>
  <NotificationItem :users="users">
    <template #title>
      <span>在{{ subjectType === 'comment' ? '评论' : '作品' }}中提及了你</span>
    </template>
    <div class="line-clamp-3 my-2">
      Lorem ipsum dolor sit amet, <span class="text-blue-500 font-medium">@用户名</span> consectetur adipisicing elit. Ad adipisci autem consequuntur.
    </div>
    <div v-if="subjectType === 'ink'" class="flex mt-2">
      <InkCover class="w-40 h-30 rounded-xl" :src="defaultBanner"></InkCover>
      <div class="ml-4 flex flex-col justify-between">
        <div>作品标题</div>
        <div>
          <el-tag type="info">#标签</el-tag>
        </div>
        <div class="text-sm text-gray-500">2024年11月2日 13:45</div>
      </div>
    </div>
  </NotificationItem>
</template>

<style scoped lang="scss"></style> 