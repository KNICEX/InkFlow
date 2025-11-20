<script setup lang="ts">
import { MilkdownProvider } from '@milkdown/vue'
import MilkdownEditor from '@/components/editor/milkdown/MilkdownEditor.vue'
import type { Ctx } from '@milkdown/ctx'
import type { Node as ProseNode } from 'prosemirror-model'
import { useTemplateRef } from 'vue'
defineProps({
  readOnly: {
    type: Boolean,
    default: false,
  },
  padding: {
    type: Boolean,
    default: true,
  },
})
const milkdownRef = useTemplateRef<InstanceType<typeof MilkdownEditor>>('milkdownRef')
const getMarkdown = () => {
  return milkdownRef.value?.getMarkdown()
}
const setContent = (content: string) => {
  return milkdownRef.value?.setContent(content)
}
defineExpose({
  getMarkdown,
  setContent,
})

const emit = defineEmits(['update'])
const handleUpdate = (ctx: Ctx, doc: ProseNode, prevDoc: ProseNode | null) => {
  emit('update', ctx, doc, prevDoc)
}
</script>

<template>
  <div :class="padding ? '' : 'milkdown-no-padding'">
    <MilkdownProvider>
      <MilkdownEditor
        :read-only="readOnly"
        ref="milkdownRef"
        @update="handleUpdate"
      ></MilkdownEditor>
    </MilkdownProvider>
  </div>
</template>

<style lang="scss"></style>
