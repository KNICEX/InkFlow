<script setup lang="ts">
import { provide, ref, watch } from 'vue'
import DashboardContent from '@/views/dashboard/DashboardContent.vue'
import DashboardInkList from '@/components/list/dashboard/DashboardInkList.vue'
import { type Ink, InkStatus, inkStatusProp } from '@/types/ink.ts'
import {
  listDraft,
  listRejected,
  listPending,
  listPrivate,
  list,
  deleteLive,
  deleteDraft,
  deletePrivate,
  publish,
} from '@/service/ink.ts'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user.ts'
import { notification } from '@/utils/notification.ts'
import { confirm } from '@/utils/message.ts'
import { parseRouteParam } from '@/utils/parse.ts'
import { wrapOffsetPagedFunc } from '@/utils/pagedLoadWrap.ts'
const title = ref('Ink')
const activeType = ref('published')
const inks = ref<Ink[]>([])
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const limit = 15
const load = async (offset: number) => {
  let res: Ink[] = []
  switch (activeType.value) {
    case 'published':
      res = await list({
        authorId: userStore.getActiveUser()?.user.id ?? '0',
        offset: offset,
        limit: limit,
      })
      break
    case 'pending':
      res = await listPending({
        offset: offset,
        limit: limit,
      })
      break
    case 'rejected':
      res = await listRejected({
        offset: offset,
        limit: limit,
      })
      break
    case 'draft':
      res = await listDraft({
        offset: offset,
        limit: limit,
      })
      break
    case 'private':
      res = await listPrivate({
        offset: offset,
        limit: limit,
      })
      break
    default:
      break
  }

  if (res.length == 0) {
    return 0
  }
  inks.value.push(...res)
  return res.length
}

const { loadMore, reset: resetPage } = wrapOffsetPagedFunc(async (offset: number) => {
  return await load(offset)
}, limit)
watch(
  () => route.params,
  () => {
    let type = parseRouteParam(route.params.status)
    if (!type || type === '') {
      type = 'published'
    }
    activeType.value = type
    resetPage()
    inks.value = []
    loadMore()
  },
  { immediate: true },
)

const handleTabChange = (tab: string) => {
  router.push({
    name: 'dashboard-ink',
    params: {
      status: tab,
    },
  })
}

const handleNewInk = () => {
  router.push({
    name: 'editor',
  })
}

const confirmBox = (onConfirm: () => void) => {
  confirm({
    title: 'warning',
    message: '删除后无法恢复，确认删除吗😰?',
    confirmed: onConfirm,
  })
}
const delSuccess = () => {
  notification({
    type: 'success',
    title: 'Success',
    message: '内容删除成功',
  })
}

const handleInkOp = async (ink: Ink, op: 'delete' | 'preview' | 'edit' | 'publish') => {
  switch (op) {
    case 'delete':
      switch (ink.status) {
        case InkStatus.Published:
          confirmBox(async () => {
            await deleteLive(ink.id)
            delSuccess()
          })
          break
        case InkStatus.UnPublished:
          confirmBox(async () => {
            await deleteDraft(ink.id)
            delSuccess()
          })
          break
        case InkStatus.Private:
          await deletePrivate(ink.id)
          break
      }
      break
    case 'preview':
      await router.push({
        name: 'ink-detail',
        params: {
          id: ink.id,
        },
        query: {
          status: inkStatusProp(ink.status),
        },
      })
      break
    case 'edit':
      console.log('ink', ink)
      await router.push({
        name: 'editor',
        params: {
          id: ink.id,
        },
      })
      break
    case 'publish':
      await publish(ink.id)
      await router.push({
        name: 'dashboard-ink',
        params: {
          status: inkStatusProp(InkStatus.Pending),
        },
      })
  }
}

provide('handle-ink-op', handleInkOp)
</script>

<template>
  <DashboardContent :title="title">
    <div class="flex overflow-hidden flex-1 flex-col">
      <div>
        <el-button class="w-30" size="large" round type="primary" @click="handleNewInk"
          >新建文章</el-button
        >
      </div>
      <div class="pt-4 text-lg">
        <el-tabs v-model="activeType" @tab-change="handleTabChange">
          <el-tab-pane label="已发布" name="published"> </el-tab-pane>
          <el-tab-pane label="审核中" name="pending"> </el-tab-pane>
          <el-tab-pane label="审核失败" name="rejected"> </el-tab-pane>
          <el-tab-pane label="草稿" name="draft"> </el-tab-pane>
          <el-tab-pane label="私有" name="private"> </el-tab-pane>
        </el-tabs>
      </div>
      <DashboardInkList :load-more="loadMore" :inks="inks" class="flex-1 overflow-y-auto">
      </DashboardInkList>
    </div>
  </DashboardContent>
</template>

<style scoped lang="scss"></style>
