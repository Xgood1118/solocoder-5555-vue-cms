<template>
  <div class="home-page">
    <section class="hero">
      <div class="container">
        <h1>欢迎来到轻量 CMS</h1>
        <p>一个基于 Vue 3 + Go 的现代化内容管理系统</p>
        <div class="hero-actions">
          <router-link to="/article/editor" class="btn btn-primary btn-lg" v-if="userStore.isLoggedIn">
            ✏️ 开始写作
          </router-link>
          <router-link to="/register" class="btn btn-outline btn-lg" v-else>
            立即注册
          </router-link>
        </div>
      </div>
    </section>
    
    <div class="container main-content">
      <div class="content-grid">
        <div class="articles-section">
          <div class="section-header">
            <h2>最新文章</h2>
            <router-link to="/search" class="view-all">查看全部 →</router-link>
          </div>
          
          <div v-if="loading" class="loading-wrap">
            <div class="loading"></div>
          </div>
          
          <div v-else-if="articles.length === 0" class="empty-state">
            <p>暂无文章</p>
          </div>
          
          <div v-else class="article-list">
            <article 
              v-for="article in articles" 
              :key="article.id" 
              class="article-card"
            >
              <div v-if="article.cover" class="article-cover">
                <img :src="article.cover" :alt="article.title" />
              </div>
              
              <div class="article-body">
                <div class="article-meta">
                  <span v-if="getCategory(article.category_id)" class="category">
                    {{ getCategory(article.category_id)?.name }}
                  </span>
                  <span class="date">{{ formatDate(article.publish_at || article.created_at) }}</span>
                  <span class="views">👁️ {{ article.view_count }}</span>
                </div>
                
                <h3 class="article-title">
                  <router-link :to="`/article/${article.slug || article.id}`">
                    {{ article.title }}
                  </router-link>
                </h3>
                
                <p class="article-summary">{{ article.summary }}</p>
                
                <div class="article-tags">
                  <router-link 
                    v-for="tag in article.tags?.slice(0, 3)" 
                    :key="tag"
                    :to="`/tag/${tag}`"
                    class="tag"
                  >
                    {{ tag }}
                  </router-link>
                </div>
              </div>
            </article>
          </div>
          
          <div v-if="total > pageSize" class="pagination">
            <button 
              class="btn btn-secondary btn-sm"
              :disabled="page <= 1"
              @click="page--; loadArticles()"
            >
              上一页
            </button>
            <span class="page-info">第 {{ page }} / {{ totalPages }} 页</span>
            <button 
              class="btn btn-secondary btn-sm"
              :disabled="page >= totalPages"
              @click="page++; loadArticles()"
            >
              下一页
            </button>
          </div>
        </div>
        
        <aside class="sidebar-section">
          <div class="card sidebar-card">
            <h3>热门文章</h3>
            <ul class="hot-list">
              <li v-for="(article, index) in hotArticles.slice(0, 5)" :key="article.id">
                <span class="rank" :class="{ top: index < 3 }">{{ index + 1 }}</span>
                <router-link :to="`/article/${article.slug || article.id}`">
                  {{ article.title }}
                </router-link>
              </li>
            </ul>
          </div>
          
          <div class="card sidebar-card">
            <h3>标签云</h3>
            <div class="tag-cloud">
              <router-link 
                v-for="tag in tagCloud.slice(0, 15)" 
                :key="tag.id"
                :to="`/tag/${tag.slug}`"
                class="tag"
                :style="{ fontSize: Math.min(12 + tag.count * 2, 20) + 'px' }"
              >
                {{ tag.name }}
              </router-link>
            </div>
          </div>
          
          <div class="card sidebar-card">
            <h3>分类目录</h3>
            <ul class="category-list">
              <li v-for="cat in categories.slice(0, 8)" :key="cat.id">
                <router-link :to="`/category/${cat.slug}`">
                  <span class="cat-name">{{ cat.name }}</span>
                </router-link>
              </li>
            </ul>
          </div>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { formatDate } from '@/utils'
import request from '@/utils/request'

const userStore = useUserStore()

const articles = ref([])
const hotArticles = ref([])
const categories = ref([])
const tagCloud = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

function getCategory(categoryId) {
  return categories.value.find(c => c.id === categoryId)
}

async function loadArticles() {
  loading.value = true
  try {
    const res = await request.get('/articles', {
      params: {
        page: page.value,
        page_size: pageSize.value,
        status: 'published'
      }
    })
    articles.value = res.data.articles || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadHotArticles() {
  try {
    const res = await request.get('/articles', {
      params: { status: 'published', page_size: 10 }
    })
    const list = res.data.articles || []
    hotArticles.value = list.sort((a, b) => b.view_count - a.view_count).slice(0, 5)
  } catch (e) {
    console.error(e)
  }
}

async function loadCategories() {
  try {
    const res = await request.get('/categories')
    categories.value = res.data.categories || []
  } catch (e) {
    console.error(e)
  }
}

async function loadTagCloud() {
  try {
    const res = await request.get('/tags/cloud')
    tagCloud.value = res.data.tags || []
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadArticles()
  loadHotArticles()
  loadCategories()
  loadTagCloud()
})
</script>

<style lang="scss" scoped>
.home-page {
  min-height: calc(100vh - var(--header-height) - 200px);
}

.hero {
  background: linear-gradient(135deg, var(--primary-color) 0%, #7c3aed 100%);
  color: white;
  padding: var(--spacing-2xl) 0;
  text-align: center;
  margin-bottom: var(--spacing-2xl);
  
  h1 {
    font-size: var(--font-4xl);
    margin-bottom: var(--spacing-md);
    color: white;
  }
  
  p {
    font-size: var(--font-lg);
    opacity: 0.9;
    margin-bottom: var(--spacing-lg);
  }
}

.hero-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: center;
}

.main-content {
  padding-bottom: var(--spacing-2xl);
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: var(--spacing-xl);
  
  @media (max-width: 1024px) {
    grid-template-columns: 1fr;
  }
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-lg);
  
  h2 {
    font-size: var(--font-2xl);
  }
}

.view-all {
  font-size: var(--font-sm);
  color: var(--primary-color);
}

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: var(--spacing-2xl) 0;
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
}

.article-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--border-color);
  transition: all var(--transition-normal);
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-lg);
  }
}

.article-cover {
  width: 100%;
  height: 200px;
  overflow: hidden;
  background: var(--bg-secondary);
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.article-body {
  padding: var(--spacing-lg);
}

.article-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-sm);
  font-size: var(--font-xs);
  color: var(--text-tertiary);
  
  .category {
    color: var(--primary-color);
    font-weight: 500;
  }
}

.article-title {
  font-size: var(--font-xl);
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
  margin-bottom: var(--spacing-md);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-md);
  margin-top: var(--spacing-xl);
}

.page-info {
  font-size: var(--font-sm);
  color: var(--text-tertiary);
}

.sidebar-section {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.sidebar-card {
  padding: var(--spacing-lg);
  
  h3 {
    font-size: var(--font-lg);
    margin-bottom: var(--spacing-md);
    padding-bottom: var(--spacing-sm);
    border-bottom: 1px solid var(--border-color);
  }
}

.hot-list {
  li {
    display: flex;
    align-items: flex-start;
    gap: var(--spacing-sm);
    padding: var(--spacing-sm) 0;
    border-bottom: 1px dashed var(--border-light);
    
    &:last-child {
      border-bottom: none;
    }
    
    a {
      font-size: var(--font-sm);
      color: var(--text-secondary);
      line-height: 1.4;
      
      &:hover {
        color: var(--primary-color);
      }
    }
  }
}

.rank {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  background: var(--bg-secondary);
  font-size: var(--font-xs);
  font-weight: 600;
  color: var(--text-tertiary);
  
  &.top {
    background: var(--primary-color);
    color: white;
  }
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
}

.category-list {
  li {
    padding: var(--spacing-xs) 0;
    
    a {
      font-size: var(--font-sm);
      color: var(--text-secondary);
      
      &:hover {
        color: var(--primary-color);
      }
    }
  }
}
</style>
