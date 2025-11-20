<template>
  <div>
    <el-form ref="pwdLoginFormRef" label-position="top" :model="pwdLoginInfo" :rules="pwdRules">
      <el-form-item label="Email" prop="email" class="h-26">
        <FormInput v-model="pwdLoginInfo.email"></FormInput>
      </el-form-item>
      <el-form-item label="Password" prop="password" class="h-26">
        <FormInput type="password" v-model="pwdLoginInfo.password"></FormInput>
      </el-form-item>
    </el-form>
    <div class="flex justify-between items-center mt-3">
      <el-link class="flex-1" :underline="false" @click="emit('noAccount')">
        没有账号或忘记密码？
      </el-link>
      <el-button
        size="large"
        @click="handlePwdLogin"
        style="height: 3rem; width: 8rem"
        class="ellipse-button"
        bg
        text
        >Login
      </el-button>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref, useTemplateRef } from 'vue'
import FormInput from '@/components/form/FormInput.vue'
import { ElForm, type FormRules } from 'element-plus'
import { loginEmailPwd } from '@/service/user.ts'
import { notifySuccessLogin } from '@/utils/notification.ts'

const pwdLoginInfo = ref({
  email: '',
  password: '',
})

const emit = defineEmits(['reportSize', 'reportTitle', 'noAccount'])
onMounted(() => {
  emit('reportSize', '24rem', '24rem')
  emit('reportTitle', '邮箱密码登录')
})

type ElFormType = InstanceType<typeof ElForm>
const pwdLoginFormRef = useTemplateRef<ElFormType>('pwdLoginFormRef')

const pwdRules = reactive<FormRules>({
  email: [
    { required: true, message: 'Please input email' },
    { type: 'email', message: 'Please input correct email' },
  ],
  password: [
    { required: true, message: 'Please input password' },
    { min: 6, message: 'Password length should be at least 6' },
  ],
})

const handlePwdLogin = () => {
  console.log('tes')
  pwdLoginFormRef?.value?.validate(async (valid: boolean) => {
    if (valid) {
      await loginEmailPwd({
        email: pwdLoginInfo.value.email,
        password: pwdLoginInfo.value.password,
      })
      notifySuccessLogin()
      window.location.reload()
    } else {
      // return
    }
  })
}
</script>
<style scoped lang="scss"></style>
