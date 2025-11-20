<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { Ink } from '@/types/ink.ts'
import InkList from '@/components/list/ink/InkList.vue'
import { wrapMaxIdPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { listLiked } from '@/service/ink.ts'
import BackTop from '@/components/BackTop.vue'

const inks = ref<Ink[]>([])

const limit = 15
const { loadMore, loading } = wrapMaxIdPagedFunc(async (maxId: string) => {
  const res = await listLiked({
    maxId,
    limit,
  })
  if (res.length == 0) {
    return '0'
  }
  inks.value = [...inks.value, ...res]
  if (res.length < limit) {
    return '0'
  }
  return inks.value[inks.value.length - 1].id
})
onMounted(() => {
  loadMore()
})
</script>

<template>
  <div>
    <InkList :inks="inks" :max-cols="4" :load-more="loadMore" :loading="loading"></InkList>
  </div>
</template>

<style scoped lang="scss"></style>
