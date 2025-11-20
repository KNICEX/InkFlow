<template>
  <div class="flex-col">
    <el-popover></el-popover>
    <div class="flex fixed w-full header-height white-bg z-10 top-0 items-center justify-center">
      <HeaderView class="line-padding"> </HeaderView>
    </div>
    <div class="flex justify-center pt-[var(--header-height)]">
      <router-view v-slot="{ Component }">
        <keep-alive :include="['SearchView', 'HomeView', 'UserView']">
          <Component :is="Component"></Component>
        </keep-alive>
      </router-view>
      <!--      <router-view></router-view>-->
    </div>
    <div class="border-t-1 border-gray-200 mt-10 flex justify-center">
      <FooterView></FooterView>
    </div>
  </div>
</template>
<script setup lang="ts">
import HeaderView from '@/views/header/HeaderView.vue'
import FooterView from '@/views/FooterView.vue'
</script>
<style scoped lang="scss">
.ink-aside {
  @media (max-width: 1024px) {
    display: none;
  }
}
</style>
