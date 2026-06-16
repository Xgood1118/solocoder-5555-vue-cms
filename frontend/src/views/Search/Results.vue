<template>
  <div class="search-page">
    <div class="container">
      <div class="search-header">
        <h1>搜索</h1>
        <div class="search-box">
          <input 
            v-model="keyword" 
            type="text" 
            placeholder="输入关键词搜索文章..." 
            @keyup.enter="doSearch"
            class="search-input"
          />
          <button class="btn btn-primary" @click="doSearch">搜索</button>
        </div>
      </div>
      
      <div v-if="loading" class="loading-wrap">
        <div class="loading"></div>
      </div>
      
      <div v-else-if="searched && articles.length === 0" class="empty-state">
        <p>未找到与 "<strong>{{ keyword }}</strong>" 相关的文章</p>
        <p class="hint">试试其他关键词吧</p>
      </div>
      
      <template v-else-if="articles.length > 0">
        <div class="search-results">
          <p class="result-count">找到 {{ total }} 篇相关文章</p>
          
          <div class="article-list">
            <article 
              v-for="article in articles" 
              :key="article.id" 
              class="article-item"
            >
              <h3 class="article-title">
                <router-link :to="`/article/${article.slug || article.id}`" v-html="highlightTitle(article.title)"></router-link>
              </h3>
              
              <p class="article-summary" v-html="article.summary"></p>
              
              <div class="article-meta">
                <span>{{ formatDate(article.publish_at || article.created_at) }}</span>
                <span>·</span>
                <span>{{ article.view_count }} 阅读</span>
              </div>
            </article>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { formatDate, highlightText } from '@/utils'
import request from '@/utils/request'

const route = useRoute()
const keyword = ref('')
const articles = ref([])
const total = ref(0)
const loading = ref(false)
const searched = ref(false)

function highlightTitle(title) {
  return highlightText(title, keyword.value)
}

async function doSearch() {
  if (!keyword.value.trim()) return
  
  loading.value = true
  searched.value = true
  
  try {
    const res = await request.get('/articles/search', {
      params: { q: keyword.value }
    })
    articles.value = res.data.articles || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  const q = route.query.q
  if (q) {
    keyword.value = q
    doSearch()
  }
})
</script>

<style lang="scss" scoped>
.search-page {
  padding: var(--spacing-xl) 0;
  min-height: calc(100vh - var(--header-height) - 200px);
}

.search-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  
  h1 {
    font-size: var(--font-3xl);
    margin-bottom: var(--spacing-lg);
  }
}

.search-box {
  display: flex;
  gap: var(--spacing-sm);
  max-width: 600px;
  margin: 0 auto;
}

.search-input {
  flex: 1;
  padding: var(--spacing-sm) var(--spacing-md);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  font-size: var(--font-base);
  background: var(--bg-primary);
  color: var(--text-primary);
  
  &:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
  }
}

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: var(--spacing-2xl) 0;
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl) 0;
  
  p {
    color: var(--text-secondary);
    
    &.hint {
      color: var(--text-tertiary);
      font-size: var(--font-sm);
      margin-top: var(--spacing-sm);
    }
  }
  
  strong {
    color: var(--primary-color);
  }
}

.search-results {
  max-width: 800px;
  margin: 0 auto;
}

.result-count {
  color: var(--text-secondary);
  font-size: var(--font-sm);
  margin-bottom: var(--spacing-lg);
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.article-item {
  padding: var(--spacing-md);
  background: var(--bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  
  &:hover {
    box-shadow: var(--shadow-sm);
  }
}

.article-title {
  font-size: var(--font-lg);
  margin-bottom: var(--spacing-sm);
  
  a {
    color: var(--text-primary);
    
    &:hover {
      color: var(--primary-color);
    }
  }
  
  :deep(mark) {
    background: #fef08a;
    padding: 0 2px;
    border-radius: 2px;
  }
}

.theme-dark .article-title :deep(mark) {
  background: #854d0e;
}

.article-summary {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: var(--spacing-sm);
  
  :deep(mark) {
    background: #fef08a;
    padding: 0 2px;
    border-radius: 2px;
  }
}

.theme-dark .article-summary :deep(mark) {
  background: #854d0e;
}

.article-meta {
  display: flex;
  gap: var(--spacing-sm);
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}
</style>
