<template>
  <div class="users-manage">
    <div class="page-header">
      <h1>用户管理</h1>
    </div>
    
    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>
    
    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>用户名</th>
            <th>邮箱</th>
            <th>角色</th>
            <th>注册时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td class="user-cell">
              <span class="avatar">{{ user.username?.charAt(0)?.toUpperCase() || '?' }}</span>
              <span class="username">{{ user.username }}</span>
            </td>
            <td>{{ user.email }}</td>
            <td>
              <span class="role-badge" :class="user.role">
                {{ getRoleLabel(user.role) }}
              </span>
            </td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td class="actions">
              <select 
                :value="user.role" 
                @change="changeRole(user, $event.target.value)"
                class="role-select"
              >
                <option value="author">作者</option>
                <option value="editor">编辑</option>
                <option value="admin">管理员</option>
              </select>
              <button 
                class="btn btn-danger btn-sm" 
                @click="deleteUser(user)"
                :disabled="user.username === 'admin'"
              >
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="users.length === 0" class="empty-state">
        暂无用户
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { formatDate } from '@/utils'
import request from '@/utils/request'

const users = ref([])
const loading = ref(false)

function getRoleLabel(role) {
  const labels = {
    admin: '管理员',
    editor: '编辑',
    author: '作者',
    guest: '游客'
  }
  return labels[role] || role
}

async function loadUsers() {
  loading.value = true
  try {
    const res = await request.get('/admin/users')
    users.value = res.data.users || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function changeRole(user, newRole) {
  if (!confirm(`确定要将用户 "${user.username}" 的角色改为 "${getRoleLabel(newRole)}" 吗？`)) {
    loadUsers()
    return
  }
  
  try {
    await request.put(`/admin/users/${user.id}/role`, { role: newRole })
    loadUsers()
  } catch (e) {
    alert('操作失败')
    loadUsers()
  }
}

async function deleteUser(user) {
  if (!confirm(`确定要删除用户 "${user.username}" 吗？`)) return
  
  try {
    await request.delete(`/admin/users/${user.id}`)
    loadUsers()
  } catch (e) {
    alert('删除失败')
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<style lang="scss" scoped>
.users-manage {
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

.user-cell {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.avatar {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--primary-color);
  color: white;
  font-size: var(--font-xs);
  font-weight: 600;
}

.username {
  font-weight: 500;
  color: var(--text-primary);
}

.role-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: var(--radius-full);
  font-size: var(--font-xs);
  font-weight: 500;
  
  &.admin {
    background: #fef3c7;
    color: #92400e;
  }
  
  &.editor {
    background: #dbeafe;
    color: #1e40af;
  }
  
  &.author {
    background: #d1fae5;
    color: #065f46;
  }
}

.theme-dark .role-badge {
  &.admin {
    background: rgba(251, 191, 36, 0.2);
    color: #fbbf24;
  }
  
  &.editor {
    background: rgba(59, 130, 246, 0.2);
    color: #60a5fa;
  }
  
  &.author {
    background: rgba(16, 185, 129, 0.2);
    color: #34d399;
  }
}

.actions {
  display: flex;
  gap: var(--spacing-xs);
  align-items: center;
}

.role-select {
  padding: 4px 8px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  font-size: var(--font-xs);
  background: var(--bg-primary);
  color: var(--text-primary);
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--text-tertiary);
  font-size: var(--font-sm);
}
</style>
