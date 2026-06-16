<template>
  <div class="tag-page">
    <div class="container">
      <div v-if="loading" class="loading-wrap">
        <div class="loading"></div>
      </div>
      
      <template v-else>
        <div class="page-header">
          <h1>🏷️ {{ tag?.name || '标签' }}</h1>
          <p>共 {{ articles.length }} 篇文章</p>
        </div>
        
        <div v-if="articles.length === 0" class="empty-state">
          <p>该标签下暂无文章</p>
        </div>
        
        <div v-else class="article-grid">
          <article 
            v-for="article in articles" 
            :key="article.id" 
            class="article-card"
          >
            <div v-if="article.cover" class="article-cover">
              <img :src="article.cover" :alt="article.title" />
            </div>
            
            <div class="article-body">
              <h3 class="article-title">
                <router-link :to="`/article/${article.slug || article.id}`">
                  {{ article.title }}
                </router-link>
              </h3>
              
              <p class="article-summary">{{ article.summary }}</p>
              
              <div class="article-meta">
                <span>{{ formatDate(article.publish_at || article.created_at) }}</span>
              </div>
            </div>
          </article>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { formatDate } from '@/utils'
import request from '@/utils/request'

const route = useRoute()
const tag = ref(null)
const articles = ref([])
const loading = ref(false)

async function loadTag() {
  loading.value = true
  try {
    const id = route.params.id
    const res = await request.get(`/tags/${id}`)
    tag.value = res.data.tag
    articles.value = res.data.articles || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadTag()
})

watch(() => route.params.id, () => {
  loadTag()
})
</script>

<style lang="scss" scoped>
.tag-page {
  padding: var(--spacing-xl) 0;
  min-height: calc(100vh - var(--header-height) - 200px);
}

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: var(--spacing-2xl) 0;
}

.page-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  
  h1 {
    font-size: var(--font-3xl);
    margin-bottom: var(--spacing-sm);
  }
  
  p {
    color: var(--text-secondary);
  }
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl) 0;
  color: var(--text-tertiary);
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.article-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--border-color);
  transition: all var(--transition-fast);
  
  &:hover {
    transform: translateY(-3px);
    box-shadow: var(--shadow-lg);
  }
}

.article-cover {
  width: 100%;
  height: 160px;
  overflow: hidden;
  background: var(--bg-secondary);
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.article-body {
  padding: var(--spacing-md);
}

.article-title {
  font-size: var(--font-base);
  margin-bottom: var(--spacing-sm);
  font-weight: 600;
  
  a {
    color: var(--text-primary);
    
    &:hover {
      color: var(--primary-color);
    }
  }
}

.article-summary {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: var(--spacing-sm);
}

.article-meta {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}
</style>
