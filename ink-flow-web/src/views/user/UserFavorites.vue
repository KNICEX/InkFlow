<script setup lang="ts">
import { onBeforeMount, ref, watch } from 'vue'
import InkList from '@/components/list/ink/InkList.vue'
import { getFavorites } from '@/service/interactive.ts'
import type { Favorite } from '@/types/interactive.ts'
import type { Ink } from '@/types/ink.ts'
import { listFavorited } from '@/service/ink.ts'
const activeIndex = ref(0)
const favorites = ref<Favorite[]>([])
const inks = ref<Ink[]>([])
onBeforeMount(async () => {
  favorites.value = await getFavorites('ink')
})

watch(
  activeIndex,
  async () => {
    inks.value = await listFavorited({
      fid: activeIndex.value,
      maxId: '0',
      limit: 25,
    })
  },
  { immediate: true },
)
</script>

<template>
  <div class="flex w-full">
    <div class="w-60">
      <div class="sticky top-44 pr-4">
        <el-menu :default-active="activeIndex">
          <el-menu-item :index="0">
            <div class="flex justify-between items-center w-full">
              <span class="text-base">默认收藏夹</span>
              <span></span>
            </div>
          </el-menu-item>
          <el-menu-item v-for="fav in favorites" :index="fav.id" :key="fav.id">
            <div class="flex justify-between items-center w-full">
              <span class="text-base">{{ fav.name }}</span>
              <span>2</span>
            </div>
          </el-menu-item>
        </el-menu>
      </div>
    </div>
    <div class="flex-1">
      <InkList :max-cols="4" :inks="inks"></InkList>
    </div>
  </div>
</template>

<style scoped lang="scss">
.el-menu {
  border: none;
}
.el-menu-item {
  border-radius: 0.8rem;
}
.is-active {
  background-color: var(--el-color-primary-light-9);
}
</style>
