import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/article/:id',
    name: 'ArticleDetail',
    component: () => import('@/views/Article/Detail.vue'),
    meta: { title: '文章详情' }
  },
  {
    path: '/article/editor',
    name: 'ArticleEditor',
    component: () => import('@/views/Article/Editor.vue'),
    meta: { title: '写文章', requiresAuth: true }
  },
  {
    path: '/article/editor/:id',
    name: 'ArticleEditorEdit',
    component: () => import('@/views/Article/Editor.vue'),
    meta: { title: '编辑文章', requiresAuth: true }
  },
  {
    path: '/category/:id',
    name: 'CategoryList',
    component: () => import('@/views/Category/List.vue'),
    meta: { title: '分类文章' }
  },
  {
    path: '/tag/:id',
    name: 'TagDetail',
    component: () => import('@/views/Tag/Detail.vue'),
    meta: { title: '标签文章' }
  },
  {
    path: '/search',
    name: 'SearchResults',
    component: () => import('@/views/Search/Results.vue'),
    meta: { title: '搜索结果' }
  },
  {
    path: '/categories',
    name: 'CategoryListAll',
    component: () => import('@/views/Category/All.vue'),
    meta: { title: '分类目录' }
  },
  {
    path: '/tags',
    name: 'TagCloud',
    component: () => import('@/views/Tag/Cloud.vue'),
    meta: { title: '标签云' }
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('@/views/About.vue'),
    meta: { title: '关于' }
  },
  {
    path: '/dashboard',
    name: 'DashboardIndex',
    component: () => import('@/views/Dashboard/Index.vue'),
    meta: { title: '仪表板', requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/dashboard/articles',
    name: 'DashboardArticles',
    component: () => import('@/views/Dashboard/Articles.vue'),
    meta: { title: '文章管理', requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/dashboard/categories',
    name: 'DashboardCategories',
    component: () => import('@/views/Dashboard/Categories.vue'),
    meta: { title: '分类管理', requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/dashboard/comments',
    name: 'DashboardComments',
    component: () => import('@/views/Dashboard/Comments.vue'),
    meta: { title: '评论审核', requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/dashboard/tags',
    name: 'DashboardTags',
    component: () => import('@/views/Dashboard/Tags.vue'),
    meta: { title: '标签管理', requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/dashboard/users',
    name: 'DashboardUsers',
    component: () => import('@/views/Dashboard/Users.vue'),
    meta: { title: '用户管理', requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/dashboard/audit-logs',
    name: 'DashboardAuditLogs',
    component: () => import('@/views/Dashboard/AuditLogs.vue'),
    meta: { title: '审计日志', requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/User/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/User/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/User/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.meta.title) {
    document.title = `${to.meta.title} - 轻量 CMS`
  }

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next({ name: 'Home' })
    return
  }

  next()
})

export default router
