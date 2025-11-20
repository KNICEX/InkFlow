<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { computed, h, nextTick, ref, render, useTemplateRef, watch } from 'vue'
import { parseRouteParam, parseRouteQuery } from '@/utils/parse.ts'
import { emptyInk, type Ink, InkStatus, inkStatusFromProp } from '@/types/ink.ts'
import {
  cancelFavorite,
  cancelLike,
  detail,
  draftDetail,
  favorite,
  like,
  pendingDetail,
  privateDetail,
  rejectedDetail,
} from '@/service/ink.ts'
import MilkdownWrapper from '@/components/editor/milkdown/MilkdownWrapper.vue'
import { notification } from '@/utils/notification.ts'
import UserCard from '@/components/UserCard.vue'
import '@milkdown/crepe/theme/common/style.css'
import '@milkdown/crepe/theme/frame.css'
import RecommendCard from '@/components/ink/RecommendCard.vue'
import InkInteractive from '@/components/ink/InkInteractive.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import InkPopover from '@/components/popover/InkPopover.vue'
import { formatDate } from '@/utils/date.ts'
import { ElImage } from 'element-plus'
import CommentView from '@/views/ink/CommentView.vue'
import MoreOperation from '@/components/button/MoreOperation.vue'
import BackTop from '@/components/BackTop.vue'
import { useUserStore } from '@/stores/user.ts'
import { similarInks } from '@/service/recommend.ts'

const milkdownRef = useTemplateRef<InstanceType<typeof MilkdownWrapper>>('milkdownRef')
const contentRef = useTemplateRef<HTMLElement>('contentRef')
const ink = ref<Ink>(emptyInk())
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 处理图片预览
const replaceImage = () => {
  const images = contentRef.value?.querySelectorAll('img')
  if (images) {
    images.forEach((img) => {
      const vNode = h(ElImage, {
        src: img.src,
        fit: 'cover',
        hideOnClickModal: true,
        previewSrcList: [img.src],
        class: 'rounded-image',
      })
      // 渲染到img之前
      render(vNode, img.parentNode as HTMLElement)
      img.remove()
    })
  }
}

const loading = ref(false)
const loadInk = async (id: string, status: InkStatus) => {
  loading.value = true
  switch (status) {
    case InkStatus.Published:
      ink.value = await detail(id)
      break
    case InkStatus.UnPublished:
      ink.value = await draftDetail(id)
      break
    case InkStatus.Private:
      ink.value = await privateDetail(id)
      break
    case InkStatus.Pending:
      ink.value = await pendingDetail(id)
      break
    case InkStatus.Rejected:
      ink.value = await rejectedDetail(id)
      break
    default:
      loading.value = false
      notification({
        title: 'Error',
        message: 'Ink status not supported',
        type: 'error',
      })
      return
  }

  similarInks(ink.value.id, 0, 5).then((si) => {
    recommendInks.value = si
  })

  console.log(ink.value)
  milkdownRef.value?.setContent(ink.value.contentMeta)

  nextTick(() => {
    replaceImage()
    loading.value = false
  })
}

watch(
  () => route.params,
  async () => {
    const inkId = parseRouteParam(route.params.id)
    const statusStr = parseRouteQuery(route.query.status)
    let inkStatus = InkStatus.Published
    if (statusStr != '') {
      inkStatus = inkStatusFromProp(statusStr)
    }
    await loadInk(inkId, inkStatus).catch(() => {
      loading.value = false
    })
  },
  { immediate: true },
)

// TODO 动态加载
const recommendInks = ref<Ink[]>([])

const handleLike = async () => {
  await like(ink.value.id)
  ink.value.interactive.likeCnt = ink.value.interactive.likeCnt + 1
  ink.value.interactive.liked = !ink.value.interactive.liked
}

const handleFavorite = async () => {
  await favorite(ink.value.id, '0')
  ink.value.interactive.favoriteCnt = ink.value.interactive.favoriteCnt + 1
  ink.value.interactive.favorited = !ink.value.interactive.favorited
}

const handleCancelLike = async () => {
  await cancelLike(ink.value.id)
  ink.value.interactive.likeCnt = Math.max(0, ink.value.interactive.likeCnt - 1)
  ink.value.interactive.liked = !ink.value.interactive.liked
}

const handleCancelFavorite = async () => {
  console.log('interactive: ', ink.value.interactive)
  await cancelFavorite(ink.value.id)
  ink.value.interactive.favoriteCnt = Math.max(0, ink.value.interactive.favoriteCnt - 1)
  ink.value.interactive.favorited = !ink.value.interactive.favorited
}

const ops = computed(() => {
  const baseOps = [
    {
      name: '举报',
      action: () => {
        console.log('举报')
      },
    },
  ]
  if (ink.value.author.id == userStore.getActiveUser()?.user.id) {
    baseOps.unshift({
      name: '编辑',
      action: () => {
        router.push(`/dashboard/${ink.value.author.account}/editor/${ink.value.id}`)
      },
    })
    baseOps.push({
      name: '删除',
      action: () => {
        console.log('删除')
      },
    })
  }
  return baseOps
})
</script>

<template>
  <div class="max-screen-w flex items-start line-padding">
    <div class="flex-1">
      <div class="text-xl flex items-center mb-4 semibold-text">
        <el-tooltip content="返回">
          <el-button text circle @click="router.back()"
            ><span class="material-symbols-outlined"> arrow_back </span>
          </el-button>
        </el-tooltip>
        <span class="ml-5 cursor-pointer">帖子</span>
      </div>
      <div v-loading="loading">
        <div class="flex justify-between">
          <div class="flex cursor-pointer">
            <InkPopover place="bottom">
              <template #reference>
                <UserAvatar :src="ink.author.avatar"></UserAvatar>
              </template>
              <template #content>
                <UserCard :user="ink.author"></UserCard>
              </template>
            </InkPopover>
            <router-link :to="`/user/${ink.author.account}`">
              <div class="ml-3">
                <div class="semibold-text">{{ ink.author.username }}</div>
                <el-link class="nav-text hover:underline">@{{ ink.author.account }}</el-link>
              </div>
            </router-link>
          </div>
          <MoreOperation :horizon="true" :operations="ops"></MoreOperation>
        </div>
        <!--      <MilkdownWrapper :padding="false" :read-only="true" ref="milkdownRef" class="mt-4">-->
        <!--      </MilkdownWrapper>-->
        <div
          class="prose mt-4 prose-slate dark:prose-invert"
          v-html="ink.contentHtml"
          ref="contentRef"
        ></div>
        <div class="nav-text mt-2">{{ formatDate(ink.createdAt) }}</div>
        <div class="mt-6 flex">
          <InkInteractive
            :interactive="ink.interactive"
            @like="handleLike"
            @cancel-like="handleCancelLike"
            @favorite="handleFavorite"
            @cancel-favorite="handleCancelFavorite"
          ></InkInteractive>
        </div>
      </div>
      <CommentView biz="ink" :biz-id="ink.id"></CommentView>
    </div>
    <div class="w-90 flex-col sticky-top line-padding ml-10">
      <div v-show="recommendInks.length > 0">
        <!--        <UserCard :user="ink.author" :auto-padding="false"></UserCard>-->
        <div class="mt-6">相关推荐</div>
        <div>
          <RecommendCard
            class="my-6"
            v-for="ink in recommendInks"
            :ink="ink"
            :key="ink.id"
          ></RecommendCard>
        </div>
      </div>
    </div>
    <BackTop></BackTop>
  </div>
</template>

<style scoped lang="scss"></style>
