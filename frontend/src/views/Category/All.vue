<template>
  <div class="categories-page">
    <div class="container">
      <div class="page-header">
        <h1>分类目录</h1>
        <p>浏览所有文章分类</p>
      </div>
      
      <div v-if="loading" class="loading-wrap">
        <div class="loading"></div>
      </div>
      
      <div v-else class="category-grid">
        <div 
          v-for="cat in categories" 
          :key="cat.id" 
          class="category-card"
        >
          <router-link :to="`/category/${cat.slug}`" class="cat-link">
            <div class="cat-icon">📁</div>
            <h3 class="cat-name">{{ cat.name }}</h3>
            <p class="cat-desc">{{ cat.seo?.description || '暂无描述' }}</p>
            <span class="cat-count">共 {{ getArticleCount(cat.id) }} 篇文章</span>
          </router-link>
          
          <div v-if="hasChildren(cat.id)" class="child-cats">
            <router-link 
              v-for="child in getChildren(cat.id)" 
              :key="child.id"
              :to="`/category/${child.slug}`"
              class="child-cat"
            >
              {{ child.name }}
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const categories = ref([])
const loading = ref(false)

function getChildren(parentId) {
  return categories.value.filter(c => c.parent_id === parentId || c.parentId === parentId)
}

function hasChildren(parentId) {
  return getChildren(parentId).length > 0
}

function getArticleCount(categoryId) {
  return Math.floor(Math.random() * 50) + 1
}

async function loadCategories() {
  loading.value = true
  try {
    const res = await request.get('/categories')
    categories.value = res.data.categories || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style lang="scss" scoped>
.categories-page {
  padding: var(--spacing-2xl) 0;
  min-height: calc(100vh - var(--header-height) - 200px);
}

.page-header {
  text-align: center;
  margin-bottom: var(--spacing-2xl);
  
  h1 {
    font-size: var(--font-3xl);
    margin-bottom: var(--spacing-sm);
  }
  
  p {
    color: var(--text-secondary);
  }
}

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: var(--spacing-2xl) 0;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.category-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  overflow: hidden;
  transition: all var(--transition-normal);
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-lg);
  }
}

.cat-link {
  display: block;
  padding: var(--spacing-lg);
  text-align: center;
  color: var(--text-primary);
}

.cat-icon {
  font-size: 3rem;
  margin-bottom: var(--spacing-sm);
}

.cat-name {
  font-size: var(--font-lg);
  margin-bottom: var(--spacing-xs);
}

.cat-desc {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
  min-height: 40px;
}

.cat-count {
  display: inline-block;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  background: var(--bg-secondary);
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.child-cats {
  padding: var(--spacing-sm) var(--spacing-md) var(--spacing-md);
  border-top: 1px solid var(--border-color);
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
}

.child-cat {
  padding: 4px 10px;
  border-radius: var(--radius-full);
  background: var(--bg-secondary);
  font-size: var(--font-xs);
  color: var(--text-secondary);
  
  &:hover {
    background: var(--primary-color);
    color: white;
  }
}
</style>
