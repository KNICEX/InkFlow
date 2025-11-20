<script setup lang="ts">
import InkList from '@/components/list/ink/InkList.vue'
import { ref, watch } from 'vue'
import type { Ink } from '@/types/ink.ts'
import { wrapOffsetPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { searchInk } from '@/service/search.ts'
const props = defineProps({
  keyword: {
    type: String,
    default: '',
  },
})

const inks = ref<Ink[]>([])
const limit = 15

const { loadMore, reset, loading } = wrapOffsetPagedFunc(async (offset: number) => {
  console.log('load more ink', offset)
  if (props.keyword == '') {
    return
  }
  const res = await searchInk({
    keyword: props.keyword,
    offset,
    limit,
  })
  if (res.length == 0) {
    return 0
  }
  inks.value = [...inks.value, ...res]
  return res.length
}, limit)

watch(
  () => props.keyword,
  () => {
    console.log('search keyword change', props.keyword)
    if (props.keyword == '') {
      inks.value = []
    } else {
      reset()
      inks.value = []
      loadMore()
    }
  },
  { immediate: true },
)
</script>

<template>
  <div>
    <InkList
      :load-more="loadMore"
      :loading="loading"
      :inks="inks"
      :max-cols="5"
      empty-desc="没有相关结果"
    ></InkList>
  </div>
</template>

<style scoped lang="scss"></style>
