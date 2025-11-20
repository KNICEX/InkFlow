<script setup lang="ts">
import { ref } from 'vue'
import InkList from '@/components/list/ink/InkList.vue'
import type { Ink } from '@/types/ink.ts'
import { onBeforeMount } from 'vue'
import { listViewed } from '@/service/ink.ts'
import { wrapMaxIdPagedFunc } from '@/utils/pagedLoadWrap.ts'
const inks = ref<Ink[]>([])

const limit = 15
const { loadMore, loading } = wrapMaxIdPagedFunc(async (maxId: string) => {
  const res = await listViewed({
    maxId,
    limit,
  })
  if (res.length == 0) {
    return '0'
  }
  inks.value = [...inks.value, ...res]
  return inks.value[inks.value.length - 1].id
})

onBeforeMount(() => {
  loadMore()
})
</script>

<template>
  <div>
    <InkList :inks="inks" :max-cols="4" :loading="loading"></InkList>
  </div>
</template>

<style scoped lang="scss"></style>
