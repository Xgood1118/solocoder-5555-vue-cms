<template>
  <div class="tags-page">
    <div class="container">
      <div class="page-header">
        <h1>标签云</h1>
        <p>按标签浏览文章</p>
      </div>
      
      <div v-if="loading" class="loading-wrap">
        <div class="loading"></div>
      </div>
      
      <div v-else class="tag-cloud">
        <router-link 
          v-for="tag in tags" 
          :key="tag.id"
          :to="`/tag/${tag.slug}`"
          class="tag-item"
          :style="{ fontSize: getTagSize(tag.count) + 'px' }"
        >
          {{ tag.name }}
          <span class="count">({{ tag.count }})</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const tags = ref([])
const loading = ref(false)

function getTagSize(count) {
  const baseSize = 14
  const maxSize = 32
  const maxCount = Math.max(...tags.value.map(t => t.count || 0), 1)
  return baseSize + (count / maxCount) * (maxSize - baseSize)
}

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

onMounted(() => {
  loadTags()
})
</script>

<style lang="scss" scoped>
.tags-page {
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

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-xl);
  background: var(--bg-card);
  border-radius: var(--radius-xl);
  border: 1px solid var(--border-color);
  min-height: 200px;
}

.tag-item {
  display: inline-flex;
  align-items: baseline;
  gap: 4px;
  padding: 6px 16px;
  border-radius: var(--radius-full);
  background: var(--bg-secondary);
  color: var(--text-secondary);
  transition: all var(--transition-fast);
  
  &:hover {
    background: var(--primary-color);
    color: white;
    transform: scale(1.05);
  }
  
  .count {
    font-size: 0.75em;
    opacity: 0.7;
  }
}
</style>
