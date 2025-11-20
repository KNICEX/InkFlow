<script setup lang="ts">
import { onMounted, useTemplateRef } from 'vue'
import MilkdownWrapper from '@/components/editor/milkdown/MilkdownWrapper.vue'
import { emptyInk, type Ink } from '@/types/ink.ts'
const milkdownRef = useTemplateRef<InstanceType<typeof MilkdownWrapper>>('milkdownRef')
const props = defineProps({
  ink: {
    type: Object as () => Ink,
    default: emptyInk(),
  },
  inkId: {
    type: Number,
    default: 0,
  },
})
onMounted(() => {
  milkdownRef.value?.setContent(props.ink.contentMeta)
})
</script>

<template>
  <div>
    <!--    <MilkdownWrapper :read-only="true" ref="milkdownRef"> </MilkdownWrapper>-->
    <div v-html="ink.contentHtml"></div>
  </div>
</template>

<style scoped lang="scss"></style>
