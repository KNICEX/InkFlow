<script setup lang="ts">
import DashboardNav from '@/views/dashboard/DashboardNav.vue'
</script>

<template>
  <div class="flex w-full h-screen bg-slate-100 dark:bg-gray-800">
    <el-popover></el-popover>
    <DashboardNav> </DashboardNav>

    <router-view></router-view>
  </div>
</template>

<style scoped lang="scss"></style>
