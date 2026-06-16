<template>
  <div class="comment-section">
    <h3>评论 ({{ totalComments }})</h3>
    
    <div class="comment-form">
      <h4>发表评论</h4>
      
      <div class="form-row">
        <div class="form-group">
          <label>昵称 *</label>
          <input v-model="form.author_name" type="text" placeholder="请输入昵称" />
        </div>
        <div class="form-group">
          <label>邮箱 * (不会公开)</label>
          <input v-model="form.author_email" type="email" placeholder="请输入邮箱" />
        </div>
      </div>
      
      <div class="form-group">
        <label>评论内容 *</label>
        <textarea 
          v-model="form.content" 
          :placeholder="replyingTo ? `回复 @${replyingTo.author_name}` : '写下你的评论...'"
          rows="4"
        ></textarea>
      </div>
      
      <div class="form-actions">
        <button v-if="replyingTo" class="btn btn-secondary btn-sm" @click="cancelReply">
          取消回复
        </button>
        <button class="btn btn-primary" :disabled="submitting" @click="submitComment">
          {{ submitting ? '提交中...' : '提交评论' }}
        </button>
      </div>
      
      <p class="tip">💡 评论需要审核后才能显示，请耐心等待</p>
    </div>
    
    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>
    
    <div v-else class="comment-list">
      <div v-if="comments.length === 0" class="empty-comments">
        暂无评论，快来抢沙发吧~
      </div>
      
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <div class="comment-header">
          <span class="comment-author">{{ comment.author_name }}</span>
          <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
        </div>
        <div class="comment-content">{{ comment.content }}</div>
        <div class="comment-actions">
          <button class="reply-btn" @click="startReply(comment)">回复</button>
        </div>
        
        <div v-if="comment.replies && comment.replies.length > 0" class="comment-replies">
          <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
            <div class="comment-header">
              <span class="comment-author">{{ reply.author_name }}</span>
              <span class="comment-date">{{ formatDate(reply.created_at) }}</span>
            </div>
            <div class="comment-content">{{ reply.content }}</div>
            <div class="comment-actions">
              <button v-if="reply.depth < 3" class="reply-btn" @click="startReply(reply)">回复</button>
            </div>
            
            <div v-if="reply.replies && reply.replies.length > 0" class="comment-replies level-2">
              <div v-for="r2 in reply.replies" :key="r2.id" class="reply-item">
                <div class="comment-header">
                  <span class="comment-author">{{ r2.author_name }}</span>
                  <span class="comment-date">{{ formatDate(r2.created_at) }}</span>
                </div>
                <div class="comment-content">{{ r2.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { formatDate } from '@/utils'
import request from '@/utils/request'

const props = defineProps({
  articleId: {
    type: String,
    required: true
  }
})

const comments = ref([])
const loading = ref(false)
const submitting = ref(false)
const replyingTo = ref(null)

const form = ref({
  author_name: '',
  author_email: '',
  content: '',
  parent_id: ''
})

const totalComments = computed(() => {
  let count = 0
  function countComments(list) {
    for (const c of list) {
      count++
      if (c.replies) {
        countComments(c.replies)
      }
    }
  }
  countComments(comments.value)
  return count
})

async function loadComments() {
  loading.value = true
  try {
    const res = await request.get('/comments', {
      params: { article_id: props.articleId }
    })
    comments.value = res.data.comments || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function startReply(comment) {
  replyingTo.value = comment
  form.value.parent_id = comment.id
}

function cancelReply() {
  replyingTo.value = null
  form.value.parent_id = ''
}

async function submitComment() {
  if (!form.value.author_name || !form.value.author_email || !form.value.content) {
    alert('请填写完整信息')
    return
  }
  
  submitting.value = true
  try {
    await request.post('/comments', {
      article_id: props.articleId,
      parent_id: form.value.parent_id,
      author_name: form.value.author_name,
      author_email: form.value.author_email,
      content: form.value.content
    })
    
    alert('评论提交成功，等待审核')
    form.value.content = ''
    cancelReply()
  } catch (e) {
    alert(e.response?.data?.error || '提交失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadComments()
})
</script>

<style lang="scss" scoped>
.comment-section {
  max-width: 800px;
  margin: var(--spacing-2xl) auto 0;
  padding-top: var(--spacing-xl);
  border-top: 1px solid var(--border-color);
  
  h3 {
    font-size: var(--font-xl);
    margin-bottom: var(--spacing-lg);
  }
}

.comment-form {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-xl);
  border: 1px solid var(--border-color);
  
  h4 {
    font-size: var(--font-base);
    margin-bottom: var(--spacing-md);
  }
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-md);
  
  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }
}

.form-group {
  margin-bottom: var(--spacing-md);
  
  label {
    display: block;
    margin-bottom: var(--spacing-xs);
    font-size: var(--font-sm);
    color: var(--text-secondary);
  }
  
  input, textarea {
    width: 100%;
    padding: var(--spacing-sm) var(--spacing-md);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    font-size: var(--font-sm);
    background: var(--bg-primary);
    color: var(--text-primary);
    
    &:focus {
      outline: none;
      border-color: var(--primary-color);
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }
  }
  
  textarea {
    resize: vertical;
    min-height: 80px;
  }
}

.form-actions {
  display: flex;
  gap: var(--spacing-sm);
  justify-content: flex-end;
}

.tip {
  margin-top: var(--spacing-sm);
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: var(--spacing-xl) 0;
}

.empty-comments {
  text-align: center;
  padding: var(--spacing-xl) 0;
  color: var(--text-tertiary);
  font-size: var(--font-sm);
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.comment-item {
  background: var(--bg-card);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  border: 1px solid var(--border-color);
}

.reply-item {
  padding: var(--spacing-sm);
  margin-top: var(--spacing-sm);
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
}

.comment-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-xs);
}

.comment-author {
  font-weight: 600;
  font-size: var(--font-sm);
  color: var(--text-primary);
}

.comment-date {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.comment-content {
  font-size: var(--font-sm);
  line-height: 1.6;
  color: var(--text-secondary);
}

.comment-actions {
  margin-top: var(--spacing-xs);
}

.reply-btn {
  font-size: var(--font-xs);
  color: var(--primary-color);
  background: none;
  padding: 0;
  
  &:hover {
    text-decoration: underline;
  }
}

.comment-replies {
  margin-top: var(--spacing-sm);
  padding-left: var(--spacing-md);
  border-left: 2px solid var(--border-color);
  
  &.level-2 {
    margin-top: var(--spacing-xs);
  }
}
</style>
