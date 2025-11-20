<script setup lang="ts">
import { reactive, ref } from 'vue'
import UserAvatar from '@/components/UserAvatar.vue'
import type { User } from '@/types/user.ts'
import type { ElForm, FormRules } from 'element-plus'
import { editProfile } from '@/service/user.ts'
import { notification } from '@/utils/notification.ts'
import { defaultBanner } from '@/consts/default.ts'
import ImageCropper, { type CropperResult } from '@/components/image/ImageCropper.vue'
import InkDialog from '@/components/InkDialog.vue'
import { uploadAvatar, uploadImage } from '@/service/file.ts'
import { PromiseGroup } from '@/utils/promise.ts'

const props = defineProps({
  user: {
    type: Object as () => User,
    required: true,
  },
})

const userInfoForm = ref({
  username: props.user.username,
  aboutMe: props.user.aboutMe,
  banner: props.user.banner,
  avatar: props.user.avatar,
  links: props.user.links,
})

const userInfoFormRef = ref<InstanceType<typeof ElForm>>()

const rules = reactive<FormRules>({
  username: {
    required: true,
    message: '昵称不能为空',
    trigger: 'blur',
  },
})

const showCropper = ref(false)
const tempImageSrc = ref('')
const tempImageName = ref('')
const cropperAspectRatio = ref<number | undefined>(1)
const cropperType = ref<'avatar' | 'banner'>('avatar')

const triggerCropper = (ratio?: number) => {
  const fileInput = document.createElement('input')
  fileInput.type = 'file'
  fileInput.accept = 'image/*'
  fileInput.onchange = (event) => {
    const file = (event.target as HTMLInputElement).files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        tempImageSrc.value = e.target?.result as string
        tempImageName.value = file.name
        cropperAspectRatio.value = ratio
        showCropper.value = true
      }
      reader.readAsDataURL(file)
    }
  }
  fileInput.click()
}
const handleAvatarClick = () => {
  cropperType.value = 'avatar'
  triggerCropper(1)
}

// Banner update
const handleBannerClick = () => {
  cropperType.value = 'banner'
  triggerCropper()
}

const base64ToFile = (src: string, name: string) => {
  const arr = src.split(',')
  const mime = arr[0].match(/:(.*?);/)?.[1]
  const bstr = atob(arr[1])
  let n = bstr.length
  const u8arr = new Uint8Array(n)
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n)
  }
  return new File([u8arr], name, { type: mime })
}

let bannerName = ''
let avatarName = ''
const handleCropConfirm = (result: CropperResult) => {
  // const file = base64ToFile(result.src, result.name)

  if (cropperType.value === 'banner') {
    userInfoForm.value.banner = result.canvas.toDataURL()
    bannerName = result.name
    showCropper.value = false
    return
  }

  userInfoForm.value.avatar = result.canvas.toDataURL()
  avatarName = result.name
  showCropper.value = false
}

const handleCropCancel = () => {
  showCropper.value = false
}

const loading = ref(false)
const handleSave = () => {
  userInfoFormRef.value
    ?.validate(async (valid) => {
      if (valid) {
        // 比对banner和avatar
        const pg = new PromiseGroup()
        loading.value = true
        if (userInfoForm.value.banner != props.user.banner) {
          pg.add(
            uploadImage(base64ToFile(userInfoForm.value.banner, bannerName)).then((url) => {
              userInfoForm.value.banner = url
            }),
          )
        }
        if (userInfoForm.value.avatar != props.user.avatar) {
          pg.add(
            uploadAvatar(base64ToFile(userInfoForm.value.avatar, avatarName)).then((url) => {
              userInfoForm.value.avatar = url
            }),
          )
        }

        await pg.wait()

        await editProfile(userInfoForm.value)
        loading.value = false
        notification({
          message: '保存成功',
        })
        location.reload()
      }
    })
    .catch(() => {
      loading.value = false
    })
}
</script>

<template>
  <div class="relative" v-loading="loading" element-loading-text="上传图片中...">
    <div class="relative">
      <el-image
        fit="cover"
        class="w-full h-40"
        :src="userInfoForm.banner == '' ? defaultBanner : userInfoForm.banner"
      ></el-image>
      <div
        title="切换照片"
        class="cursor-pointer absolute w-full h-40 left-0 top-0 flex justify-center items-center bg-black/20"
        @click="handleBannerClick"
      >
        <span class="material-symbols-outlined"> add_a_photo </span>
      </div>
    </div>
    <div class="absolute left-4 top-30">
      <div class="relative inline-block">
        <UserAvatar
          class="edit-avatar"
          :src="userInfoForm.avatar"
          :size="80"
          :border="true"
        ></UserAvatar>
        <div
          title="切换头像"
          class="absolute z-10 cursor-pointer absolute left-0 top-0 w-full h-full flex justify-center items-center bg-black/20 rounded-full"
          @click="handleAvatarClick"
        >
          <span class="material-symbols-outlined text-white"> add_a_photo </span>
        </div>
      </div>
    </div>
    <div class="mt-12 mb-4">
      <el-form :model="userInfoForm" label-position="top" :rules="rules" ref="userInfoFormRef">
        <el-form-item label="昵称" prop="username">
          <el-input size="large" v-model="userInfoForm.username"></el-input>
        </el-form-item>
        <el-form-item label="个人简介" prop="aboutMe">
          <el-input
            size="large"
            type="textarea"
            maxlength="160"
            show-word-limit
            v-model="userInfoForm.aboutMe"
            rows="3"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="链接" prop="links">
          <el-input-tag
            size="large"
            v-model="userInfoForm.links"
            draggable
            placeholder="输入后按回车添加"
          />
        </el-form-item>
      </el-form>
    </div>
    <div class="mb-2 flex justify-end">
      <el-button class="w-24" round size="large" type="primary" @click="handleSave">保存</el-button>
    </div>

    <!-- Image Cropper Dialog -->
    <InkDialog v-model="showCropper" title="裁剪图片">
      <ImageCropper
        :src="tempImageSrc"
        :image-name="tempImageName"
        :aspect-ratio="cropperAspectRatio"
        @confirm="handleCropConfirm"
        @cancel="handleCropCancel"
      />
    </InkDialog>
  </div>
</template>

<style scoped lang="scss">
.edit-avatar {
  &:hover {
    .avatar-overlay {
      opacity: 1;
    }
  }
}
</style>
