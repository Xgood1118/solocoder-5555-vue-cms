<template>
  <div class="audit-logs">
    <div class="page-header">
      <h1>审计日志</h1>
    </div>
    
    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>
    
    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>时间</th>
            <th>用户</th>
            <th>操作</th>
            <th>资源类型</th>
            <th>资源ID</th>
            <th>IP</th>
            <th>详情</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in logs" :key="log.id">
            <td>{{ formatDate(log.created_at) }}</td>
            <td>{{ log.username || '系统' }}</td>
            <td>
              <span class="action-badge" :class="getActionClass(log.action)">
                {{ log.action }}
              </span>
            </td>
            <td>{{ log.resource }}</td>
            <td class="resource-id">{{ log.resource_id || '-' }}</td>
            <td>{{ log.ip }}</td>
            <td class="detail">{{ log.detail || '-' }}</td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="logs.length === 0" class="empty-state">
        暂无日志记录
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { formatDate } from '@/utils'
import request from '@/utils/request'

const logs = ref([])
const loading = ref(false)

function getActionClass(action) {
  if (action.includes('create') || action.includes('register')) return 'create'
  if (action.includes('update') || action.includes('edit')) return 'update'
  if (action.includes('delete') || action.includes('remove')) return 'delete'
  if (action.includes('approve')) return 'approve'
  if (action.includes('login')) return 'login'
  return 'default'
}

async function loadLogs() {
  loading.value = true
  try {
    const res = await request.get('/admin/audit-logs')
    logs.value = res.data.logs || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadLogs()
})
</script>

<style lang="scss" scoped>
.audit-logs {
  padding: var(--spacing-lg);
}

.page-header {
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

.table-container {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  overflow: hidden;
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  min-width: 800px;
  
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
    white-space: nowrap;
  }
  
  tbody tr:hover {
    background: var(--bg-secondary);
  }
}

.action-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  font-size: var(--font-xs);
  font-weight: 500;
  
  &.create {
    background: #d1fae5;
    color: #065f46;
  }
  
  &.update {
    background: #dbeafe;
    color: #1e40af;
  }
  
  &.delete {
    background: #fecaca;
    color: #991b1b;
  }
  
  &.approve {
    background: #fef3c7;
    color: #92400e;
  }
  
  &.login {
    background: #e9d5ff;
    color: #6b21a8;
  }
  
  &.default {
    background: var(--bg-secondary);
    color: var(--text-secondary);
  }
}

.theme-dark .action-badge {
  &.create {
    background: rgba(16, 185, 129, 0.2);
    color: #34d399;
  }
  
  &.update {
    background: rgba(59, 130, 246, 0.2);
    color: #60a5fa;
  }
  
  &.delete {
    background: rgba(239, 68, 68, 0.2);
    color: #f87171;
  }
  
  &.approve {
    background: rgba(251, 191, 36, 0.2);
    color: #fbbf24;
  }
  
  &.login {
    background: rgba(168, 85, 247, 0.2);
    color: #c084fc;
  }
}

.resource-id {
  font-family: monospace;
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.detail {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--text-tertiary);
  font-size: var(--font-sm);
}
</style>
