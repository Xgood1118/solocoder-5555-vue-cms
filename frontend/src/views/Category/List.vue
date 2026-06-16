<template>
  <div class="category-page">
    <div class="container">
      <div v-if="loading" class="loading-wrap">
        <div class="loading"></div>
      </div>
      
      <template v-else>
        <div class="page-header">
          <h1>{{ category?.name || '分类' }}</h1>
          <p v-if="category?.seo?.description">{{ category.seo.description }}</p>
        </div>
        
        <div v-if="articles.length === 0" class="empty-state">
          <p>该分类下暂无文章</p>
        </div>
        
        <div v-else class="article-list">
          <article 
            v-for="article in articles" 
            :key="article.id" 
            class="article-item"
          >
            <div v-if="article.cover" class="article-cover">
              <img :src="article.cover" :alt="article.title" />
            </div>
            
            <div class="article-info">
              <h3 class="article-title">
                <router-link :to="`/article/${article.slug || article.id}`">
                  {{ article.title }}
                </router-link>
              </h3>
              
              <p class="article-summary">{{ article.summary }}</p>
              
              <div class="article-meta">
                <span>{{ formatDate(article.publish_at || article.created_at) }}</span>
                <span>·</span>
                <span>{{ article.view_count }} 阅读</span>
              </div>
              
              <div class="article-tags">
                <span v-for="tag in article.tags?.slice(0, 3)" :key="tag" class="tag">
                  {{ tag }}
                </span>
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
const category = ref(null)
const articles = ref([])
const loading = ref(false)

async function loadCategory() {
  loading.value = true
  try {
    const id = route.params.id
    const res = await request.get(`/categories/${id}`)
    category.value = res.data.category
    articles.value = res.data.articles || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategory()
})

watch(() => route.params.id, () => {
  loadCategory()
})
</script>

<style lang="scss" scoped>
.category-page {
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
  padding-bottom: var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
  
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

.article-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  max-width: 900px;
  margin: 0 auto;
}

.article-item {
  display: flex;
  gap: var(--spacing-lg);
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
  transition: all var(--transition-fast);
  
  &:hover {
    box-shadow: var(--shadow-md);
  }
  
  @media (max-width: 640px) {
    flex-direction: column;
  }
}

.article-cover {
  flex-shrink: 0;
  width: 240px;
  height: 150px;
  border-radius: var(--radius-md);
  overflow: hidden;
  background: var(--bg-secondary);
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  @media (max-width: 640px) {
    width: 100%;
    height: 200px;
  }
}

.article-info {
  flex: 1;
  min-width: 0;
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
}

.article-summary {
  color: var(--text-secondary);
  font-size: var(--font-sm);
  line-height: 1.6;
  margin-bottom: var(--spacing-sm);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: var(--font-xs);
  color: var(--text-tertiary);
  margin-bottom: var(--spacing-sm);
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
}
</style>
