<script setup lang="ts">
import InkList from '@/components/list/ink/InkList.vue'
import type { Ink } from '@/types/ink.ts'
import { list } from '@/service/ink.ts'
import { ref, watch } from 'vue'
import { wrapOffsetPagedFunc } from '@/utils/pagedLoadWrap.ts'
const props = defineProps({
  uid: {
    type: Number,
    default: 0,
  },
})

const limit = 15
const inks = ref<Ink[]>([])

const { loadMore, reset, loading } = wrapOffsetPagedFunc(async (offset: number) => {
  if (props.uid == 0) {
    return
  }
  console.log('user ink list load more', offset)
  const res = await list({
    authorId: props.uid,
    offset: offset,
    limit: limit,
  })
  inks.value = [...inks.value, ...res]
  return res.length
}, limit)

watch(
  () => props.uid,
  () => {
    console.log('user ink list props uid changed', props.uid)
    if (!props.uid) {
      inks.value = []
      reset()
      return
    }

    console.log('user ink list props uid changed, load more')
    inks.value = []
    reset()
    loadMore(props.uid.toString())
  },
  { immediate: true },
)
</script>

<template>
  <div class="w-full">
    <InkList
      :max-cols="4"
      :inks="inks"
      :load-more="() => loadMore(uid.toString())"
      :loading="loading"
    ></InkList>
  </div>
</template>

<style scoped lang="scss"></style>
