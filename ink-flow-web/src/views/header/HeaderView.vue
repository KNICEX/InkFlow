<template>
  <div class="max-screen-w flex header bg items-center justify-between">
    <div class="flex items-center">
      <InkLogo></InkLogo>
      <div class="ml-5 text-lg">
        <router-link to="/" class="ml-10">主页</router-link>
        <span class="ml-10">关于</span>
      </div>
    </div>
    <div>
      <div v-if="userStore.getActiveUser()?.user?.id" class="flex items-center">
        <div class="flex items-center text-gray-500">
          <InkPopover place="bottom" padding="0">
            <template #reference>
              <div class="relative flex items-center">
                <span class="material-symbols-outlined cursor-pointer"> notifications </span>
                <span v-show="true" class="absolute flex -right-1 -top-1 size-3">
                  <span
                    class="absolute inline-flex h-full w-full animate-ping rounded-full bg-sky-400 opacity-75"
                  ></span>
                  <span class="relative inline-flex size-3 rounded-full bg-sky-500"></span>
                </span>
              </div>
            </template>
            <template #content>
              <NotificationPop></NotificationPop>
            </template>
          </InkPopover>
          <router-link
            title="新建"
            :to="`/dashboard/${userStore.getActiveUser()?.user.account}`"
            class="material-symbols-outlined ml-3 cursor-pointer"
          >
            <span @click="toDashboard"> add_circle </span>
          </router-link>
        </div>
        <el-divider direction="vertical" style="height: 2rem" border-style="solid" />
        <UserPanel></UserPanel>
      </div>
      <div v-else>
        <el-button round size="large" class="w-28" color="black" @click="handleLogin">
          登录
        </el-button>
      </div>
    </div>
    <LoginView v-model="loginStore.showLogin"></LoginView>
  </div>
</template>
<script setup lang="ts">
import { useShowLoginStore, useUserStore } from '@/stores/user.ts'
import LoginView from '@/views/login/LoginView.vue'
import InkPopover from '@/components/popover/InkPopover.vue'
import NotificationPop from '@/views/header/NotificationPop.vue'
import { useRouter } from 'vue-router'
import InkLogo from '@/components/icons/InkLogo.vue'
import UserPanel from '@/views/header/UserPanel.vue'
const router = useRouter()
const userStore = useUserStore()

const loginStore = useShowLoginStore()
const handleLogin = () => {
  loginStore.showLogin = true
}

const toDashboard = () => {
  router.push(`/dashboard/${userStore.getActiveUser()?.user.account}`)
}
</script>
<style scoped lang="scss"></style>
