<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-header">
        <h1>登录</h1>
        <p>欢迎回来</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="auth-form">
        <div class="form-group">
          <label>用户名或邮箱</label>
          <input 
            v-model="form.username" 
            type="text" 
            placeholder="请输入用户名或邮箱"
            required
          />
        </div>
        
        <div class="form-group">
          <label>密码</label>
          <input 
            v-model="form.password" 
            type="password" 
            placeholder="请输入密码"
            required
          />
        </div>
        
        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
      
      <div class="auth-divider">
        <span>或</span>
      </div>
      
      <button class="btn btn-secondary btn-block github-btn" @click="loginWithGitHub">
        <span class="icon">🐙</span>
        GitHub 登录
      </button>
      
      <div class="auth-footer">
        还没有账号？
        <router-link to="/register">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: ''
})
const loading = ref(false)

async function handleLogin() {
  loading.value = true
  try {
    const data = await userStore.login(form.value)
    
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (e) {
    alert(e.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}

function loginWithGitHub() {
  alert('GitHub OAuth 需要配置环境变量')
}
</script>

<style lang="scss" scoped>
.auth-page {
  min-height: calc(100vh - var(--header-height));
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
  padding: var(--spacing-lg);
}

.auth-card {
  width: 100%;
  max-width: 400px;
  background: var(--bg-card);
  border-radius: var(--radius-xl);
  padding: var(--spacing-2xl);
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--border-color);
}

.auth-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  
  h1 {
    font-size: var(--font-2xl);
    margin-bottom: var(--spacing-xs);
  }
  
  p {
    color: var(--text-secondary);
    font-size: var(--font-sm);
  }
}

.auth-form {
  margin-bottom: var(--spacing-lg);
}

.form-group {
  margin-bottom: var(--spacing-md);
  
  label {
    display: block;
    margin-bottom: var(--spacing-xs);
    font-size: var(--font-sm);
    color: var(--text-secondary);
    font-weight: 500;
  }
  
  input {
    width: 100%;
    padding: var(--spacing-sm) var(--spacing-md);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    font-size: var(--font-sm);
    background: var(--bg-primary);
    color: var(--text-primary);
    transition: all var(--transition-fast);
    
    &:focus {
      outline: none;
      border-color: var(--primary-color);
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }
  }
}

.auth-divider {
  position: relative;
  text-align: center;
  margin: var(--spacing-lg) 0;
  
  span {
    background: var(--bg-card);
    padding: 0 var(--spacing-md);
    color: var(--text-tertiary);
    font-size: var(--font-xs);
    position: relative;
    z-index: 1;
  }
  
  &::before {
    content: '';
    position: absolute;
    left: 0;
    right: 0;
    top: 50%;
    height: 1px;
    background: var(--border-color);
  }
}

.github-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  
  .icon {
    font-size: var(--font-lg);
  }
}

.auth-footer {
  text-align: center;
  margin-top: var(--spacing-lg);
  font-size: var(--font-sm);
  color: var(--text-secondary);
  
  a {
    color: var(--primary-color);
    font-weight: 500;
  }
}
</style>
