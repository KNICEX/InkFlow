<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue'
import MilkdownWrapper from '@/components/editor/milkdown/MilkdownWrapper.vue'
import { draftDetail, publish, saveDraft } from '@/service/ink.ts'
import { useRoute, useRouter } from 'vue-router'
import { notification } from '@/utils/notification.ts'
import { parseRouteParam } from '@/utils/parse.ts'
import DashboardContent from '@/views/dashboard/DashboardContent.vue'
import { InkStatus, inkStatusProp } from '@/types/ink.ts'

const route = useRoute()
const router = useRouter()
const milkdownRef = useTemplateRef<InstanceType<typeof MilkdownWrapper>>('milkdownRef')

const coverUrl = ref('')
let draftId = parseRouteParam(route.params.id)
const tags = ref<string[]>([])
const title = ref('')
let lastSave = Date.now()
let saveClicked = false
onMounted(async () => {
  if (draftId != '') {
    const draft = await draftDetail(draftId)
    milkdownRef.value?.setContent(draft.contentMeta)
    coverUrl.value = draft.cover
    tags.value = draft.tags
    title.value = draft.title
  }
})

const save = async () => {
  const markdown = milkdownRef.value?.getMarkdown() ?? ''
  if (markdown == '') {
    notification({
      title: 'Error',
      message: '内容不能为空',
      type: 'error',
    })
    return
  }
  const { id } = await saveDraft({
    id: draftId == '' ? '0' : draftId,
    title: title.value,
    contentMeta: milkdownRef.value?.getMarkdown() || '',
    contentHtml: '',
    tags: tags.value,
    cover: coverUrl.value,
  })
  draftId = id
  router.replace({
    name: 'editor',
    params: {
      id: id,
    },
  })
  lastSave = Date.now()
}
const handleUpdate = async () => {
  console.log('update')
  if (saveClicked && Date.now() - lastSave > 1000 * 60) {
    await save()
  }
}
const handleSave = async () => {
  saveClicked = true
  // cover为空，设置第一个图片为cover
  if (coverUrl.value == '') {
    const markdown = milkdownRef.value?.getMarkdown() ?? ''
    const firstImage = markdown.match(/!\[.*?\]\((.*?)\)/)
    if (firstImage) {
      coverUrl.value = firstImage[1]
    }
  }

  // 标题为空, 设置第一段文字为标题
  if (title.value == '') {
    const markdown = milkdownRef.value?.getMarkdown() ?? ''
    title.value = markdown.split('\n')[0].replace(/[#*`~]/g, '')
  }

  await save()
  notification({
    message: '保存成功',
  })
}
const handlePublish = async () => {
  await handleSave()
  await publish(draftId)
  notification({
    message: '发布成功',
  })
  await router.push({
    name: 'dashboard-ink',
    params: {
      status: inkStatusProp(InkStatus.Pending),
    },
  })
}
</script>

<template>
  <DashboardContent title="Edit">
    <div class="flex w-full flex-1 grow-1 overflow-hidden">
      <div class="w-40 h-full hidden 2xl:inline">
        <span> 自动目录区 </span>
        <span>待实现</span>
      </div>
      <MilkdownWrapper
        :read-only="false"
        @update="handleUpdate"
        ref="milkdownRef"
        class="flex-1 overflow-y-auto overflow-x-visible"
      ></MilkdownWrapper>
      <div class="w-80 h-full pl-4 hidden xl:inline">
        <!--        <div>-->
        <!--          <div class="label-text">封面</div>-->
        <!--          <el-upload-->
        <!--            class="w-full bg-gray-50 h-50 upload mt-3 overflow-hidden relative"-->
        <!--            action="#"-->
        <!--            list-type="picture-card"-->
        <!--            :show-file-list="false"-->
        <!--            :on-success="handleCoverUploadSuccess"-->
        <!--            :before-upload="beforeCoverUpload"-->
        <!--          >-->
        <!--            <el-image-->
        <!--              v-if="coverUrl"-->
        <!--              fit="cover"-->
        <!--              :src="coverUrl"-->
        <!--              class="w-full h-full"-->
        <!--              alt="cover"-->
        <!--            />-->
        <!--            <div class="absolute w-4 h-4 right-2 top-2" v-if="coverUrl">-->
        <!--              <span class="material-symbols-outlined"> close </span>-->
        <!--            </div>-->
        <!--            <div v-else class="h-full flex justify-center items-center nav-text">点击上传封面</div>-->
        <!--          </el-upload>-->
        <!--        </div>-->
        <div class="my-6">
          <div class="label-text mb-3">标签</div>
          <el-input-tag size="large" v-model="tags" draggable placeholder="" />
          <div class="text-sm mt-2 text-gray-400 dark:text-gray-500">输入#标签 并回车</div>
        </div>
        <div class="my-6">
          <div class="label-text">标题</div>
          <el-input
            class="mt-3"
            v-model="title"
            :autosize="{ minRows: 2, maxRows: 4 }"
            type="textarea"
            placeholder="Please input"
          />
        </div>
        <div class="flex flex-col items-center">
          <el-button type="primary" size="large" class="w-full mt-6" @click="handleSave"
            >保存草稿</el-button
          >
          <el-button type="primary" size="large" class="w-full mt-6" @click="handlePublish"
            >发布</el-button
          >
        </div>
        <div>1分钟前自动保存</div>
      </div>
    </div>
  </DashboardContent>
</template>

<style scoped lang="scss">
.upload {
  &:deep(.el-upload) {
    @apply w-full;
    @apply h-full;
  }
  &:deep(.el-upload-dragger) {
    @apply w-full;
    @apply h-full;
  }
}
</style>
