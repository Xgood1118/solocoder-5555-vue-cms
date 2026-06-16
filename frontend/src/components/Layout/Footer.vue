<template>
  <footer class="app-footer">
    <div class="container footer-inner">
      <div class="footer-links">
        <router-link to="/">首页</router-link>
        <a href="/api/rss">RSS</a>
        <a href="/api/sitemap.xml">Sitemap</a>
        <router-link to="/about">关于</router-link>
      </div>
      
      <div class="footer-info">
        <p>© {{ currentYear }} 轻量 CMS. 基于 Vue 3 + Go 构建</p>
        <p class="powered">PWA 支持 · 离线可访问</p>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { computed } from 'vue'

const currentYear = computed(() => new Date().getFullYear())
</script>

<style lang="scss" scoped>
.app-footer {
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
  padding: var(--spacing-xl) 0;
  margin-top: var(--spacing-2xl);
}

.footer-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-md);
}

.footer-links {
  display: flex;
  gap: var(--spacing-lg);
  
  a {
    color: var(--text-secondary);
    font-size: var(--font-sm);
    
    &:hover {
      color: var(--primary-color);
    }
  }
}

.footer-info {
  text-align: center;
  
  p {
    color: var(--text-tertiary);
    font-size: var(--font-sm);
    margin-bottom: 4px;
  }
  
  .powered {
    font-size: var(--font-xs);
    opacity: 0.7;
  }
}
</style>
