<template>
  <div class="article-editor">
    <div class="editor-header">
      <input 
        v-model="article.title" 
        type="text" 
        class="title-input"
        placeholder="输入文章标题..."
        @blur="generateSlug"
      />
      
      <div class="editor-actions">
        <div class="mode-switch">
          <button 
            :class="{ active: editorMode === 'markdown' }" 
            @click="editorMode = 'markdown'"
            class="mode-btn"
          >
            Markdown
          </button>
          <button 
            :class="{ active: editorMode === 'richtext' }" 
            @click="editorMode = 'richtext'"
            class="mode-btn"
          >
            富文本
          </button>
        </div>
        
        <span v-if="lastSaved" class="save-status">
          {{ lastSavedText }}
        </span>
        
        <button class="btn btn-secondary" @click="showVersions = true">
          📜 版本历史
        </button>
        
        <button class="btn btn-secondary" @click="saveDraft">
          💾 存草稿
        </button>
        
        <button class="btn btn-primary" @click="publishArticle" :disabled="publishing">
          {{ publishing ? '发布中...' : (isEditing ? '更新文章' : '发布文章') }}
        </button>
      </div>
    </div>
    
    <div class="editor-body">
      <div class="main-editor">
        <div v-if="editorMode === 'richtext'" class="richtext-editor">
          <div class="toolbar">
            <button @click="execCommand('bold')" title="加粗"><b>B</b></button>
            <button @click="execCommand('italic')" title="斜体"><i>I</i></button>
            <button @click="execCommand('underline')" title="下划线"><u>U</u></button>
            <span class="divider"></span>
            <button @click="execCommand('insertUnorderedList')" title="无序列表">• 列表</button>
            <button @click="execCommand('insertOrderedList')" title="有序列表">1. 列表</button>
            <span class="divider"></span>
            <button @click="execCommand('formatBlock', 'h1')" title="标题1">H1</button>
            <button @click="execCommand('formatBlock', 'h2')" title="标题2">H2</button>
            <button @click="execCommand('formatBlock', 'h3')" title="标题3">H3</button>
            <span class="divider"></span>
            <button @click="execCommand('justifyLeft')" title="左对齐">⬅</button>
            <button @click="execCommand('justifyCenter')" title="居中">⬌</button>
            <button @click="execCommand('justifyRight')" title="右对齐">➡</button>
            <span class="divider"></span>
            <button @click="insertLink" title="插入链接">🔗</button>
            <button @click="insertImage" title="插入图片">🖼️</button>
            <button @click="insertCodeBlock" title="代码块">&lt;/&gt;</button>
            <button @click="execCommand('formatBlock', 'blockquote')" title="引用">❝</button>
          </div>
          
          <div 
            ref="richTextEditor"
            class="editor-content"
            contenteditable="true"
            @input="onRichTextInput"
            placeholder="开始写作..."
          ></div>
        </div>
        
        <div v-else class="markdown-editor">
          <div class="markdown-toolbar">
            <button @click="insertMarkdown('**', '**')" title="加粗"><b>B</b></button>
            <button @click="insertMarkdown('*', '*')" title="斜体"><i>I</i></button>
            <span class="divider"></span>
            <button @click="insertMarkdown('# ', '')" title="标题">#</button>
            <button @click="insertMarkdown('## ', '')" title="二级标题">##</button>
            <span class="divider"></span>
            <button @click="insertMarkdown('- ', '')" title="无序列表">•</button>
            <button @click="insertMarkdown('1. ', '')" title="有序列表">1.</button>
            <span class="divider"></span>
            <button @click="insertMarkdownLink" title="链接">🔗</button>
            <button @click="insertMarkdownImage" title="图片">🖼️</button>
            <button @click="insertMarkdownCode" title="代码">{'<>'}</button>
            <button @click="insertMarkdown('> ', '')" title="引用">❝</button>
          </div>
          
          <div class="markdown-split">
            <div class="markdown-input">
              <textarea 
                ref="markdownTextarea"
                v-model="markdownContent"
                class="markdown-textarea"
                placeholder="用 Markdown 开始写作...

# 一级标题
## 二级标题

**粗体** *斜体*

- 列表项1
- 列表项2

> 引用文本

`代码`

```
代码块
```

[链接文字](url)
![图片描述](url)"
                @input="onMarkdownInput"
              ></textarea>
            </div>
            
            <div class="markdown-preview">
              <div class="preview-label">预览</div>
              <div class="preview-content" v-html="renderedMarkdown"></div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="sidebar">
        <div class="sidebar-section">
          <h4>封面图</h4>
          <div v-if="article.cover" class="cover-preview">
            <img :src="article.cover" alt="封面" />
            <button class="remove-cover" @click="article.cover = ''">×</button>
          </div>
          <div class="cover-actions">
            <button class="btn btn-secondary btn-sm btn-block" @click="triggerUpload">
              📤 上传封面
            </button>
            <input 
              ref="coverInput" 
              type="file" 
              accept="image/*" 
              style="display: none" 
              @change="handleCoverUpload"
            />
            <button class="btn btn-secondary btn-sm btn-block" @click="showUrlInput = !showUrlInput">
              🔗 图片URL
            </button>
          </div>
          <div v-if="showUrlInput" class="url-input">
            <input 
              v-model="coverUrl" 
              type="text" 
              placeholder="输入图片URL"
              @keyup.enter="applyCoverUrl"
            />
            <button class="btn btn-primary btn-sm" @click="applyCoverUrl">确定</button>
          </div>
        </div>
        
        <div class="sidebar-section">
          <h4>Slug</h4>
          <input 
            v-model="article.slug" 
            type="text" 
            placeholder="自动生成"
            @blur="validateSlug"
          />
          <p class="hint">文章URL的友好标识</p>
        </div>
        
        <div class="sidebar-section">
          <h4>分类</h4>
          <select v-model="article.category_id" class="category-select">
            <option value="">请选择分类</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
        </div>
        
        <div class="sidebar-section">
          <h4>标签</h4>
          <div class="tags-input">
            <span v-for="(tag, index) in article.tags" :key="tag" class="tag">
              {{ tag }}
              <button class="remove-tag" @click="removeTag(index)">×</button>
            </span>
            <input 
              v-model="newTag" 
              type="text" 
              placeholder="添加标签，回车确认"
              @keyup.enter="addTag"
              class="tag-input"
            />
          </div>
        </div>
        
        <div class="sidebar-section">
          <h4>摘要</h4>
          <textarea 
            v-model="article.summary" 
            rows="4"
            placeholder="自动从正文截取200字，可手动修改"
          ></textarea>
        </div>
        
        <div class="sidebar-section">
          <h4>发布状态</h4>
          <select v-model="article.status" class="status-select">
            <option value="draft">草稿</option>
            <option value="published">立即发布</option>
            <option value="scheduled">定时发布</option>
          </select>
          
          <div v-if="article.status === 'scheduled'" class="schedule-input">
            <label>发布时间</label>
            <input 
              v-model="publishTime" 
              type="datetime-local"
            />
          </div>
        </div>
        
        <div class="sidebar-section">
          <h4>阅读权限</h4>
          <select v-model="article.permissions" class="perm-select">
            <option value="public">公开</option>
            <option value="login">登录可见</option>
            <option value="role">指定角色可见</option>
            <option value="closed">关闭评论</option>
          </select>
        </div>
      </div>
    </div>
    
    <div v-if="showVersions" class="version-modal">
      <div class="modal-backdrop" @click="showVersions = false"></div>
      <div class="modal-content">
        <div class="modal-header">
          <h3>版本历史</h3>
          <button class="close-btn" @click="showVersions = false">×</button>
        </div>
        <div class="modal-body">
          <div v-if="versions.length === 0" class="empty-versions">
            暂无历史版本
          </div>
          <div v-else class="version-list">
            <div 
              v-for="(version, index) in versions.slice().reverse()" 
              :key="version.id" 
              class="version-item"
            >
              <div class="version-info">
                <span class="version-num">v{{ version.version }}</span>
                <span class="version-time">{{ formatDate(version.created_at) }}</span>
              </div>
              <div class="version-title">{{ version.title }}</div>
              <div class="version-actions">
                <button class="btn btn-secondary btn-sm" @click="previewVersion(version)">
                  预览
                </button>
                <button class="btn btn-primary btn-sm" @click="revertVersion(version)">
                  回滚到此版本
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked } from 'marked'
import hljs from 'highlight.js'
import { saveDraft, getDraft, deleteDraft } from '@/utils/draftDB'
import { formatDate, slugify, truncateText, extractFirstParagraph } from '@/utils'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()

const editorMode = ref('markdown')
const isEditing = ref(false)
const publishing = ref(false)
const lastSaved = ref(null)
const showVersions = ref(false)
const versions = ref([])
const categories = ref([])
const markdownContent = ref('')

const richTextEditor = ref(null)
const markdownTextarea = ref(null)
const coverInput = ref(null)

const showUrlInput = ref(false)
const coverUrl = ref('')
const newTag = ref('')
const publishTime = ref('')

const article = ref({
  id: '',
  title: '',
  slug: '',
  content: '',
  summary: '',
  cover: '',
  category_id: '',
  tags: [],
  status: 'draft',
  permissions: 'public',
  publish_at: null
})

const lastSavedText = computed(() => {
  if (!lastSaved.value) return ''
  const diff = Math.floor((Date.now() - new Date(lastSaved.value).getTime()) / 1000)
  if (diff < 60) return '刚刚保存'
  if (diff < 3600) return `${Math.floor(diff / 60)} 分钟前保存`
  return formatDate(lastSaved.value)
})

const renderedMarkdown = computed(() => {
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
    return marked.parse(markdownContent.value)
  } catch (e) {
    return markdownContent.value
  }
})

function generateSlug() {
  if (article.value.title && !article.value.slug) {
    article.value.slug = slugify(article.value.title)
  }
}

function validateSlug() {
  if (article.value.slug) {
    article.value.slug = slugify(article.value.slug)
  }
}

function onRichTextInput() {
  if (richTextEditor.value) {
    article.value.content = richTextEditor.value.innerHTML
  }
  autoSave()
}

function onMarkdownInput() {
  article.value.content = markdownContent.value
  autoSave()
}

let autoSaveTimer = null
function autoSave() {
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(() => {
    saveToLocal()
  }, 2000)
}

let intervalTimer = null

function saveToLocal() {
  const draftId = isEditing.value ? article.value.id : 'new'
  saveDraft(draftId, article.value)
  lastSaved.value = new Date().toISOString()
}

async function loadDraft() {
  const draftId = isEditing.value ? article.value.id : 'new'
  const draft = await getDraft(draftId)
  if (draft && draft.content) {
    if (confirm('检测到未保存的草稿，是否恢复？')) {
      article.value = { ...article.value, ...draft }
      if (editorMode.value === 'markdown') {
        markdownContent.value = article.value.content
      } else if (richTextEditor.value) {
        richTextEditor.value.innerHTML = article.value.content
      }
    } else {
      await deleteDraft(draftId)
    }
  }
}

function saveDraftManual() {
  saveToLocal()
  alert('草稿已保存到本地')
}

async function publishArticle() {
  if (!article.value.title) {
    alert('请输入文章标题')
    return
  }
  if (!article.value.content) {
    alert('请输入文章内容')
    return
  }
  
  if (!article.value.summary) {
    article.value.summary = truncateText(extractFirstParagraph(article.value.content), 200)
  }
  
  if (article.value.status === 'scheduled' && publishTime.value) {
    article.value.publish_at = new Date(publishTime.value).toISOString()
  }
  
  publishing.value = true
  try {
    if (isEditing.value) {
      await request.put(`/articles/${article.value.id}`, article.value)
      alert('文章更新成功')
    } else {
      const res = await request.post('/articles', article.value)
      article.value.id = res.data.article.id
      isEditing.value = true
      alert('文章发布成功')
    }
    
    const draftId = isEditing.value ? article.value.id : 'new'
    await deleteDraft(draftId)
    
    router.push(`/article/${article.value.slug || article.value.id}`)
  } catch (e) {
    alert(e.response?.data?.error || '保存失败')
  } finally {
    publishing.value = false
  }
}

async function loadArticle() {
  const id = route.params.id
  if (!id) return
  
  try {
    const res = await request.get(`/articles/${id}`)
    article.value = res.data.article
    isEditing.value = true
    
    if (editorMode.value === 'markdown') {
      markdownContent.value = article.value.content
    } else if (richTextEditor.value) {
      richTextEditor.value.innerHTML = article.value.content
    }
    
    loadVersions()
  } catch (e) {
    console.error(e)
  }
}

async function loadVersions() {
  if (!article.value.id) return
  
  try {
    const res = await request.get(`/articles/${article.value.id}/versions`)
    versions.value = res.data.versions || []
  } catch (e) {
    console.error(e)
  }
}

function previewVersion(version) {
  alert(`版本 v${version.version}\n标题: ${version.title}`)
}

async function revertVersion(version) {
  if (!confirm(`确定要回滚到版本 v${version.version} 吗？`)) return
  
  try {
    await request.post(`/articles/${article.value.id}/versions/${version.id}/revert`)
    
    const res = await request.get(`/articles/${article.value.id}`)
    article.value = res.data.article
    
    if (editorMode.value === 'markdown') {
      markdownContent.value = article.value.content
    } else if (richTextEditor.value) {
      richTextEditor.value.innerHTML = article.value.content
    }
    
    loadVersions()
    showVersions.value = false
    alert('回滚成功')
  } catch (e) {
    alert('回滚失败')
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

function addTag() {
  const tag = newTag.value.trim()
  if (tag && !article.value.tags.includes(tag)) {
    article.value.tags.push(tag)
  }
  newTag.value = ''
}

function removeTag(index) {
  article.value.tags.splice(index, 1)
}

function triggerUpload() {
  coverInput.value?.click()
}

async function handleCoverUpload(e) {
  const file = e.target.files?.[0]
  if (!file) return
  
  const formData = new FormData()
  formData.append('file', file)
  
  try {
    const res = await request.post('/upload/cover', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    article.value.cover = res.data.url
  } catch (e) {
    alert('上传失败')
  }
}

function applyCoverUrl() {
  if (coverUrl.value) {
    article.value.cover = coverUrl.value
    showUrlInput.value = false
    coverUrl.value = ''
  }
}

function execCommand(command, value = null) {
  document.execCommand(command, false, value)
  richTextEditor.value?.focus()
  onRichTextInput()
}

function insertLink() {
  const url = prompt('请输入链接地址:', 'https://')
  if (url) {
    execCommand('createLink', url)
  }
}

function insertImage() {
  const url = prompt('请输入图片地址:', 'https://')
  if (url) {
    execCommand('insertImage', url)
  }
}

function insertCodeBlock() {
  const code = prompt('请输入代码:')
  if (code) {
    execCommand('insertHTML', `<pre><code>${code}</code></pre>`)
  }
}

function insertMarkdown(prefix, suffix) {
  const textarea = markdownTextarea.value
  if (!textarea) return
  
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = textarea.value.substring(start, end)
  
  const before = textarea.value.substring(0, start)
  const after = textarea.value.substring(end)
  
  const newText = before + prefix + selectedText + suffix + after
  markdownContent.value = newText
  
  nextTick(() => {
    textarea.focus()
    const newPos = start + prefix.length + selectedText.length + suffix.length
    textarea.setSelectionRange(newPos, newPos)
  })
}

function insertMarkdownLink() {
  const url = prompt('请输入链接地址:', 'https://')
  if (url) {
    insertMarkdown('[', `](${url})`)
  }
}

function insertMarkdownImage() {
  const url = prompt('请输入图片地址:', 'https://')
  if (url) {
    insertMarkdown('![图片描述', `](${url})`)
  }
}

function insertMarkdownCode() {
  insertMarkdown('`', '`')
}

watch(editorMode, (newMode) => {
  if (newMode === 'markdown') {
    markdownContent.value = article.value.content
  } else if (newMode === 'richtext' && richTextEditor.value) {
    richTextEditor.value.innerHTML = article.value.content
  }
})

onMounted(async () => {
  await loadCategories()
  
  if (route.params.id) {
    await loadArticle()
  }
  
  loadDraft()
  
  intervalTimer = setInterval(() => {
    saveToLocal()
  }, 30000)
})

onUnmounted(() => {
  if (intervalTimer) {
    clearInterval(intervalTimer)
  }
  if (autoSaveTimer) {
    clearTimeout(autoSaveTimer)
  }
})
</script>

<style lang="scss" scoped>
.article-editor {
  min-height: calc(100vh - var(--header-height));
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-card);
  gap: var(--spacing-md);
  flex-wrap: wrap;
  
  @media (max-width: 768px) {
    flex-direction: column;
    align-items: stretch;
  }
}

.title-input {
  flex: 1;
  min-width: 300px;
  border: none;
  outline: none;
  font-size: var(--font-xl);
  font-weight: 600;
  background: transparent;
  color: var(--text-primary);
  
  &::placeholder {
    color: var(--text-tertiary);
  }
}

.editor-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  flex-wrap: wrap;
}

.mode-switch {
  display: flex;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  padding: 2px;
}

.mode-btn {
  padding: var(--spacing-xs) var(--spacing-md);
  border-radius: var(--radius-sm);
  font-size: var(--font-xs);
  color: var(--text-secondary);
  transition: all var(--transition-fast);
  
  &.active {
    background: var(--primary-color);
    color: white;
  }
}

.save-status {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.editor-body {
  flex: 1;
  display: flex;
  overflow: hidden;
  
  @media (max-width: 1024px) {
    flex-direction: column;
  }
}

.main-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.richtext-editor {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  padding: var(--spacing-sm) var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-card);
  
  button {
    padding: 4px 8px;
    border-radius: var(--radius-sm);
    font-size: var(--font-sm);
    color: var(--text-secondary);
    
    &:hover {
      background: var(--bg-secondary);
      color: var(--text-primary);
    }
  }
  
  .divider {
    width: 1px;
    background: var(--border-color);
    margin: 0 var(--spacing-xs);
  }
}

.editor-content {
  flex: 1;
  padding: var(--spacing-lg);
  overflow-y: auto;
  outline: none;
  line-height: 1.8;
  font-size: var(--font-base);
  color: var(--text-primary);
  
  &:empty::before {
    content: attr(placeholder);
    color: var(--text-tertiary);
  }
}

.markdown-editor {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.markdown-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  padding: var(--spacing-sm) var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-card);
  
  button {
    padding: 4px 8px;
    border-radius: var(--radius-sm);
    font-size: var(--font-sm);
    color: var(--text-secondary);
    
    &:hover {
      background: var(--bg-secondary);
      color: var(--text-primary);
    }
  }
  
  .divider {
    width: 1px;
    background: var(--border-color);
    margin: 0 var(--spacing-xs);
  }
}

.markdown-split {
  flex: 1;
  display: flex;
  overflow: hidden;
  
  @media (max-width: 768px) {
    flex-direction: column;
  }
}

.markdown-input {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--border-color);
  
  @media (max-width: 768px) {
    border-right: none;
    border-bottom: 1px solid var(--border-color);
  }
}

.markdown-textarea {
  flex: 1;
  width: 100%;
  padding: var(--spacing-lg);
  border: none;
  outline: none;
  resize: none;
  font-family: 'Fira Code', Monaco, Consolas, monospace;
  font-size: var(--font-sm);
  line-height: 1.6;
  background: var(--bg-primary);
  color: var(--text-primary);
  
  &::placeholder {
    color: var(--text-tertiary);
  }
}

.markdown-preview {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-secondary);
}

.preview-label {
  padding: var(--spacing-xs) var(--spacing-md);
  font-size: var(--font-xs);
  color: var(--text-tertiary);
  background: var(--bg-tertiary);
}

.preview-content {
  flex: 1;
  padding: var(--spacing-md);
  overflow-y: auto;
  font-size: var(--font-sm);
  line-height: 1.6;
  color: var(--text-primary);
  
  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    margin-top: var(--spacing-md);
    margin-bottom: var(--spacing-sm);
    font-weight: 600;
  }
  
  :deep(h1) { font-size: var(--font-xl); }
  :deep(h2) { font-size: var(--font-lg); }
  :deep(h3) { font-size: var(--font-base); }
  
  :deep(p) {
    margin-bottom: var(--spacing-sm);
  }
  
  :deep(pre) {
    background: var(--bg-primary);
    padding: var(--spacing-sm);
    border-radius: var(--radius-sm);
    overflow-x: auto;
    margin-bottom: var(--spacing-sm);
    
    code {
      font-size: var(--font-xs);
    }
  }
  
  :deep(code) {
    background: var(--bg-primary);
    padding: 1px 4px;
    border-radius: 3px;
    font-size: 0.9em;
  }
  
  :deep(blockquote) {
    border-left: 3px solid var(--primary-color);
    padding-left: var(--spacing-sm);
    margin: var(--spacing-sm) 0;
    color: var(--text-secondary);
  }
  
  :deep(ul), :deep(ol) {
    margin-bottom: var(--spacing-sm);
    padding-left: var(--spacing-md);
  }
  
  :deep(img) {
    max-width: 100%;
    border-radius: var(--radius-sm);
  }
}

.sidebar {
  width: 280px;
  border-left: 1px solid var(--border-color);
  background: var(--bg-card);
  padding: var(--spacing-md);
  overflow-y: auto;
  
  @media (max-width: 1024px) {
    width: 100%;
    border-left: none;
    border-top: 1px solid var(--border-color);
    max-height: 300px;
  }
}

.sidebar-section {
  margin-bottom: var(--spacing-lg);
  
  h4 {
    font-size: var(--font-sm);
    font-weight: 600;
    margin-bottom: var(--spacing-sm);
    color: var(--text-secondary);
  }
}

.cover-preview {
  position: relative;
  margin-bottom: var(--spacing-sm);
  
  img {
    width: 100%;
    height: 120px;
    object-fit: cover;
    border-radius: var(--radius-md);
  }
  
  .remove-cover {
    position: absolute;
    top: 4px;
    right: 4px;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.5);
    color: white;
    font-size: var(--font-sm);
  }
}

.cover-actions {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.btn-block {
  display: block;
  width: 100%;
}

.url-input {
  margin-top: var(--spacing-sm);
  display: flex;
  gap: var(--spacing-xs);
  
  input {
    flex: 1;
    padding: var(--spacing-xs) var(--spacing-sm);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    font-size: var(--font-xs);
    background: var(--bg-primary);
    color: var(--text-primary);
  }
}

.sidebar-section input[type="text"],
.sidebar-section textarea,
.sidebar-section select {
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

.sidebar-section textarea {
  resize: vertical;
  min-height: 60px;
}

.hint {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
  margin-top: 4px;
}

.tags-input {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  background: var(--bg-primary);
}

.tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  background: var(--primary-color);
  color: white;
  font-size: var(--font-xs);
  
  .remove-tag {
    color: white;
    opacity: 0.7;
    font-size: var(--font-sm);
    line-height: 1;
    
    &:hover {
      opacity: 1;
    }
  }
}

.tag-input {
  flex: 1;
  min-width: 80px;
  border: none;
  outline: none;
  font-size: var(--font-xs);
  background: transparent;
  color: var(--text-primary);
}

.schedule-input {
  margin-top: var(--spacing-sm);
  
  label {
    display: block;
    font-size: var(--font-xs);
    color: var(--text-secondary);
    margin-bottom: 4px;
  }
  
  input {
    width: 100%;
    padding: var(--spacing-xs) var(--spacing-sm);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    font-size: var(--font-xs);
    background: var(--bg-primary);
    color: var(--text-primary);
  }
}

.version-modal {
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
  max-height: 80vh;
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
    font-size: var(--font-lg);
    color: var(--text-secondary);
    
    &:hover {
      background: var(--bg-secondary);
    }
  }
}

.modal-body {
  padding: var(--spacing-md);
  max-height: 60vh;
  overflow-y: auto;
}

.empty-versions {
  text-align: center;
  padding: var(--spacing-xl) 0;
  color: var(--text-tertiary);
  font-size: var(--font-sm);
}

.version-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.version-item {
  padding: var(--spacing-sm);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
}

.version-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-xs);
}

.version-num {
  font-weight: 600;
  font-size: var(--font-sm);
  color: var(--primary-color);
}

.version-time {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.version-title {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
}

.version-actions {
  display: flex;
  gap: var(--spacing-xs);
  justify-content: flex-end;
}
</style>
