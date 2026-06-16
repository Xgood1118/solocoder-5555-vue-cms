<template>
  <div class="article-detail">
    <div class="container">
      <div v-if="loading" class="loading-wrap">
        <div class="loading"></div>
      </div>
      
      <template v-else-if="article">
        <article class="article-content">
          <header class="article-header">
            <div class="article-breadcrumb">
              <router-link to="/">首页</router-link>
              <span class="sep">/</span>
              <span>{{ article.title }}</span>
            </div>
            
            <h1 class="article-title">{{ article.title }}</h1>
            
            <div class="article-meta">
              <span class="author">✍️ {{ getAuthorName(article.author_id) }}</span>
              <span class="date">📅 {{ formatDate(article.publish_at || article.created_at) }}</span>
              <span class="views">👁️ {{ article.view_count }} 阅读</span>
              <span class="category" v-if="category">
                📁 <router-link :to="`/category/${category.slug}`">{{ category.name }}</router-link>
              </span>
            </div>
            
            <div v-if="article.cover" class="article-cover">
              <img :src="article.cover" :alt="article.title" />
            </div>
          </header>
          
          <div class="article-body">
            <div v-html="renderedContent" class="markdown-body"></div>
          </div>
          
          <footer class="article-footer">
            <div class="article-tags">
              <span class="label">标签：</span>
              <router-link 
                v-for="tag in article.tags" 
                :key="tag"
                :to="`/tag/${tag}`"
                class="tag"
              >
                {{ tag }}
              </router-link>
            </div>
            
            <div class="article-actions">
              <button class="btn btn-secondary" @click="copyLink">🔗 复制链接</button>
            </div>
          </footer>
        </article>
        
        <CommentSection :article-id="article.id" />
      </template>
      
      <div v-else class="empty-state">
        <p>文章不存在或已删除</p>
        <router-link to="/" class="btn btn-primary">返回首页</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'
import { formatDate, copyToClipboard } from '@/utils'
import request from '@/utils/request'
import CommentSection from '@/components/Comment/CommentSection.vue'

const route = useRoute()
const article = ref(null)
const category = ref(null)
const loading = ref(false)

const renderedContent = computed(() => {
  if (!article.value) return ''
  
  let content = article.value.content || ''
  
  if (content.startsWith('<') && content.includes('<p')) {
    return content
  }
  
  try {
    marked.setOptions({
      highlight: function(code, lang) {
        if (lang && hljs.getLanguage(lang)) {
          try {
            return hljs.highlight(code, { language: lang }).value
          } catch (e) {}
        }
        return hljs.highlightAuto(code).value
      },
      breaks: true,
      gfm: true
    })
    return marked.parse(content)
  } catch (e) {
    return content
  }
})

function getAuthorName(authorId) {
  return '作者'
}

function copyLink() {
  copyToClipboard(window.location.href)
  alert('链接已复制')
}

async function loadArticle() {
  loading.value = true
  try {
    const id = route.params.id
    const res = await request.get(`/articles/${id}`)
    article.value = res.data.article
    
    if (article.value?.category_id) {
      const catRes = await request.get(`/categories/${article.value.category_id}`)
      category.value = catRes.data.category
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadArticle()
})

watch(() => route.params.id, () => {
  loadArticle()
})
</script>

<style lang="scss" scoped>
.article-detail {
  padding: var(--spacing-xl) 0;
  min-height: calc(100vh - var(--header-height) - 200px);
}

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: var(--spacing-2xl) 0;
}

.article-content {
  max-width: 800px;
  margin: 0 auto;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-2xl);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
}

.article-header {
  margin-bottom: var(--spacing-xl);
  padding-bottom: var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.article-breadcrumb {
  font-size: var(--font-sm);
  color: var(--text-tertiary);
  margin-bottom: var(--spacing-md);
  
  .sep {
    margin: 0 var(--spacing-xs);
  }
  
  a {
    color: var(--text-secondary);
    
    &:hover {
      color: var(--primary-color);
    }
  }
}

.article-title {
  font-size: var(--font-3xl);
  margin-bottom: var(--spacing-md);
  line-height: 1.3;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-md);
  font-size: var(--font-sm);
  color: var(--text-tertiary);
  
  .category a {
    color: var(--primary-color);
  }
}

.article-cover {
  margin-top: var(--spacing-lg);
  border-radius: var(--radius-md);
  overflow: hidden;
  
  img {
    width: 100%;
    height: auto;
  }
}

.article-body {
  font-size: var(--font-base);
  line-height: 1.8;
  color: var(--text-primary);
}

.article-footer {
  margin-top: var(--spacing-xl);
  padding-top: var(--spacing-lg);
  border-top: 1px solid var(--border-color);
}

.article-tags {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-lg);
  
  .label {
    font-size: var(--font-sm);
    color: var(--text-secondary);
  }
}

.article-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl) 0;
  
  p {
    color: var(--text-tertiary);
    margin-bottom: var(--spacing-lg);
  }
}

.markdown-body {
  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    margin-top: var(--spacing-lg);
    margin-bottom: var(--spacing-md);
    font-weight: 600;
  }
  
  :deep(h1) { font-size: var(--font-2xl); }
  :deep(h2) { font-size: var(--font-xl); }
  :deep(h3) { font-size: var(--font-lg); }
  
  :deep(p) {
    margin-bottom: var(--spacing-md);
  }
  
  :deep(pre) {
    background: var(--bg-secondary);
    border-radius: var(--radius-md);
    padding: var(--spacing-md);
    overflow-x: auto;
    margin-bottom: var(--spacing-md);
    position: relative;
    
    code {
      font-family: 'Fira Code', Monaco, Consolas, monospace;
      font-size: var(--font-sm);
    }
  }
  
  :deep(code) {
    background: var(--bg-secondary);
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 0.9em;
  }
  
  :deep(pre code) {
    background: none;
    padding: 0;
  }
  
  :deep(blockquote) {
    border-left: 4px solid var(--primary-color);
    padding-left: var(--spacing-md);
    margin: var(--spacing-md) 0;
    color: var(--text-secondary);
    font-style: italic;
  }
  
  :deep(ul), :deep(ol) {
    margin-bottom: var(--spacing-md);
    padding-left: var(--spacing-lg);
    
    li {
      margin-bottom: var(--spacing-xs);
    }
  }
  
  :deep(ul) {
    list-style: disc;
  }
  
  :deep(ol) {
    list-style: decimal;
  }
  
  :deep(img) {
    max-width: 100%;
    border-radius: var(--radius-md);
    margin: var(--spacing-md) 0;
  }
  
  :deep(a) {
    color: var(--primary-color);
    text-decoration: underline;
  }
  
  :deep(table) {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: var(--spacing-md);
    
    th, td {
      border: 1px solid var(--border-color);
      padding: var(--spacing-sm);
      text-align: left;
    }
    
    th {
      background: var(--bg-secondary);
      font-weight: 600;
    }
  }
  
  :deep(hr) {
    border: none;
    border-top: 1px solid var(--border-color);
    margin: var(--spacing-lg) 0;
  }
}

@media (max-width: 768px) {
  .article-content {
    padding: var(--spacing-lg);
  }
  
  .article-title {
    font-size: var(--font-2xl);
  }
}
</style>
