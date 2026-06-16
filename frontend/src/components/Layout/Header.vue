<template>
  <header class="app-header">
    <div class="container header-inner">
      <div class="logo-area">
        <router-link to="/" class="logo">
          <span class="logo-icon">📝</span>
          <span class="logo-text">轻量 CMS</span>
        </router-link>
        
        <nav class="main-nav">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/categories" class="nav-link">分类</router-link>
          <router-link to="/tags" class="nav-link">标签</router-link>
          <router-link to="/search" class="nav-link">搜索</router-link>
        </nav>
      </div>
      
      <div class="header-actions">
        <div class="search-box">
          <input 
            v-model="searchKeyword" 
            type="text" 
            placeholder="搜索文章..." 
            @keyup.enter="goSearch"
            class="search-input"
          />
          <button class="search-btn" @click="goSearch">🔍</button>
        </div>
        
        <button class="theme-toggle" @click="toggleTheme" :title="themeStore.isDark ? '切换亮色' : '切换暗色'">
          {{ themeStore.isDark ? '☀️' : '🌙' }}
        </button>
        
        <template v-if="userStore.isLoggedIn">
          <router-link to="/article/editor" class="btn btn-primary btn-sm write-btn">
            ✏️ 写文章
          </router-link>
          
          <div class="user-menu">
            <button class="user-avatar" @click="showMenu = !showMenu">
              {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}
            </button>
            
            <div v-if="showMenu" class="dropdown-menu" @click.self="showMenu = false">
              <router-link to="/profile" class="dropdown-item" @click="showMenu = false">
                👤 个人中心
              </router-link>
              <router-link v-if="userStore.isAdmin" to="/dashboard" class="dropdown-item" @click="showMenu = false">
                📊 管理后台
              </router-link>
              <div class="dropdown-divider"></div>
              <button class="dropdown-item danger" @click="handleLogout">
                🚪 退出登录
              </button>
            </div>
          </div>
        </template>
        
        <template v-else>
          <router-link to="/login" class="btn btn-secondary btn-sm">登录</router-link>
          <router-link to="/register" class="btn btn-primary btn-sm">注册</router-link>
        </template>
        
        <button class="mobile-menu-btn" @click="showMobileMenu = !showMobileMenu">
          {{ showMobileMenu ? '✕' : '☰' }}
        </button>
      </div>
    </div>
    
    <div v-if="showMobileMenu" class="mobile-menu">
      <div class="mobile-nav">
        <router-link to="/" class="mobile-nav-link" @click="showMobileMenu = false">首页</router-link>
        <router-link to="/categories" class="mobile-nav-link" @click="showMobileMenu = false">分类</router-link>
        <router-link to="/tags" class="mobile-nav-link" @click="showMobileMenu = false">标签</router-link>
        <router-link to="/search" class="mobile-nav-link" @click="showMobileMenu = false">搜索</router-link>
        <router-link v-if="userStore.isLoggedIn" to="/article/editor" class="mobile-nav-link" @click="showMobileMenu = false">写文章</router-link>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useThemeStore } from '@/stores/theme'
import { useUserStore } from '@/stores/user'

const themeStore = useThemeStore()
const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const searchKeyword = ref('')
const showMenu = ref(false)
const showMobileMenu = ref(false)

function toggleTheme() {
  themeStore.toggleTheme()
}

function goSearch() {
  if (searchKeyword.value.trim()) {
    router.push({ name: 'SearchResults', query: { q: searchKeyword.value } })
  }
}

function handleLogout() {
  userStore.logout()
  showMenu.value = false
  router.push({ name: 'Home' })
}
</script>

<style lang="scss" scoped>
.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
}

.header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--header-height);
}

.logo-area {
  display: flex;
  align-items: center;
  gap: var(--spacing-xl);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: var(--font-xl);
  font-weight: 700;
  color: var(--text-primary);
  
  &:hover {
    color: var(--primary-color);
  }
}

.logo-icon {
  font-size: var(--font-2xl);
}

.main-nav {
  display: flex;
  gap: var(--spacing-lg);
  
  @media (max-width: 768px) {
    display: none;
  }
}

.nav-link {
  color: var(--text-secondary);
  font-size: var(--font-sm);
  font-weight: 500;
  transition: color var(--transition-fast);
  
  &:hover, &.router-link-active {
    color: var(--primary-color);
  }
}

.header-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.search-box {
  display: flex;
  align-items: center;
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  padding: 4px 4px 4px var(--spacing-md);
  
  @media (max-width: 768px) {
    display: none;
  }
}

.search-input {
  border: none;
  background: transparent;
  outline: none;
  width: 160px;
  font-size: var(--font-sm);
  color: var(--text-primary);
  
  &::placeholder {
    color: var(--text-tertiary);
  }
}

.search-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--primary-color);
  color: white;
  font-size: var(--font-xs);
}

.theme-toggle {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--bg-secondary);
  font-size: var(--font-base);
  transition: all var(--transition-fast);
  
  &:hover {
    background: var(--bg-tertiary);
    transform: rotate(15deg);
  }
}

.write-btn {
  @media (max-width: 768px) {
    display: none;
  }
}

.user-menu {
  position: relative;
}

.user-avatar {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--primary-color);
  color: white;
  font-weight: 600;
  font-size: var(--font-sm);
  text-transform: uppercase;
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 180px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  width: 100%;
  padding: var(--spacing-sm) var(--spacing-md);
  text-align: left;
  font-size: var(--font-sm);
  color: var(--text-primary);
  transition: background var(--transition-fast);
  
  &:hover {
    background: var(--bg-secondary);
  }
  
  &.danger {
    color: var(--danger-color);
  }
}

.dropdown-divider {
  height: 1px;
  background: var(--border-color);
  margin: 4px 0;
}

.mobile-menu-btn {
  display: none;
  width: 36px;
  height: 36px;
  align-items: center;
  justify-content: center;
  font-size: var(--font-xl);
  
  @media (max-width: 768px) {
    display: flex;
  }
}

.mobile-menu {
  display: none;
  
  @media (max-width: 768px) {
    display: block;
    background: var(--bg-primary);
    border-bottom: 1px solid var(--border-color);
  }
}

.mobile-nav {
  display: flex;
  flex-direction: column;
  padding: var(--spacing-sm) 0;
}

.mobile-nav-link {
  padding: var(--spacing-sm) var(--spacing-lg);
  color: var(--text-secondary);
  font-size: var(--font-sm);
  
  &:hover, &.router-link-active {
    color: var(--primary-color);
    background: var(--bg-secondary);
  }
}
</style>
