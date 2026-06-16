<template>
  <div class="dashboard-page">
    <h1 class="page-title">仪表板</h1>
    
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📝</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.total_articles || 0 }}</span>
          <span class="stat-label">总文章数</span>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">👁️</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.total_views || 0 }}</span>
          <span class="stat-label">总阅读量</span>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">👥</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.total_users || 0 }}</span>
          <span class="stat-label">用户总数</span>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">💬</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.pending_comments || 0 }}</span>
          <span class="stat-label">待审核评论</span>
        </div>
      </div>
    </div>
    
    <div class="today-stats">
      <div class="today-card">
        <h3>今日数据</h3>
        <div class="today-grid">
          <div class="today-item">
            <span class="today-label">今日 UV</span>
            <span class="today-value">{{ stats.today_uv || 0 }}</span>
          </div>
          <div class="today-item">
            <span class="today-label">今日 PV</span>
            <span class="today-value">{{ stats.today_pv || 0 }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="chart-section">
      <h3>最近 7 天访问趋势</h3>
      <div class="chart-container">
        <div class="bar-chart">
          <div 
            v-for="(day, index) in last7Days" 
            :key="index" 
            class="bar-item"
          >
            <div class="bar" :style="{ height: getBarHeight(day.pv) + '%' }"></div>
            <span class="bar-label">{{ day.date }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="popular-section">
      <h3>热门文章 Top 10</h3>
      <div class="popular-list">
        <div 
          v-for="(article, index) in popularArticles.slice(0, 10)" 
          :key="article.id" 
          class="popular-item"
        >
          <span class="rank" :class="{ top: index < 3 }">{{ index + 1 }}</span>
          <router-link :to="`/article/${article.slug || article.id}`" class="title">
            {{ article.title }}
          </router-link>
          <span class="views">{{ article.view_count }} 阅读</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import request from '@/utils/request'

const stats = ref({})
const popularArticles = ref([])

const last7Days = computed(() => {
  const days = stats.value.last_7_days || []
  return days.map(d => ({
    date: d.date?.slice(5) || '',
    pv: d.pv || 0,
    uv: d.uv || 0
  }))
})

function getBarHeight(pv) {
  const max = Math.max(...last7Days.value.map(d => d.pv), 1)
  return (pv / max) * 100
}

async function loadStats() {
  try {
    const res = await request.get('/admin/dashboard')
    stats.value = res.data || {}
    popularArticles.value = res.data.popular_articles || []
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadStats()
})
</script>

<style lang="scss" scoped>
.dashboard-page {
  padding: var(--spacing-lg);
}

.page-title {
  font-size: var(--font-2xl);
  margin-bottom: var(--spacing-lg);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-xl);
}

.stat-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
}

.stat-icon {
  font-size: var(--font-3xl);
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: var(--font-2xl);
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  font-size: var(--font-sm);
  color: var(--text-secondary);
}

.today-stats {
  margin-bottom: var(--spacing-xl);
}

.today-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
  
  h3 {
    font-size: var(--font-lg);
    margin-bottom: var(--spacing-md);
  }
}

.today-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-lg);
}

.today-item {
  text-align: center;
  padding: var(--spacing-md);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.today-label {
  display: block;
  font-size: var(--font-sm);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-xs);
}

.today-value {
  font-size: var(--font-2xl);
  font-weight: 700;
  color: var(--primary-color);
}

.chart-section {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-xl);
  border: 1px solid var(--border-color);
  
  h3 {
    font-size: var(--font-lg);
    margin-bottom: var(--spacing-lg);
  }
}

.chart-container {
  height: 250px;
  display: flex;
  align-items: flex-end;
}

.bar-chart {
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  width: 100%;
  height: 200px;
  padding: 0 var(--spacing-md);
}

.bar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  height: 100%;
  justify-content: flex-end;
}

.bar {
  width: 60%;
  max-width: 40px;
  background: linear-gradient(to top, var(--primary-color), var(--primary-light));
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
  transition: height 0.3s ease;
  min-height: 2px;
}

.bar-label {
  margin-top: var(--spacing-xs);
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.popular-section {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
  
  h3 {
    font-size: var(--font-lg);
    margin-bottom: var(--spacing-md);
  }
}

.popular-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.popular-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm);
  border-radius: var(--radius-md);
  transition: background var(--transition-fast);
  
  &:hover {
    background: var(--bg-secondary);
  }
}

.rank {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  background: var(--bg-secondary);
  font-size: var(--font-sm);
  font-weight: 600;
  color: var(--text-tertiary);
  
  &.top {
    background: var(--primary-color);
    color: white;
  }
}

.title {
  flex: 1;
  font-size: var(--font-sm);
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  
  &:hover {
    color: var(--primary-color);
  }
}

.views {
  flex-shrink: 0;
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .stat-card {
    padding: var(--spacing-md);
  }
  
  .stat-icon {
    font-size: var(--font-2xl);
  }
  
  .stat-value {
    font-size: var(--font-xl);
  }
}
</style>
