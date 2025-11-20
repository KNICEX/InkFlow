import axios, { type AxiosRequestConfig } from 'axios'
import { AxiosError } from 'axios'
import { useUserStore } from '@/stores/user.ts'
import { notification } from '@/utils/notification.ts'
import { emptyUser } from '@/types/user.ts'
import { refreshToken } from '@/service/user.ts'
import router from '@/router/index.ts'
import { parseResponseDate } from '@/utils/parse.ts'

interface ApiResponse<T> {
  code: number
  msg: string
  data: T
}

interface CustomRequestConfig<T> extends AxiosRequestConfig {
  onRespError?: (error: AxiosError | Error) => void
  onReqError?: (error: Error) => void
  onSuccess?: (data: T) => void
}

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 5000,
})

request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    const activeUser = userStore.getActiveUser()
    if (activeUser) {
      config.headers['Authorization'] = `Bearer ${activeUser.token}`
    }
    return config
  },
  (error) => {
    console.error('Request Error: ', error)
    error.config?.onReqError?.(error)
    return Promise.reject(error)
  },
)

let refreshedToken = false
request.interceptors.response.use(
  (response) => {
    // 处理服务端刷新token
    const newToken = response.headers['x-access-token']
    const newRefreshToken = response.headers['x-refresh-token']
    if (newToken || newRefreshToken) {
      const userStore = useUserStore()
      const activeUser = userStore.getActiveUser()
      // TODO 还没测试，这里可能拿到空user
      if (activeUser) {
        let tokenUpdated = false
        if (newToken) {
          tokenUpdated = true
          activeUser.token = newToken
        }
        if (newRefreshToken) {
          tokenUpdated = true
          activeUser.refreshToken = newRefreshToken
        }
        if (tokenUpdated) {
          userStore.setActiveUser(activeUser)
        }
      } else {
        // 登陆完毕
        userStore.setActiveUser({
          user: emptyUser(),
          token: newToken,
          refreshToken: newRefreshToken,
        })
      }
    }

    const { data } = response
    if (data.code != 0) {
      // TODO 业务错误码，暂定
      notification({
        title: 'Error',
        message: data.msg,
        type: 'error',
      })
      console.error('Response Error: ', data)

      return Promise.reject(response)
    }

    // 处理时间类型
    parseResponseDate(data.data)

    return response
  },
  async (error) => {
    console.error('Response Error: ', error)
    if (error.response && error.response.status === 401) {
      // 1. 临时token过期，这时要调用refreshToken 2. refreshToken过期，这时要重新登录
      if (refreshedToken) {
        // refreshToken也过期，删除当前登录账户的信息
        const userStore = useUserStore()
        userStore.setActiveUser(null)
        refreshedToken = false
        return
      }
      refreshedToken = true
      await refreshToken()
    }
    return Promise.reject(error)
  },
)

const catchError = <T>(error: unknown, config?: CustomRequestConfig<T>) => {
  if (error instanceof AxiosError) {
    config?.onReqError?.(error)
  } else if (error instanceof Error) {
    config?.onReqError?.(error)
  }
  // TODO 可能会有更多错误类型
}

async function get<T>(url: string, config?: CustomRequestConfig<T>): Promise<T> {
  try {
    const response = await request.get<ApiResponse<T>>(url, config)
    const respData = response.data
    config?.onSuccess?.(respData?.data)
    return respData.data
  } catch (error) {
    catchError<T>(error)
    throw error
  }
}

async function post<T, D>(url: string, data: D, config?: CustomRequestConfig<T>): Promise<T> {
  try {
    const response = await request.post<ApiResponse<T>>(url, data, config)
    const respData = response.data
    config?.onSuccess?.(respData?.data)
    return respData.data
  } catch (error) {
    catchError<T>(error)
    throw error
  }
}

async function put<T, D>(url: string, data: D, config?: CustomRequestConfig<T>): Promise<T> {
  try {
    const response = await request.put<ApiResponse<T>>(url, data, config)
    const respData = response.data
    config?.onSuccess?.(respData?.data)
    return respData.data
  } catch (error) {
    catchError<T>(error)
    throw error
  }
}

async function del<T>(url: string, config?: CustomRequestConfig<T>): Promise<T> {
  try {
    const response = await request.delete<ApiResponse<T>>(url, config)
    const respData = response.data
    config?.onSuccess?.(respData?.data)
    return respData.data
  } catch (error) {
    catchError<T>(error)
    throw error
  }
}

export { get, post, put, del }
