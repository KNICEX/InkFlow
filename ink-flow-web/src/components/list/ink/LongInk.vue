<script setup lang="ts">
import type { Ink } from '@/types/ink.ts'
import InkCard from '@/components/card/InkCard.vue'

defineProps({
  ink: {
    type: Object as () => Ink,
    required: true,
  },
})

const handleCoverClick = () => {}
</script>

<template>
  <div>
    <InkCard>
      <el-image
        v-if="ink.cover != ''"
        @click="handleCoverClick"
        class="rounded-xl w-full aspect-square"
        loading="lazy"
        fit="cover"
        :src="ink.cover"
      >
      </el-image>
      <div
        v-else
        @click="handleCoverClick"
        class="flex justify-center items-center rounded-xl w-full aspect-square"
      >
        <div>{{ ink.title }}</div>
      </div>
    </InkCard>
  </div>
</template>

<style scoped lang="scss"></style>
