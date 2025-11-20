<template>
  <el-form
    status-icon
    v-loading="loading"
    ref="codeLoginFormRef"
    :model="codeLoginInfo"
    :rules="codeRules"
  >
    <el-form-item ref="emailFormItemRef" label="Email" prop="email" class="h-26 relative">
      <FormInput v-model="codeLoginInfo.email"> </FormInput>
      <div @click="sendEmailCode" :class="verifyClass">
        {{ verifyPrompt }}
      </div>
    </el-form-item>
    <el-form-item label="Code" prop="code" class="h-26">
      <FormInput :maxlength="6" v-model="codeLoginInfo.code"></FormInput>
    </el-form-item>
    <div class="flex justify-between items-center mt-3">
      <el-link class="flex-1" :underline="false" @click="handleHasAccount"> 已有账号 </el-link>
      <el-button
        size="large"
        @click="handleCodeLogin"
        style="height: 3rem; width: 12rem"
        class="ellipse-button"
        bg
        text
        >Login or Register
      </el-button>
    </div>
  </el-form>
</template>
<script setup lang="ts">
import { ElForm, ElFormItem, type FormRules } from 'element-plus'
import FormInput from '@/components/form/FormInput.vue'
import { loginEmailCode, sendLoginCode } from '@/service/user.ts'
import { computed, onMounted, reactive, ref, useTemplateRef } from 'vue'
import { useRegisterTokenStore } from '@/stores/user.ts'
import clsx from 'clsx'
import { notification } from '@/utils/notification.ts'
const registerTokenStore = useRegisterTokenStore()
const codeLoginFormRef = useTemplateRef<InstanceType<typeof ElForm>>('codeLoginFormRef')
const emailFormItemRef = useTemplateRef<InstanceType<typeof ElFormItem>>('emailFormItemRef')
const emit = defineEmits(['reportSize', 'reportTitle', 'hasAccount', 'nextStep'])
onMounted(() => {
  emit('reportSize', '24rem', '24rem')
  emit('reportTitle', '注册 1/2')
})

const handleHasAccount = () => {
  emit('hasAccount')
}

const codeSent = ref(false)
const loading = ref(false)

const verifyClass = computed(() => {
  return emailFormItemRef.value?.validateState == 'success' && !codeSent.value
    ? clsx(
        'absolute text-white hover:bg-green-500 flex justify-center cursor-pointer items-center verify bg-green-400',
      )
    : clsx(
        'absolute text-white flex justify-center cursor-not-allowed items-center verify bg-green-300',
      )
})
const verifyPrompt = ref('Verify')

const codeLoginInfo = ref({
  email: '',
  code: '',
})
const codeRules = reactive<FormRules>({
  email: [
    { required: true, message: 'Please input email' },
    { type: 'email', message: 'Please input correct email' },
  ],
  code: [
    { required: true, message: 'Please input code' },
    { min: 6, message: 'Code length should be at least 6' },
  ],
})

const sendEmailCode = async () => {
  if (emailFormItemRef.value?.validateState !== 'success' || codeSent.value) {
    return
  }
  await sendLoginCode(codeLoginInfo.value.email)
  notification({
    message: '验证码发送成功',
  })
  let count = 60
  codeSent.value = true
  verifyPrompt.value = `Try Again ${count}`
  const timer = setInterval(() => {
    count--
    verifyPrompt.value = `Try Again ${count}`
    if (count === 0) {
      codeSent.value = false
      clearInterval(timer)
      verifyPrompt.value = 'Verify'
    }
  }, 1000)
}
const handleCodeLogin = async () => {
  codeLoginFormRef?.value?.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true

      const token = await loginEmailCode({
        email: codeLoginInfo.value.email,
        code: codeLoginInfo.value.code,
      }).catch(() => {
        loading.value = false
      })
      loading.value = false
      if (!token) {
        // 登录成功
        return
      }
      // 保存临时token，进入注册流程
      registerTokenStore.setEmailToken(codeLoginInfo.value.email, token)
      notification({
        message: '邮箱验证成功，请注册账号',
      })
      emit('nextStep')
    }
  })
}
</script>
<style scoped lang="scss">
.verify {
  width: 6rem;
  right: 0;
  top: 1.3rem;
  height: 3rem;
  border-radius: 0 1rem 1rem 0;
}
</style>
