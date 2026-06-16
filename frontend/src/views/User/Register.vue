<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-header">
        <h1>注册</h1>
        <p>创建一个新账号</p>
      </div>
      
      <form @submit.prevent="handleRegister" class="auth-form">
        <div class="form-group">
          <label>用户名</label>
          <input 
            v-model="form.username" 
            type="text" 
            placeholder="请输入用户名"
            required
          />
        </div>
        
        <div class="form-group">
          <label>邮箱</label>
          <input 
            v-model="form.email" 
            type="email" 
            placeholder="请输入邮箱"
            required
          />
        </div>
        
        <div class="form-group">
          <label>密码</label>
          <input 
            v-model="form.password" 
            type="password" 
            placeholder="至少6位，包含字母和数字"
            required
          />
          <p class="hint">密码至少6位，必须包含字母和数字</p>
        </div>
        
        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>
      
      <div class="auth-footer">
        已有账号？
        <router-link to="/login">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  email: '',
  password: ''
})
const loading = ref(false)

async function handleRegister() {
  loading.value = true
  try {
    await userStore.register(form.value)
    alert('注册成功，请登录')
    router.push('/login')
  } catch (e) {
    alert(e.response?.data?.error || '注册失败')
  } finally {
    loading.value = false
  }
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
  
  .hint {
    margin-top: var(--spacing-xs);
    font-size: var(--font-xs);
    color: var(--text-tertiary);
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
