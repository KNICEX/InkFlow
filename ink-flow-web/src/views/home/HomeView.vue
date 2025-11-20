<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watchEffect } from 'vue'
import HomeAside from '@/views/home/HomeAside.vue'
import HomeList from '@/views/home/HomeList.vue'
import { useUserStore } from '@/stores/user.ts'
import type { TabPaneName } from 'element-plus'

const userStore = useUserStore()
const loggedIn = computed(() => {
  return userStore.getActiveUser()?.user.id != null
})

const activeNav = ref('recommend')
onMounted(() => {
  if (!loggedIn.value) {
    activeNav.value = 'hot'
  }
})

// 滚动条缓存
const tabScrollTop: Map<string, number> = new Map([
  ['recommend', 0],
  ['hot', 0],
  ['follow', 0],
])

const handleBeforeLeave = (activeName: TabPaneName, oldName: TabPaneName) => {
  // 记录滚动条位置
  console.log(
    'document top',
    document.documentElement.scrollTop,
    'body top',
    document.body.scrollTop,
  )
  tabScrollTop.set(oldName as string, document.documentElement.scrollTop || document.body.scrollTop)
  return true
}
const handleNavChange = (activeName: TabPaneName) => {
  // 恢复滚动条位置
  const scrollTop = tabScrollTop.get(activeName as string)
  if (scrollTop != undefined) {
    nextTick(() => {
      document.documentElement.scrollTop = scrollTop
    })
  }
}
</script>
<template>
  <div class="max-screen-w flex items-start">
    <div class="flex-1 h-full">
      <div class="line-padding z-1">
        <el-tabs
          v-model="activeNav"
          @tab-change="handleNavChange"
          :before-leave="handleBeforeLeave"
        >
          <el-tab-pane label="推荐" name="recommend" lazy v-if="loggedIn">
            <HomeList type="recommend" :show="activeNav == 'recommend'"></HomeList>
          </el-tab-pane>
          <el-tab-pane label="热门" name="hot" lazy>
            <HomeList type="hot" :show="activeNav == 'hot'"></HomeList>
          </el-tab-pane>
          <el-tab-pane label="关注" name="follow" lazy v-if="loggedIn">
            <HomeList type="follow" :show="activeNav == 'follow'"></HomeList>
          </el-tab-pane>
          <!--          <el-tab-pane label="最新" name="latest"></el-tab-pane>-->
        </el-tabs>
      </div>
      <!--      <div class="line-padding relative home-list">-->
      <!--        <InkList></InkList>-->
      <!--        <TwitterStyleInkList :inks="inks" @load-more="loadMore"></TwitterStyleInkList>-->
      <!--      </div>-->
    </div>
    <el-aside width="22rem" class="sticky-top line-padding ml-10 ink-aside hidden-scrollbar">
      <HomeAside></HomeAside>
    </el-aside>
  </div>
</template>
<style scoped lang="scss">
.ink-aside {
  @media (max-width: 1024px) {
    display: none;
  }
}

:deep(.el-tabs__nav-wrap::after) {
  background-color: transparent;
}

:deep(.el-tabs__header) {
  position: sticky;
  top: var(--header-height);
  background-color: var(--white-background);
}
</style>
