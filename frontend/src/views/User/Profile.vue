<template>
  <div class="profile-page">
    <div class="container">
      <div class="profile-card">
        <div class="profile-header">
          <div class="avatar">
            {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}
          </div>
          <div class="profile-info">
            <h2>{{ userStore.userInfo?.username }}</h2>
            <p class="email">{{ userStore.userInfo?.email }}</p>
            <span class="role-badge" :class="userStore.userInfo?.role">
              {{ getRoleLabel(userStore.userInfo?.role) }}
            </span>
          </div>
        </div>
        
        <div class="profile-stats">
          <div class="stat-item">
            <span class="stat-value">{{ articleCount }}</span>
            <span class="stat-label">文章</span>
          </div>
        </div>
        
        <div class="profile-actions">
          <router-link to="/article/editor" class="btn btn-primary">
            ✏️ 写文章
          </router-link>
          <router-link v-if="userStore.isAdmin" to="/dashboard" class="btn btn-secondary">
            📊 管理后台
          </router-link>
          <button class="btn btn-danger" @click="logout">
            🚪 退出登录
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'

const router = useRouter()
const userStore = useUserStore()
const articleCount = ref(0)

function getRoleLabel(role) {
  const roles = {
    admin: '管理员',
    editor: '编辑',
    author: '作者',
    guest: '游客'
  }
  return roles[role] || role
}

async function loadUserArticles() {
  try {
    const res = await request.get('/articles', {
      params: { author_id: userStore.userInfo?.id, status: 'all', page_size: 100 }
    })
    articleCount.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  }
}

function logout() {
  if (confirm('确定要退出登录吗？')) {
    userStore.logout()
    router.push('/')
  }
}

onMounted(() => {
  loadUserArticles()
})
</script>

<style lang="scss" scoped>
.profile-page {
  padding: var(--spacing-2xl) 0;
  min-height: calc(100vh - var(--header-height) - 200px);
}

.profile-card {
  max-width: 600px;
  margin: 0 auto;
  background: var(--bg-card);
  border-radius: var(--radius-xl);
  padding: var(--spacing-2xl);
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
}

.profile-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-xl);
  
  @media (max-width: 640px) {
    flex-direction: column;
    text-align: center;
  }
}

.avatar {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--primary-color);
  color: white;
  font-size: var(--font-2xl);
  font-weight: 600;
}

.profile-info {
  flex: 1;
  
  h2 {
    font-size: var(--font-xl);
    margin-bottom: var(--spacing-xs);
  }
  
  .email {
    color: var(--text-secondary);
    font-size: var(--font-sm);
    margin-bottom: var(--spacing-sm);
  }
}

.role-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-size: var(--font-xs);
  font-weight: 500;
  
  &.admin {
    background: #fef3c7;
    color: #92400e;
  }
  
  &.editor {
    background: #dbeafe;
    color: #1e40af;
  }
  
  &.author {
    background: #d1fae5;
    color: #065f46;
  }
}

.theme-dark .role-badge {
  &.admin {
    background: rgba(251, 191, 36, 0.2);
    color: #fbbf24;
  }
  
  &.editor {
    background: rgba(59, 130, 246, 0.2);
    color: #60a5fa;
  }
  
  &.author {
    background: rgba(16, 185, 129, 0.2);
    color: #34d399;
  }
}

.profile-stats {
  display: flex;
  gap: var(--spacing-xl);
  padding: var(--spacing-lg) 0;
  border-top: 1px solid var(--border-color);
  border-bottom: 1px solid var(--border-color);
  margin-bottom: var(--spacing-lg);
}

.stat-item {
  text-align: center;
  
  .stat-value {
    display: block;
    font-size: var(--font-2xl);
    font-weight: 700;
    color: var(--text-primary);
  }
  
  .stat-label {
    font-size: var(--font-sm);
    color: var(--text-secondary);
  }
}

.profile-actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}
</style>
