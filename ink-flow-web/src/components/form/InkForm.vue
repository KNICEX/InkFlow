<template>
    <form>

    </form>
</template>
<script setup lang="ts" generic="T">
import type { PropType } from 'vue'
import { defineProps, ref } from 'vue'
import type {  FormRules } from 'element-plus'

const props = defineProps({
  formItems: {
    type: Array as PropType<FormItem[]>,
    default: () => []
  },
  rules:{
    type : Object as PropType<FormRules>,
    default: () => ({})
  }
})

interface FormItem {
  label: string
  prop: string // prop是T的key 并且是string类型
  type: string
  default?: string
}
const formData = ref({} as T)
for (const item of props.formItems) {
  formData.value[item.prop] = item.default ?? ''
}



</script>
<style scoped lang="scss">

</style>
