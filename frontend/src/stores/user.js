import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import request from '@/utils/request'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || 'null'))
  const refreshToken = ref(localStorage.getItem('refreshToken') || '')

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')

  function setToken(newToken, newRefreshToken) {
    token.value = newToken
    refreshToken.value = newRefreshToken || refreshToken.value
    localStorage.setItem('token', newToken)
    if (newRefreshToken) {
      localStorage.setItem('refreshToken', newRefreshToken)
    }
  }

  function setUserInfo(info) {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  async function login(credentials) {
    const res = await request.post('/auth/login', credentials)
    setToken(res.data.token, res.data.refreshToken)
    setUserInfo(res.data.user)
    return res.data
  }

  async function register(data) {
    const res = await request.post('/auth/register', data)
    return res.data
  }

  async function fetchUserInfo() {
    const res = await request.get('/user/info')
    setUserInfo(res.data)
    return res.data
  }

  function logout() {
    token.value = ''
    refreshToken.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('userInfo')
  }

  async function refreshAccessToken() {
    if (!refreshToken.value) return false
    try {
      const res = await request.post('/auth/refresh', { refreshToken: refreshToken.value })
      setToken(res.data.token)
      return true
    } catch {
      logout()
      return false
    }
  }

  return {
    token,
    userInfo,
    refreshToken,
    isLoggedIn,
    isAdmin,
    login,
    register,
    fetchUserInfo,
    logout,
    refreshAccessToken,
    setToken
  }
})
