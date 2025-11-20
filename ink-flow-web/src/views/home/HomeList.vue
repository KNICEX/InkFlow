<script setup lang="ts">
import { onMounted, type Ref, ref } from 'vue'
import TwitterStyleInkList from '@/components/list/ink/TwitterStyleInkList.vue'
import { wrapMaxIdTimestampPagedFunc, wrapOffsetPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { followInks, hotInks, recommendInks } from '@/service/feed.ts'
import type { Ink } from '@/types/ink.ts'
import { useProvideInkInteractiveHandler } from '@/hook/interactive.ts'

const props = defineProps({
  type: {
    type: String as () => 'recommend' | 'follow' | 'hot',
    required: true,
  },
  show: {
    type: Boolean,
    default: true,
  },
})

let loadMore: () => Promise<boolean | undefined>

const inks = ref<Ink[]>([])
let loading: Ref<boolean>
useProvideInkInteractiveHandler(inks)
const limit = 25
switch (props.type) {
  case 'recommend':
    const { loadMore: loadRecommend, loading: recommendLoading } = wrapOffsetPagedFunc(
      async (offset: number) => {
        const res = await recommendInks(offset, limit)
        if (res.length == 0) {
          return 0
        }
        inks.value = [...inks.value, ...res]
        return res.length
      },
      limit,
    )
    loadMore = loadRecommend
    loading = recommendLoading
    break
  case 'hot':
    const { loadMore: loadHot, loading: hotLoading } = wrapOffsetPagedFunc(
      async (offset: number) => {
        const res = await hotInks(offset, limit)
        if (res.length == 0) {
          return 0
        }
        inks.value = [...inks.value, ...res]
        return res.length
      },
      limit,
    )
    loadMore = loadHot
    loading = hotLoading
    break
  case 'follow':
    const { loadMore: loadFollow, loading: followLoading } = wrapMaxIdTimestampPagedFunc(
      async (maxId: string, timestamp: number) => {
        const res = await followInks(maxId, timestamp, limit)
        if (res.length == 0) {
          return {
            maxId: '0',
            timestamp: new Date(),
          }
        }
        inks.value = [...inks.value, ...res]
        return {
          maxId: res[res.length - 1].id,
          timestamp: res[res.length - 1].createdAt,
        }
      },
    )
    loadMore = loadFollow
    loading = followLoading
    break
  default:
    console.error('unknown type', props.type)
}
onMounted(() => {
  loadMore()
})
</script>

<template>
  <div>
    <TwitterStyleInkList
      :inks="inks"
      :show="show"
      :load-more="loadMore"
      :loading="loading"
    ></TwitterStyleInkList>
  </div>
</template>

<style scoped lang="scss"></style>
