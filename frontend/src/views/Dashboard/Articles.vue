<template>
  <div class="articles-manage">
    <div class="page-header">
      <h1>文章管理</h1>
      <router-link to="/article/editor" class="btn btn-primary">
        + 新建文章
      </router-link>
    </div>
    
    <div class="filters">
      <div class="filter-item">
        <label>状态</label>
        <select v-model="filterStatus" @change="loadArticles">
          <option value="">全部</option>
          <option value="published">已发布</option>
          <option value="draft">草稿</option>
          <option value="scheduled">定时发布</option>
        </select>
      </div>
      
      <div class="filter-item">
        <label>搜索</label>
        <input 
          v-model="searchKeyword" 
          type="text" 
          placeholder="搜索文章标题..."
          @keyup.enter="loadArticles"
        />
      </div>
    </div>
    
    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>
    
    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>标题</th>
            <th>状态</th>
            <th>分类</th>
            <th>阅读量</th>
            <th>发布时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="article in articles" :key="article.id">
            <td class="title-cell">
              <router-link :to="`/article/${article.slug || article.id}`">
                {{ article.title }}
              </router-link>
            </td>
            <td>
              <span class="status-badge" :class="article.status">
                {{ getStatusLabel(article.status) }}
              </span>
            </td>
            <td>{{ getCategoryName(article.category_id) }}</td>
            <td>{{ article.view_count }}</td>
            <td>{{ formatDate(article.publish_at || article.created_at) }}</td>
            <td class="actions">
              <router-link 
                :to="`/article/editor/${article.id}`" 
                class="btn btn-secondary btn-sm"
              >
                编辑
              </router-link>
              <button 
                class="btn btn-danger btn-sm" 
                @click="deleteArticle(article)"
              >
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="articles.length === 0" class="empty-state">
        暂无文章
      </div>
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
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { formatDate } from '@/utils'
import request from '@/utils/request'

const articles = ref([])
const categories = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filterStatus = ref('')
const searchKeyword = ref('')

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

function getStatusLabel(status) {
  const labels = {
    published: '已发布',
    draft: '草稿',
    scheduled: '定时发布'
  }
  return labels[status] || status
}

function getCategoryName(categoryId) {
  const cat = categories.value.find(c => c.id === categoryId)
  return cat?.name || '未分类'
}

async function loadArticles() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    
    const res = await request.get('/articles', { params })
    articles.value = res.data.articles || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
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

async function deleteArticle(article) {
  if (!confirm(`确定要删除文章 "${article.title}" 吗？`)) return
  
  try {
    await request.delete(`/articles/${article.id}`)
    loadArticles()
  } catch (e) {
    alert('删除失败')
  }
}

onMounted(() => {
  loadCategories()
  loadArticles()
})
</script>

<style lang="scss" scoped>
.articles-manage {
  padding: var(--spacing-lg);
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-lg);
  
  h1 {
    font-size: var(--font-2xl);
  }
}

.filters {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  
  label {
    font-size: var(--font-xs);
    color: var(--text-secondary);
  }
  
  select, input {
    padding: var(--spacing-xs) var(--spacing-sm);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    font-size: var(--font-sm);
    background: var(--bg-card);
    color: var(--text-primary);
    min-width: 150px;
  }
}

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: var(--spacing-2xl) 0;
}

.table-container {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  
  th, td {
    padding: var(--spacing-sm) var(--spacing-md);
    text-align: left;
    border-bottom: 1px solid var(--border-color);
    font-size: var(--font-sm);
  }
  
  th {
    background: var(--bg-secondary);
    font-weight: 600;
    color: var(--text-secondary);
  }
  
  tbody tr:hover {
    background: var(--bg-secondary);
  }
}

.title-cell {
  max-width: 300px;
  
  a {
    color: var(--text-primary);
    
    &:hover {
      color: var(--primary-color);
    }
  }
}

.status-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: var(--radius-full);
  font-size: var(--font-xs);
  font-weight: 500;
  
  &.published {
    background: #d1fae5;
    color: #065f46;
  }
  
  &.draft {
    background: #e5e7eb;
    color: #374151;
  }
  
  &.scheduled {
    background: #fef3c7;
    color: #92400e;
  }
}

.theme-dark .status-badge {
  &.published {
    background: rgba(16, 185, 129, 0.2);
    color: #34d399;
  }
  
  &.draft {
    background: rgba(156, 163, 175, 0.2);
    color: #9ca3af;
  }
  
  &.scheduled {
    background: rgba(251, 191, 36, 0.2);
    color: #fbbf24;
  }
}

.actions {
  display: flex;
  gap: var(--spacing-xs);
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--text-tertiary);
  font-size: var(--font-sm);
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-md);
  margin-top: var(--spacing-lg);
}

.page-info {
  font-size: var(--font-sm);
  color: var(--text-tertiary);
}
</style>
