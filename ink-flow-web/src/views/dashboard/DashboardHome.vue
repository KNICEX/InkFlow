<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import DashboardContent from '@/views/dashboard/DashboardContent.vue'
import InkCard from '@/components/card/InkCard.vue'
import { userActionStats } from '@/service/user.ts'
import type { UserActionStats } from '@/types/stats.ts'
import { daysFrom } from '@/utils/date.ts'
import { useUserStore } from '@/stores/user.ts'

const userStats = ref<UserActionStats>()
const userStore = useUserStore()

onMounted(async () => {
  userStats.value = await userActionStats()
})

const cardItemsInfo = computed(() => {
  return [
    {
      title: '文章',
      icon: 'article',
      count: userStats.value?.inkCount || 0,
    },
    {
      title: '评论',
      icon: 'comment',
      count: userStats.value?.commentCount || 0,
    },
    {
      title: '点赞',
      icon: 'like',
      count: userStats.value?.likeCount || 0,
    },
    {
      title: '关注者',
      icon: 'follower',
      count: userStats.value?.followerCount || 0,
    },
    {
      title: '收藏',
      icon: 'collection',
      count: userStats.value?.favoriteCount || 0,
    },
    {
      title: '正在关注',
      icon: 'following',
      count: userStats.value?.followingCount || 0,
    },
    {
      title: '浏览量',
      icon: 'view',
      count: userStats.value?.viewCount || 0,
    },
  ]
})

const joinDays = computed(() => {
  return daysFrom(userStore.getActiveUser()?.user.createdAt || new Date())
})
</script>

<template>
  <DashboardContent title="Dashboard">
    <div class="flex">
      <div class="flex-1 max-w-240 grid grid-cols-3 gap-6">
        <InkCard
          class="flex items-center"
          bg="gray"
          :border="false"
          v-for="item in cardItemsInfo"
          :key="item.title"
        >
          <div class="h-20 p-2">
            <div>
              <span>{{ item.title }}</span>
            </div>
            <div class="mt-2 text-xl font-semibold">{{ item.count }}</div>
          </div>
        </InkCard>
        <InkCard class="flex items-center" bg="gray" :border="false">
          <div class="h-20 p-2">
            <div>
              <span>已加入</span>
            </div>
            <div class="mt-2 text-xl font-semibold">{{ joinDays }}天</div>
          </div>
        </InkCard>
      </div>
      <div class="w-150 hidden 2xl:block"></div>
    </div>
  </DashboardContent>
</template>

<style scoped lang="scss"></style>
