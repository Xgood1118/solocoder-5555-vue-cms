import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import request from '@/utils/request'

export const useArticleStore = defineStore('article', () => {
  const articles = ref([])
  const currentArticle = ref(null)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const loading = ref(false)
  const categories = ref([])
  const tags = ref([])
  const hotArticles = ref([])
  const stats = ref(null)

  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  async function fetchArticles(params = {}) {
    loading.value = true
    try {
      const res = await request.get('/articles', {
        params: {
          page: page.value,
          pageSize: pageSize.value,
          ...params
        }
      })
      articles.value = res.data.list || res.data.articles || []
      total.value = res.data.total || 0
      return res.data
    } finally {
      loading.value = false
    }
  }

  async function fetchArticle(id) {
    const res = await request.get(`/articles/${id}`)
    currentArticle.value = res.data
    return res.data
  }

  async function createArticle(data) {
    const res = await request.post('/articles', data)
    return res.data
  }

  async function updateArticle(id, data) {
    const res = await request.put(`/articles/${id}`, data)
    return res.data
  }

  async function deleteArticle(id) {
    const res = await request.delete(`/articles/${id}`)
    return res.data
  }

  async function fetchCategories() {
    const res = await request.get('/categories')
    categories.value = res.data.list || res.data || []
    return categories.value
  }

  async function fetchTags() {
    const res = await request.get('/tags')
    tags.value = res.data.list || res.data || []
    return tags.value
  }

  async function fetchHotArticles(limit = 10) {
    const res = await request.get('/articles/hot', { params: { limit } })
    hotArticles.value = res.data.list || res.data || []
    return hotArticles.value
  }

  async function fetchStats() {
    const res = await request.get('/stats')
    stats.value = res.data
    return res.data
  }

  async function searchArticles(keyword, params = {}) {
    loading.value = true
    try {
      const res = await request.get('/articles/search', {
        params: { keyword, ...params }
      })
      articles.value = res.data.list || res.data.articles || []
      total.value = res.data.total || 0
      return res.data
    } finally {
      loading.value = false
    }
  }

  function setPage(p) {
    page.value = p
  }

  function resetPagination() {
    page.value = 1
    total.value = 0
    articles.value = []
  }

  return {
    articles,
    currentArticle,
    total,
    page,
    pageSize,
    totalPages,
    loading,
    categories,
    tags,
    hotArticles,
    stats,
    fetchArticles,
    fetchArticle,
    createArticle,
    updateArticle,
    deleteArticle,
    fetchCategories,
    fetchTags,
    fetchHotArticles,
    fetchStats,
    searchArticles,
    setPage,
    resetPagination
  }
})
