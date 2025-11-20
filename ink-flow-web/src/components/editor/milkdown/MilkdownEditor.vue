<script setup lang="ts">
import { Milkdown, useEditor } from '@milkdown/vue'
import { Crepe, CrepeFeature } from '@milkdown/crepe'
import { replaceAll } from '@milkdown/kit/utils'
import { commonmark } from '@milkdown/kit/preset/commonmark'
import type { Editor } from '@milkdown/kit/core'
import { upload, uploadConfig, type Uploader } from '@milkdown/plugin-upload'
import { onBeforeUnmount } from 'vue'
import type { Ctx } from '@milkdown/ctx'
import type { Node as ProseNode } from 'prosemirror-model'
import { uploadImage } from '@/service/file.ts'
import '@milkdown/crepe/theme/common/style.css'
import '@milkdown/crepe/theme/frame.css'

import { indent } from '@milkdown/plugin-indent'
const props = defineProps({
  readOnly: {
    type: Boolean,
    default: false,
  },
  defaultValue: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update'])
let crepe: Crepe
const handleUpdate: (ctx: Ctx, doc: ProseNode, prevDoc: ProseNode | null) => void = (
  ctx,
  doc,
  predoc,
) => {
  emit('update', ctx, doc, predoc)
}

const uploader: Uploader = async (files, schema) => {
  const images: File[] = []
  for (const file of files) {
    if (!file) {
      continue
    }
    if (!file.type.includes('image')) {
      continue
    }
    images.push(file)
  }

  return await Promise.all(
    images.map(async (image) => {
      const url = await uploadImage(image)
      return schema.nodes.image.createAndFill({
        src: url,
        alt: image.name,
      })!
    }),
  )
}

let editorInstance: Editor
const waitProcess: (() => void)[] = []

useEditor((root) => {
  crepe = new Crepe({
    root,
    defaultValue: props.defaultValue,
    featureConfigs: {
      [CrepeFeature.Placeholder]: {
        text: 'Type / to use slash command',
      },
      [CrepeFeature.ImageBlock]: {
        onUpload: async (file) => {
          return await uploadImage(file)
        },
      },
    },
  })
  crepe.editor.use(commonmark)
  crepe.on((listener) => {
    listener.updated(handleUpdate)
  })
  crepe.editor.use(upload).use(indent)
  crepe.editor.config((ctx) => {
    ctx.update(uploadConfig.key, (prev) => ({
      ...prev,
      uploader,
      enableHtmlFileUploader: true,
    }))
  })

  crepe.create().then((editor) => {
    editorInstance = editor
    waitProcess.forEach((fn) => {
      fn()
    })
  })
  crepe.setReadonly(props.readOnly)
  return crepe
})

const setContent = (content: string) => {
  if (editorInstance) {
    editorInstance.action(replaceAll(content))
  } else {
    waitProcess.push(() => {
      editorInstance.action(replaceAll(content))
    })
  }
}

const getMarkdown = () => {
  return crepe.getMarkdown()
}
defineExpose({
  getMarkdown,
  setContent,
})

onBeforeUnmount(() => {
  crepe.destroy()
})
</script>

<template>
  <Milkdown class="prose"></Milkdown>
</template>

<style scoped lang="scss"></style>
