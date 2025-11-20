<script setup lang="ts">
import TwitterStyleInk from '@/components/list/ink/TwitterStyleInk.vue'
import type { Ink } from '@/types/ink.ts'
import BackTop from '@/components/BackTop.vue'
import NoData from '@/components/empty/NoData.vue'
import InkLoading from '@/components/loading/InkLoading.vue'
import { useActive } from '@/hook/active.ts'
import { watchEffect } from 'vue'

const props = defineProps({
  inks: {
    type: Array as () => Ink[],
    default: () => [],
  },
  previewLength: {
    type: Number,
    default: 150,
  },
  loadMore: {
    type: Function,
    default: () => {},
  },
  loading: {
    type: Boolean,
    default: false,
  },
  show: {
    type: Boolean,
    default: true,
  },
})

const active = useActive()

watchEffect(() => {
  active.value = props.show
})
</script>

<template>
  <div>
    <div
      v-infinite-scroll="loadMore"
      :infinite-scroll-disabled="!active"
      :infinite-scroll-distance="1800"
    >
      <div v-for="ink in inks" :key="ink.id" class="mb-4">
        <TwitterStyleInk v-slide-in :ink="ink" :previewLength="previewLength"></TwitterStyleInk>
      </div>
      <NoData v-show="!loading && inks.length == 0"></NoData>
      <InkLoading v-show="loading" :loading="loading"></InkLoading>
    </div>
    <BackTop></BackTop>
  </div>
</template>

<style scoped lang="scss"></style>
