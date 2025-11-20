import { post } from '@/service/axios.ts'

export const uploadCover = async (file: File) => {
  const form = new FormData()
  form.append('cover', file)
  return await post<string, FormData>('/file/cover', form)
}

export const uploadAvatar = async (file: File) => {
  const form = new FormData()
  form.append('avatar', file)
  return await post<string, FormData>('/file/avatar', form)
}

export const uploadImage = async (file: File) => {
  const form = new FormData()
  form.append('image', file)
  return await post<string, FormData>('/file/image', form)
}
