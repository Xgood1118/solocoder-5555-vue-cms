<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useThemeStore } from '@/stores/theme'
import AppHeader from '@/components/Layout/Header.vue'
import AppSidebar from '@/components/Layout/Sidebar.vue'
import AppFooter from '@/components/Layout/Footer.vue'

const route = useRoute()
const themeStore = useThemeStore()

const isDashboard = computed(() => route.path.startsWith('/dashboard'))
const isAuthPage = computed(() => ['/login', '/register'].includes(route.path))

const htmlClass = computed(() => ({
  'theme-dark': themeStore.isDark,
  'theme-light': !themeStore.isDark
}))
</script>

<template>
  <div :class="htmlClass" class="app-container">
    <AppHeader v-if="!isAuthPage" />
    <div class="main-wrapper" :class="{ 'dashboard-layout': isDashboard }">
      <AppSidebar v-if="isDashboard" />
      <main class="main-content">
        <router-view />
      </main>
    </div>
    <AppFooter v-if="!isAuthPage && !isDashboard" />
  </div>
</template>

<style lang="scss" scoped>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: background 0.3s, color 0.3s;
}

.main-wrapper {
  flex: 1;
  display: flex;
}

.main-content {
  flex: 1;
  min-width: 0;
}

.dashboard-layout {
  .main-content {
    padding: 24px;
    background: var(--bg-secondary);
  }
}

@media (max-width: 768px) {
  .dashboard-layout .main-content {
    padding: 16px;
  }
}
</style>
