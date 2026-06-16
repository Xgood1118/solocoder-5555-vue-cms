<template>
  <div class="categories-manage">
    <div class="page-header">
      <h1>分类管理</h1>
      <button class="btn btn-primary" @click="showModal = true; formData = { name: '', slug: '', parent_id: '', template: 'default', seo: { title: '', description: '', keywords: '' } }">
        + 新建分类
      </button>
    </div>
    
    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>
    
    <div v-else class="tree-container">
      <div v-if="categories.length === 0" class="empty-state">
        暂无分类
      </div>
      
      <div class="category-tree">
        <div 
          v-for="cat in rootCategories" 
          :key="cat.id" 
          class="category-node"
        >
          <div class="node-content">
            <span class="node-name">{{ cat.name }}</span>
            <span class="node-slug">{{ cat.slug }}</span>
            <div class="node-actions">
              <button class="btn btn-secondary btn-sm" @click="editCategory(cat)">
                编辑
              </button>
              <button class="btn btn-danger btn-sm" @click="deleteCategory(cat)">
                删除
              </button>
            </div>
          </div>
          <div v-if="getChildCategories(cat.id).length > 0" class="child-categories">
            <div 
              v-for="child in getChildCategories(cat.id)" 
              :key="child.id" 
              class="category-node child"
            >
              <div class="node-content">
                <span class="node-name">└ {{ child.name }}</span>
                <span class="node-slug">{{ child.slug }}</span>
                <div class="node-actions">
                  <button class="btn btn-secondary btn-sm" @click="editCategory(child)">
                    编辑
                  </button>
                  <button class="btn btn-danger btn-sm" @click="deleteCategory(child)">
                    删除
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-if="showModal" class="modal">
      <div class="modal-backdrop" @click="showModal = false"></div>
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ isEditing ? '编辑分类' : '新建分类' }}</h3>
          <button class="close-btn" @click="showModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>分类名称 *</label>
            <input v-model="formData.name" type="text" placeholder="请输入分类名称" />
          </div>
          
          <div class="form-group">
            <label>Slug</label>
            <input v-model="formData.slug" type="text" placeholder="自动生成" />
          </div>
          
          <div class="form-group">
            <label>父级分类</label>
            <select v-model="formData.parent_id">
              <option value="">无（顶级分类）</option>
              <option v-for="cat in rootCategories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
          </div>
          
          <div class="form-group">
            <label>模板</label>
            <select v-model="formData.template">
              <option value="default">默认模板</option>
              <option value="blog">博客模板</option>
              <option value="doc">文档模板</option>
            </select>
          </div>
          
          <div class="seo-section">
            <h4>SEO 设置</h4>
            <div class="form-group">
              <label>SEO Title</label>
              <input v-model="formData.seo.title" type="text" placeholder="页面标题" />
            </div>
            <div class="form-group">
              <label>SEO Description</label>
              <textarea v-model="formData.seo.description" rows="2" placeholder="页面描述"></textarea>
            </div>
            <div class="form-group">
              <label>SEO Keywords</label>
              <input v-model="formData.seo.keywords" type="text" placeholder="关键词，逗号分隔" />
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showModal = false">取消</button>
          <button class="btn btn-primary" :disabled="saving" @click="saveCategory">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import request from '@/utils/request'

const categories = ref([])
const loading = ref(false)
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const formData = ref({
  id: '',
  name: '',
  slug: '',
  parent_id: '',
  template: 'default',
  seo: {
    title: '',
    description: '',
    keywords: ''
  }
})

const rootCategories = computed(() => {
  return categories.value.filter(c => !c.parent_id || c.parent_id === '')
})

function getChildCategories(parentId) {
  return categories.value.filter(c => c.parent_id === parentId)
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

function editCategory(cat) {
  isEditing.value = true
  formData.value = {
    id: cat.id,
    name: cat.name,
    slug: cat.slug,
    parent_id: cat.parent_id || '',
    template: cat.template || 'default',
    seo: cat.seo || { title: '', description: '', keywords: '' }
  }
  showModal.value = true
}

async function saveCategory() {
  if (!formData.value.name) {
    alert('请输入分类名称')
    return
  }
  
  saving.value = true
  try {
    if (isEditing.value) {
      await request.put(`/categories/${formData.value.id}`, formData.value)
    } else {
      await request.post('/categories', formData.value)
    }
    
    showModal.value = false
    loadCategories()
  } catch (e) {
    alert(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

async function deleteCategory(cat) {
  if (!confirm(`确定要删除分类 "${cat.name}" 吗？子分类也会被一并删除。`)) return
  
  try {
    await request.delete(`/categories/${cat.id}`)
    loadCategories()
  } catch (e) {
    alert('删除失败')
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style lang="scss" scoped>
.categories-manage {
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

.tree-container {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  padding: var(--spacing-md);
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--text-tertiary);
  font-size: var(--font-sm);
}

.category-tree {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.category-node {
  border-radius: var(--radius-md);
  overflow: hidden;
}

.node-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  
  &:hover {
    background: var(--bg-tertiary);
  }
}

.child-categories {
  margin-left: var(--spacing-lg);
  margin-top: var(--spacing-xs);
  
  .category-node.child .node-content {
    background: var(--bg-card);
    border: 1px dashed var(--border-color);
  }
}

.node-name {
  flex: 1;
  font-size: var(--font-sm);
  font-weight: 500;
  color: var(--text-primary);
}

.node-slug {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
  font-family: monospace;
}

.node-actions {
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
  max-width: 500px;
  max-height: 90vh;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
  overflow: hidden;
  display: flex;
  flex-direction: column;
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
  flex: 1;
  padding: var(--spacing-md);
  overflow-y: auto;
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
  
  input, select, textarea {
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
  
  textarea {
    resize: vertical;
  }
}

.seo-section {
  margin-top: var(--spacing-md);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--border-color);
  
  h4 {
    font-size: var(--font-sm);
    margin-bottom: var(--spacing-sm);
    color: var(--text-secondary);
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
