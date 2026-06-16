import axios from 'axios'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

const request = axios.create({
  baseURL: '/api',
  timeout: 30000
})

let csrfToken = null

async function fetchCSRFToken() {
  try {
    const res = await axios.get('/api/csrf-token')
    csrfToken = res.data.csrf_token
    return csrfToken
  } catch (e) {
    return null
  }
}

request.interceptors.request.use(
  async (config) => {
    const userStore = useUserStore()
    
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    
    if (config.method !== 'get' && config.method !== 'head' && config.method !== 'options') {
      if (!csrfToken) {
        await fetchCSRFToken()
      }
      if (csrfToken) {
        config.headers['X-CSRF-Token'] = csrfToken
      }
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const userStore = useUserStore()
    const router = useRouter()
    
    if (error.response?.status === 401) {
      if (userStore.refreshToken) {
        try {
          const success = await userStore.refreshAccessToken()
          if (success) {
            error.config.headers.Authorization = `Bearer ${userStore.token}`
            return request(error.config)
          }
        } catch {
        }
      }
      userStore.logout()
      router.push({ name: 'Login' })
    }
    
    if (error.response?.status === 403 && error.response.data?.error?.includes('CSRF')) {
      csrfToken = null
      return request(error.config)
    }
    
    return Promise.reject(error)
  }
)

export default request
