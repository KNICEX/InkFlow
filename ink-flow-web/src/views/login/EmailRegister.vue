<template>
  <div>
    <div v-if="step == 'pwd'">
      <el-form
        v-loading="loading"
        ref="registerPwdFormRef"
        :model="registerPwdInfo"
        :rules="registerPwdRules"
      >
        <el-form-item label="Password" prop="password" class="h-26">
          <FormInput v-model="registerPwdInfo.password" type="password"></FormInput>
        </el-form-item>
        <el-form-item label="Confirm Password" prop="re" class="h-26">
          <FormInput v-model="registerPwdInfo.re" type="password"></FormInput>
        </el-form-item>
        <div class="flex justify-between items-center mt-3">
          <el-link class="flex-1" :underline="false" @click="handleHasAccount"> 已有账号 </el-link>
          <el-button
            size="large"
            @click="toAccountStep"
            style="height: 3rem; width: 12rem"
            class="ellipse-button"
            bg
            text
            >Register
          </el-button>
        </div>
      </el-form>
    </div>
    <div v-else>
      <el-form
        v-loading="loading"
        ref="registerAccountFormRef"
        :model="registerAccountInfo"
        :rules="registerAccountRules"
      >
        <el-form-item ref="emailFormItemRef" label="Account" prop="account" class="h-26 relative">
          <FormInput v-model="registerAccountInfo.account"> </FormInput>
        </el-form-item>
        <el-form-item label="Username" prop="username" class="h-26">
          <FormInput :maxlength="30" v-model="registerAccountInfo.username"></FormInput>
        </el-form-item>

        <div class="flex justify-between items-center mt-3">
          <el-link class="flex-1" :underline="false" @click="handleHasAccount"> 已有账号 </el-link>
          <el-button
            size="large"
            @click="handleRegister"
            style="height: 3rem; width: 12rem"
            class="ellipse-button"
            bg
            text
            >Mint
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, useTemplateRef, reactive } from 'vue'
import { ElForm, ElFormItem, type FormRules } from 'element-plus'
import FormInput from '@/components/form/FormInput.vue'
import { useRegisterTokenStore, useUserStore } from '@/stores/user.ts'
import { notification } from '@/utils/notification.ts'
import { profile, registerByEmail } from '@/service/user.ts'

const registerPwdFormRef = useTemplateRef<InstanceType<typeof ElForm>>('registerPwdFormRef')
const registerAccountFormRef = useTemplateRef<InstanceType<typeof ElForm>>('registerAccountFormRef')
const emit = defineEmits(['reportSize', 'reportTitle', 'hasAccount', 'registerSuccess'])
onMounted(() => {
  emit('reportSize', '24rem', '26rem')
  emit('reportTitle', 'Register 2/2')
})
const handleHasAccount = () => {
  emit('hasAccount')
}

const step = ref('pwd')
const registerPwdInfo = ref({
  password: '',
  re: '',
})

const registerPwdRules = reactive<FormRules>({
  password: [
    { required: true, message: 'Please input password' },
    { min: 6, message: 'Password length should be at least 6' },
  ],
  re: {
    validator: (rule, value, callback) => {
      if (value === '') {
        callback(new Error('Please input password again'))
      } else if (value !== registerPwdInfo.value.password) {
        callback(new Error('Two passwords are inconsistent'))
      } else {
        callback()
      }
    },
  },
})
const loading = ref(false)
const toAccountStep = () => {
  registerPwdFormRef.value?.validate(async (valid: boolean) => {
    if (valid) {
      step.value = 'account'
      emit('reportTitle', 'Mint Account')
    }
  })
}

const registerAccountInfo = ref({
  account: '',
  username: '',
})

const registerAccountRules = {
  account: [{ required: true, message: 'Please input account' }],
  username: [{ required: true, message: 'Please input username' }],
}

const registerTokenStore = useRegisterTokenStore()
const handleRegister = () => {
  registerAccountFormRef.value?.validate(async (valid: boolean) => {
    if (valid) {
      const token = registerTokenStore.getEmailToken()
      console.log('token', token)
      if (!token.token) {
        // todo show message
        notification({
          type: 'error',
          title: 'Error',
          message: '邮箱验证码已过期，请重新获取验证码',
        })
        return
      }
      loading.value = true
      await registerByEmail({
        account: registerAccountInfo.value.account,
        username: registerAccountInfo.value.username,
        password: registerPwdInfo.value.password,
        token: token.token,
        email: token.email,
      }).catch(() => {
        loading.value = false
      })
      loading.value = false
      const userStore = useUserStore()
      const user = await profile()
      userStore.updateActiveUser(user)
      notification({
        message: '注册成功',
      })
      emit('registerSuccess')
    }
  })
}

const back = () => {
  if (step.value === 'account') {
    step.value = 'pwd'
    emit('reportTitle', 'Register 2/2')
    return false
  } else {
    return true
  }
}
defineExpose({
  back,
})
</script>
<style scoped lang="scss"></style>
