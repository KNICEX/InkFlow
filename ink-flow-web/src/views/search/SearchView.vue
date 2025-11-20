<script setup lang="ts">
import InkSearch from '@/components/InkSearch.vue'
import { useRoute, useRouter } from 'vue-router'
import { nextTick, useTemplateRef, watch, ref } from 'vue'
import { parseRouteQuery } from '@/utils/parse.ts'
import type { TabPaneName } from 'element-plus'
const route = useRoute()
const router = useRouter()

const searchInputRef = useTemplateRef<InstanceType<typeof InkSearch>>('searchInputRef')

const searchTarget = ref('ink')
const searchKey = ref('')

const scrollTopMap = new Map<string, number>([
  ['ink', 0],
  ['user', 0],
  ['reply', 0],
])

const handleBeforeLeave = (active: TabPaneName, old: TabPaneName) => {
  // 记录滚动条位置
  scrollTopMap.set(old as string, document.documentElement.scrollTop || document.body.scrollTop)
  return true
}

const handleTargetChange = (active: TabPaneName) => {
  // 恢复滚动条位置
  const scrollTop = scrollTopMap.get(active as string)

  router
    .push({
      name: `search-${active as string}`,
      query: {
        key: route.query.key,
      },
    })
    .then(() => {
      if (scrollTop != undefined) {
        nextTick(() => {
          document.documentElement.scrollTop = scrollTop
        })
      }
    })
}

const reset = () => {
  scrollTopMap.set('ink', 0)
  scrollTopMap.set('user', 0)
  scrollTopMap.set('reply', 0)
}

watch(
  () => route.query.key,
  () => {
    if (route.query.key == undefined) {
      return
    }
    const key = decodeURIComponent(parseRouteQuery(route.query.key))
    let type = parseRouteQuery(route.query.type)

    reset()

    if (type == '') {
      type = 'ink'
    }

    searchTarget.value = type
    searchKey.value = key

    if (searchInputRef.value) {
      searchInputRef.value.setValue(key)
    } else {
      nextTick(() => {
        searchInputRef.value?.setValue(key)
      })
    }
  },
  {
    immediate: true,
  },
)
</script>

<template>
  <div class="max-screen-w line-padding">
    <div class="mb-2 sticky-top z-1 white-bg">
      <div class="max-w-140 mb-4"><InkSearch ref="searchInputRef"></InkSearch></div>
      <div class="sticky-top white-bg z-10">
        <el-tabs
          v-model="searchTarget"
          @tab-change="handleTargetChange"
          :before-leave="handleBeforeLeave"
        >
          <!--        <el-tab-pane label="热门" name="hot"></el-tab-pane>-->
          <el-tab-pane lazy label="内容" name="ink"> </el-tab-pane>
          <el-tab-pane lazy label="用户" name="user"> </el-tab-pane>
          <el-tab-pane lazy label="回复" name="reply"> </el-tab-pane>
        </el-tabs>
      </div>
      <div class="w-full">
        <router-view v-slot="{ Component }">
          <keep-alive>
            <Component :is="Component" :keyword="searchKey"></Component>
          </keep-alive>
        </router-view>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>
