<script setup lang="ts">
import { Search } from '@element-plus/icons-vue'
import { computed, onMounted, ref, useTemplateRef } from 'vue'
import { useRouter } from 'vue-router'
import InkPopover from '@/components/popover/InkPopover.vue'
import { ElInput } from 'element-plus'
import { useSearchLog } from '@/stores/search.ts'
const props = defineProps({
  placeholder: {
    type: String,
    default: 'Search what you want',
  },
  value: {
    type: String,
    default: '',
  },
})
const router = useRouter()
const searchLogStore = useSearchLog()

const keywords = ref(props.value)
const marshaledKeywords = computed(() => {
  return encodeURIComponent(keywords.value)
})
const handleSearch = () => {
  searchLogStore.addSearchLog(keywords.value)
  router.push({
    name: 'search-ink',
    query: {
      key: marshaledKeywords.value,
    },
  })
}

const inputRef = useTemplateRef<InstanceType<typeof ElInput>>('inputRef')
const popWidth = ref('auto')
onMounted(() => {
  popWidth.value = `${inputRef.value?.$el.offsetWidth}px`
})

const setValue = (value: string) => {
  keywords.value = value
}
defineExpose({
  setValue,
})

const clearSearchLog = () => {
  searchLogStore.clearSearchLogs()
}

const stopPropagation = (evt: Event) => {
  evt.stopPropagation()
}

const searchLog = (key: string) => {
  keywords.value = key
  handleSearch()
}
</script>

<template>
  <InkPopover trigger="click" padding="0" :width="popWidth">
    <template #reference>
      <el-input
        ref="inputRef"
        v-model="keywords"
        size="large"
        clearable
        class="hover:shadow-sm rounded-input h-12"
        :placeholder="placeholder"
        :prefix-icon="Search"
        :style="{ borderRadius: '2rem' }"
        @keyup.enter="handleSearch"
        @keydown.space="stopPropagation"
        @keyup.space="stopPropagation"
      />
    </template>
    <template #content>
      <div>
        <div
          class="px-10 h-12 flex items-center justify-between"
          v-if="searchLogStore.searchLogs.length > 0"
        >
          <span class="lg-semibold-text">历史记录</span>
          <span class="cursor-pointer text-[var(--primary-color)]" @click="clearSearchLog"
            >清除全部</span
          >
        </div>
        <div v-else class="h-14 text-base flex items-center justify-center">
          <span>尝试搜索内容, 人物</span>
        </div>
        <div>
          <div
            class="popover-button"
            v-for="(item, idx) in searchLogStore.searchLogs"
            :key="idx"
            @click="searchLog(item)"
          >
            <span>{{ item }}</span>
            <span></span>
          </div>
        </div>
      </div>
    </template>
  </InkPopover>
</template>

<style scoped lang="scss"></style>
