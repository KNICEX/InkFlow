<script setup lang="ts">
import { onBeforeMount } from 'vue'
import { useUserStore } from '@/stores/user.ts'
import { profile } from '@/service/user.ts'
import { useThemeStore } from '@/stores/theme.ts'
const userStore = useUserStore()
const themeStore = useThemeStore()
onBeforeMount(async () => {
  userStore.loadFromLocalStorage()
  themeStore.loadFromLocalStorage()
  const activeUser = userStore.getActiveUser()
  if (activeUser) {
    const userInfo = await profile()
    userStore.updateActiveUser(userInfo)
  }
})
</script>

<template>
  <router-view> </router-view>
</template>
