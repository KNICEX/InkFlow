<script setup lang="ts">
import InkLogo from '@/components/icons/InkLogo.vue'
import UserPanel from '@/views/header/UserPanel.vue'
import InkNav from '@/components/nav/InkNav.vue'
import InkNavItem from '@/components/nav/InkNavItem.vue'
import { useRouter } from 'vue-router'
import { onMounted, useTemplateRef } from 'vue'
import { useUnreadNotificationStore } from '@/stores/notification.ts'

const router = useRouter()
const route = router.currentRoute
const inkNavRef = useTemplateRef<InstanceType<typeof InkNav>>('inkNavRef')
const handleNavChange = (nav: string) => {
  router.push({ name: nav, params: route.value.params })
}
onMounted(() => {
  inkNavRef.value?.setActive(route.value.name as string)
})

const noStore = useUnreadNotificationStore()
</script>

<template>
  <div class="flex flex-col w-70 items-start px-3 h-full">
    <div class="w-full">
      <div class="flex pl-2 items-center h-24 ml-4">
        <InkLogo></InkLogo>
      </div>
      <div class="flex pl-2 items-center">
        <UserPanel></UserPanel>
      </div>
    </div>

    <div class="w-full flex-1 overflow-y-auto mt-4">
      <InkNav ref="inkNavRef" default-active="dashboard-home" @on-change="handleNavChange">
        <InkNavItem prop="dashboard-home">
          <span class="material-symbols-outlined"> space_dashboard </span>
          <span class="ml-3">仪表盘</span>
        </InkNavItem>
        <InkNavItem :selectable="false">
          <span class="material-symbols-outlined"> add_circle </span>
          <span>发表</span>
        </InkNavItem>
        <InkNavItem :sub="true" prop="editor">
          <span class="material-symbols-outlined"> edit_note </span>
          <span class="ml-3">新建</span>
        </InkNavItem>
        <InkNavItem :sub="true" prop="dashboard-ink">
          <span class="material-symbols-outlined"> apps </span>
          <span class="ml-3">内容</span>
        </InkNavItem>
        <!--        <InkNavItem :sub="true" prop="new_share">-->
        <!--          <span class="material-symbols-outlined"> smb_share </span>-->
        <!--          <span class="ml-3">发布资源</span>-->
        <!--        </InkNavItem>-->
        <InkNavItem :selectable="false">
          <span class="material-symbols-outlined"> bolt </span>
          <span>互动</span>
        </InkNavItem>
        <InkNavItem :sub="true" prop="dashboard-system">
          <el-badge type="primary" :value="noStore.unreadMap['system']" :show-zero="false">
            <span class="material-symbols-outlined"> notifications </span>
          </el-badge>
          <span class="ml-3">系统通知</span>
        </InkNavItem>
        <InkNavItem :sub="true" prop="dashboard-reply">
          <el-badge
            type="primary"
            :value="noStore.unreadMap['reply']"
            :show-zero="false"
            class="flex flex-col justify-center"
          >
            <span class="material-symbols-outlined"> chat </span>
          </el-badge>
          <span class="ml-3">回复我的</span>
        </InkNavItem>
        <InkNavItem :sub="true" prop="dashboard-like">
          <el-badge type="primary" :value="noStore.unreadMap['like']" :show-zero="false">
            <span class="material-symbols-outlined"> thumb_up </span>
          </el-badge>
          <span class="ml-3">收到的赞</span>
        </InkNavItem>
        <InkNavItem :sub="true" prop="dashboard-follow">
          <el-badge type="primary" :value="noStore.unreadMap['follow']" :show-zero="false">
            <span class="material-symbols-outlined"> how_to_reg </span>
          </el-badge>
          <span class="ml-3">关注我的</span>
        </InkNavItem>
        <!--        <InkNavItem prop="setting">-->
        <!--          <span class="material-symbols-outlined"> manufacturing </span>-->
        <!--          <span class="ml-3">设置</span>-->
        <!--        </InkNavItem>-->
      </InkNav>
    </div>
    <!--    <div class="w-full h-20 flex items-center justify-center border border-solid bt">Footer</div>-->
  </div>
</template>

<style scoped lang="scss"></style>
