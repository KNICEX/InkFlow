<template>
  <InkDialog
    :width="dialogSize.width"
    :height="dialogSize.height"
    v-model="model"
    :show-back="showBack"
    @on-back="handleBack"
    :title="title"
    @on-close="handleClose"
  >
    <div class="mt-6 px-2">
      <div v-if="loginType == ''" class="flex flex-col items-center">
        <div :class="item" @click="changeType(LoginType.Email)">
          <div>使用邮箱登录</div>
          <div>图标</div>
        </div>
        <div :class="item" @click="changeType(LoginType.Github)">
          <div>使用Github登录</div>
          <div>图标</div>
        </div>
        <div :class="item" @click="changeType(LoginType.Google)">
          <div>使用Google账户登录</div>
          <div>图标</div>
        </div>
      </div>
      <EmailPwd
        v-else-if="loginType == LoginType.Email"
        @report-size="handleReportSize"
        @no-account="changeType(LoginType.EmailCode)"
        @report-title="handleReportTitle"
      ></EmailPwd>
      <EmailCode
        v-else-if="loginType == LoginType.EmailCode"
        @report-size="handleReportSize"
        @report-title="handleReportTitle"
        @next-step="changeType(LoginType.EmailRegister)"
        @has-account="changeType(LoginType.Email)"
      ></EmailCode>
      <EmailRegister
        ref="emailRegisterRef"
        v-else-if="loginType == LoginType.EmailRegister"
        @report-size="handleReportSize"
        @has-account="changeType(LoginType.Email)"
        @report-title="handleReportTitle"
        @register-success="handleRegisterSuccess"
      ></EmailRegister>
    </div>
  </InkDialog>
</template>
<script setup lang="ts">
import { computed, reactive, ref, useTemplateRef } from 'vue'
import clsx from 'clsx'
import InkDialog from '@/components/InkDialog.vue'
import EmailPwd from '@/views/login/EmailPwd.vue'
import EmailCode from '@/views/login/EmailCode.vue'
import EmailRegister from '@/views/login/EmailRegister.vue'

const item = clsx(
  'flex cursor-pointer mt-4  text-base text-gray-700 font-semibold rounded-xl bg-gray-100 hover:bg-gray-200 py-6 px-4 justify-between w-full',
)
const model = defineModel()

const emailRegisterRef = useTemplateRef<InstanceType<typeof EmailRegister>>('emailRegisterRef')
const loginType = ref('')
const dialogSize = reactive({
  width: '20rem',
  height: '28rem',
})

enum LoginType {
  Home = '',
  Email = 'email',
  EmailCode = 'email-code',
  Github = 'github',
  Google = 'google',
  EmailRegister = 'email-register',
}

const title = ref('请选择登录方式')

const showBack = computed(() => loginType.value !== '')
const changeType = (type: LoginType) => {
  if (type == LoginType.Home) {
    title.value = '请选择登录方式'
    dialogSize.width = '20rem'
    dialogSize.height = '28rem'
  }
  console.log('changeType', type)
  loginType.value = type
}

const handleBack = () => {
  switch (loginType.value) {
    case LoginType.Email:
      changeType(LoginType.Home)
      break
    case LoginType.EmailCode:
      changeType(LoginType.Email)
      break
    case LoginType.EmailRegister:
      if (emailRegisterRef.value?.back()) {
        changeType(LoginType.EmailCode)
      }
      break
  }
}
const handleReportSize = (width: string, height: string) => {
  dialogSize.width = width
  dialogSize.height = height
}
const handleReportTitle = (t: string) => {
  title.value = t
}

const handleClose = () => {
  // 关闭登录框之后，重置登录类型
  changeType(LoginType.Home)
}

const handleRegisterSuccess = () => {
  model.value = false
}
</script>
<style scoped lang="scss"></style>
