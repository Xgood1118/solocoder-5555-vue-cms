<template>
  <div class="tags-manage">
    <div class="page-header">
      <h1>标签管理</h1>
      <button class="btn btn-primary" @click="showModal = true; formData = { name: '', slug: '' }">
        + 新建标签
      </button>
    </div>
    
    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>
    
    <div v-else class="tags-grid">
      <div v-if="tags.length === 0" class="empty-state">
        暂无标签
      </div>
      
      <div 
        v-for="tag in tags" 
        :key="tag.id" 
        class="tag-card"
      >
        <div class="tag-info">
          <span class="tag-name">{{ tag.name }}</span>
          <span class="tag-count">{{ tag.count || 0 }} 篇文章</span>
        </div>
        <div class="tag-actions">
          <button class="btn btn-secondary btn-sm" @click="editTag(tag)">
            编辑
          </button>
          <button class="btn btn-danger btn-sm" @click="deleteTag(tag)">
            删除
          </button>
        </div>
      </div>
    </div>
    
    <div v-if="showModal" class="modal">
      <div class="modal-backdrop" @click="showModal = false"></div>
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ isEditing ? '编辑标签' : '新建标签' }}</h3>
          <button class="close-btn" @click="showModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标签名称 *</label>
            <input v-model="formData.name" type="text" placeholder="请输入标签名称" />
          </div>
          <div class="form-group">
            <label>Slug</label>
            <input v-model="formData.slug" type="text" placeholder="自动生成" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showModal = false">取消</button>
          <button class="btn btn-primary" :disabled="saving" @click="saveTag">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const tags = ref([])
const loading = ref(false)
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const formData = ref({ name: '', slug: '' })

async function loadTags() {
  loading.value = true
  try {
    const res = await request.get('/tags/cloud')
    tags.value = res.data.tags || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function editTag(tag) {
  isEditing.value = true
  formData.value = {
    id: tag.id,
    name: tag.name,
    slug: tag.slug
  }
  showModal.value = true
}

async function saveTag() {
  if (!formData.value.name) {
    alert('请输入标签名称')
    return
  }
  
  saving.value = true
  try {
    if (isEditing.value) {
      await request.put(`/tags/${formData.value.id}`, formData.value)
    } else {
      await request.post('/tags', formData.value)
    }
    
    showModal.value = false
    loadTags()
  } catch (e) {
    alert(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

async function deleteTag(tag) {
  if (!confirm(`确定要删除标签 "${tag.name}" 吗？`)) return
  
  try {
    await request.delete(`/tags/${tag.id}`)
    loadTags()
  } catch (e) {
    alert('删除失败')
  }
}

onMounted(() => {
  loadTags()
})
</script>

<style lang="scss" scoped>
.tags-manage {
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

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: var(--spacing-2xl) 0;
}

.tags-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: var(--spacing-md);
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--text-tertiary);
  font-size: var(--font-sm);
  grid-column: 1 / -1;
}

.tag-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
  border: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  transition: all var(--transition-fast);
  
  &:hover {
    box-shadow: var(--shadow-md);
  }
}

.tag-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tag-name {
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--text-primary);
}

.tag-count {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.tag-actions {
  display: flex;
  gap: var(--spacing-xs);
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-backdrop {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--bg-overlay);
}

.modal-content {
  position: relative;
  width: 90%;
  max-width: 400px;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  
  h3 {
    font-size: var(--font-lg);
  }
  
  .close-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    font-size: var(--font-xl);
    color: var(--text-secondary);
    
    &:hover {
      background: var(--bg-secondary);
    }
  }
}

.modal-body {
  padding: var(--spacing-md);
}

.form-group {
  margin-bottom: var(--spacing-md);
  
  label {
    display: block;
    margin-bottom: 4px;
    font-size: var(--font-xs);
    color: var(--text-secondary);
    font-weight: 500;
  }
  
  input {
    width: 100%;
    padding: var(--spacing-xs) var(--spacing-sm);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    font-size: var(--font-sm);
    background: var(--bg-primary);
    color: var(--text-primary);
    
    &:focus {
      outline: none;
      border-color: var(--primary-color);
    }
  }
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
  padding: var(--spacing-md);
  border-top: 1px solid var(--border-color);
}
</style>
