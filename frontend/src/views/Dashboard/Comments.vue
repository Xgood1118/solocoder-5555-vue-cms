<template>
  <div class="comments-manage">
    <div class="page-header">
      <h1>评论审核</h1>
      <div class="tabs">
        <button 
          :class="{ active: activeTab === 'pending' }" 
          class="tab-btn"
          @click="activeTab = 'pending'"
        >
          待审核
          <span v-if="pendingCount > 0" class="badge badge-danger">{{ pendingCount }}</span>
        </button>
        <button 
          :class="{ active: activeTab === 'approved' }" 
          class="tab-btn"
          @click="activeTab = 'approved'"
        >
          已通过
        </button>
        <button 
          :class="{ active: activeTab === 'rejected' }" 
          class="tab-btn"
          @click="activeTab = 'rejected'"
        >
          已拒绝
        </button>
      </div>
    </div>
    
    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>
    
    <div v-else class="comments-list">
      <div v-if="comments.length === 0" class="empty-state">
        暂无评论
      </div>
      
      <div 
        v-for="comment in comments" 
        :key="comment.id" 
        class="comment-card"
      >
        <div class="comment-header">
          <div class="comment-author">
            <span class="avatar">{{ comment.author_name?.charAt(0)?.toUpperCase() || '?' }}</span>
            <div class="author-info">
              <span class="name">{{ comment.author_name }}</span>
              <span class="email">{{ comment.author_email }}</span>
            </div>
          </div>
          <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
        </div>
        
        <div class="comment-content">
          {{ comment.content }}
        </div>
        
        <div class="comment-meta">
          <span class="ip">IP: {{ comment.ip }}</span>
          <span v-if="comment.article_id" class="article">
            文章: {{ comment.article_id }}
          </span>
        </div>
        
        <div class="comment-actions">
          <template v-if="comment.status === 'pending'">
            <button class="btn btn-success btn-sm" @click="approveComment(comment)">
              ✓ 通过
            </button>
            <button class="btn btn-danger btn-sm" @click="rejectComment(comment)">
              ✗ 拒绝
            </button>
          </template>
          <template v-else-if="comment.status === 'approved'">
            <button class="btn btn-secondary btn-sm" @click="rejectComment(comment)">
              取消通过
            </button>
          </template>
          <template v-else>
            <button class="btn btn-success btn-sm" @click="approveComment(comment)">
              恢复
            </button>
          </template>
          <button class="btn btn-danger btn-sm" @click="deleteComment(comment)">
            删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { formatDate } from '@/utils'
import request from '@/utils/request'

const comments = ref([])
const loading = ref(false)
const activeTab = ref('pending')
const pendingCount = ref(0)

async function loadComments() {
  loading.value = true
  try {
    const res = await request.get('/comments', {
      params: { status: activeTab.value === 'pending' ? 'pending' : activeTab.value }
    })
    comments.value = res.data.comments || []
    
    if (activeTab.value === 'pending') {
      pendingCount.value = comments.value.length
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function approveComment(comment) {
  try {
    await request.post(`/comments/${comment.id}/approve`)
    loadComments()
    loadPendingCount()
  } catch (e) {
    alert('操作失败')
  }
}

async function rejectComment(comment) {
  try {
    await request.post(`/comments/${comment.id}/reject`)
    loadComments()
    loadPendingCount()
  } catch (e) {
    alert('操作失败')
  }
}

async function deleteComment(comment) {
  if (!confirm('确定要删除这条评论吗？')) return
  
  try {
    await request.delete(`/comments/${comment.id}`)
    loadComments()
    loadPendingCount()
  } catch (e) {
    alert('删除失败')
  }
}

async function loadPendingCount() {
  try {
    const res = await request.get('/comments/pending')
    pendingCount.value = res.data.comments?.length || 0
  } catch (e) {
    console.error(e)
  }
}

watch(activeTab, () => {
  loadComments()
})

onMounted(() => {
  loadComments()
  loadPendingCount()
})
</script>

<style lang="scss" scoped>
.comments-manage {
  padding: var(--spacing-lg);
}

.page-header {
  margin-bottom: var(--spacing-lg);
  
  h1 {
    font-size: var(--font-2xl);
    margin-bottom: var(--spacing-md);
  }
}

.tabs {
  display: flex;
  gap: var(--spacing-xs);
  border-bottom: 1px solid var(--border-color);
}

.tab-btn {
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-sm);
  color: var(--text-secondary);
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  
  &.active {
    color: var(--primary-color);
    border-bottom-color: var(--primary-color);
  }
  
  &:hover {
    color: var(--primary-color);
  }
}

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: var(--spacing-2xl) 0;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--text-tertiary);
  font-size: var(--font-sm);
}

.comment-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
  border: 1px solid var(--border-color);
}

.comment-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: var(--spacing-sm);
}

.comment-author {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.avatar {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--primary-color);
  color: white;
  font-size: var(--font-sm);
  font-weight: 600;
}

.author-info {
  display: flex;
  flex-direction: column;
}

.name {
  font-size: var(--font-sm);
  font-weight: 600;
  color: var(--text-primary);
}

.email {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.comment-date {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.comment-content {
  font-size: var(--font-sm);
  line-height: 1.6;
  color: var(--text-secondary);
  padding: var(--spacing-sm);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-sm);
}

.comment-meta {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-sm);
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.comment-actions {
  display: flex;
  gap: var(--spacing-xs);
}
</style>
