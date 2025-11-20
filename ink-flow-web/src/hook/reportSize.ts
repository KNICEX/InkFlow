import { onMounted } from 'vue'

export const reportSize = (width:string, height:string) => {
  const emit = defineEmits(['reportSize'])
  onMounted(() => {
    emit('reportSize', { width, height })
  })
}
