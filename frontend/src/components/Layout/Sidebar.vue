<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <h3>管理后台</h3>
    </div>
    
    <nav class="sidebar-nav">
      <router-link to="/dashboard" class="nav-item">
        <span class="nav-icon">📊</span>
        <span class="nav-text">仪表板</span>
      </router-link>
      
      <router-link to="/dashboard/articles" class="nav-item">
        <span class="nav-icon">📝</span>
        <span class="nav-text">文章管理</span>
      </router-link>
      
      <router-link to="/dashboard/categories" class="nav-item">
        <span class="nav-icon">📁</span>
        <span class="nav-text">分类管理</span>
      </router-link>
      
      <router-link to="/dashboard/tags" class="nav-item">
        <span class="nav-icon">🏷️</span>
        <span class="nav-text">标签管理</span>
      </router-link>
      
      <router-link to="/dashboard/comments" class="nav-item">
        <span class="nav-icon">💬</span>
        <span class="nav-text">评论审核</span>
        <span v-if="pendingCount > 0" class="badge badge-danger">{{ pendingCount }}</span>
      </router-link>
      
      <router-link v-if="isAdmin" to="/dashboard/users" class="nav-item">
        <span class="nav-icon">👥</span>
        <span class="nav-text">用户管理</span>
      </router-link>
      
      <router-link v-if="isAdmin" to="/dashboard/audit-logs" class="nav-item">
        <span class="nav-icon">📋</span>
        <span class="nav-text">审计日志</span>
      </router-link>
    </nav>
    
    <div class="sidebar-footer">
      <router-link to="/" class="back-link">← 返回首页</router-link>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'

const userStore = useUserStore()
const pendingCount = ref(0)

const isAdmin = computed(() => userStore.isAdmin)

onMounted(() => {
  if (userStore.isLoggedIn) {
    loadPendingCount()
  }
})

async function loadPendingCount() {
  try {
    const res = await request.get('/comments/pending')
    pendingCount.value = res.data.comments?.length || 0
  } catch (e) {
    console.error(e)
  }
}
</script>

<style lang="scss" scoped>
.sidebar {
  width: var(--sidebar-width);
  background: var(--bg-card);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - var(--header-height));
  
  @media (max-width: 768px) {
    position: fixed;
    left: 0;
    top: var(--header-height);
    z-index: 99;
    transform: translateX(-100%);
    transition: transform var(--transition-normal);
    
    &.open {
      transform: translateX(0);
    }
  }
}

.sidebar-header {
  padding: var(--spacing-lg) var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  
  h3 {
    font-size: var(--font-lg);
    color: var(--text-primary);
  }
}

.sidebar-nav {
  flex: 1;
  padding: var(--spacing-md) 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm) var(--spacing-md);
  margin: 2px var(--spacing-sm);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-size: var(--font-sm);
  transition: all var(--transition-fast);
  
  &:hover {
    background: var(--bg-secondary);
    color: var(--text-primary);
  }
  
  &.router-link-active {
    background: var(--primary-color);
    color: var(--text-inverse);
  }
}

.nav-icon {
  font-size: var(--font-base);
  width: 20px;
  text-align: center;
}

.nav-text {
  flex: 1;
}

.badge {
  margin-left: auto;
}

.sidebar-footer {
  padding: var(--spacing-md);
  border-top: 1px solid var(--border-color);
}

.back-link {
  font-size: var(--font-sm);
  color: var(--text-tertiary);
  
  &:hover {
    color: var(--primary-color);
  }
}
</style>
