<template>
  <div class="about-page">
    <div class="container">
      <div class="about-card">
        <h1>关于轻量 CMS</h1>
        
        <section class="section">
          <h2>项目简介</h2>
          <p>
            轻量 CMS 是一个基于 <strong>Vue 3 + Go</strong> 构建的现代化内容管理系统，
            采用前后端分离架构，数据全部存储在 JSON 文件中，无需任何外部数据库依赖。
          </p>
          <p>
            专为小团队博客、文档站、技术博客等内容站点设计，
            只需一个命令即可启动，真正做到零运维部署。
          </p>
        </section>
        
        <section class="section">
          <h2>核心特性</h2>
          <div class="features-grid">
            <div class="feature">
              <span class="feature-icon">✏️</span>
              <h3>双编辑器</h3>
              <p>同时支持富文本和 Markdown 双模式编辑，满足不同写作习惯</p>
            </div>
            <div class="feature">
              <span class="feature-icon">💾</span>
              <h3>自动存稿</h3>
              <p>每30秒自动保存草稿到 IndexedDB，写文章再也不怕丢失</p>
            </div>
            <div class="feature">
              <span class="feature-icon">📜</span>
              <h3>版本管理</h3>
              <p>自动生成多个历史版本，支持任意版本回滚</p>
            </div>
            <div class="feature">
              <span class="feature-icon">⏰</span>
              <h3>定时发布</h3>
              <p>支持设定未来时间点自动发布，让内容发布更从容</p>
            </div>
            <div class="feature">
              <span class="feature-icon">🌙</span>
              <h3>暗色主题</h3>
              <p>完美支持明暗主题切换，护眼模式随时切换</p>
            </div>
            <div class="feature">
              <span class="feature-icon">📱</span>
              <h3>响应式设计</h3>
              <p>完美适配桌面和移动端，随时随地浏览内容</p>
            </div>
            <div class="feature">
              <span class="feature-icon">🔍</span>
              <h3>全文搜索</h3>
              <p>支持标题、正文、标签、作者四个维度搜索，关键词高亮显示</p>
            </div>
            <div class="feature">
              <span class="feature-icon">💬</span>
              <h3>评论系统</h3>
              <p>支持嵌套评论、垃圾过滤、审核机制，打造优质交流环境</p>
            </div>
          </div>
        </section>
        
        <section class="section">
          <h2>技术栈</h2>
          <div class="tech-stack">
            <div class="tech-column">
              <h3>前端</h3>
              <ul>
                <li>Vue 3 (Composition API)</li>
                <li>Vite</li>
                <li>Pinia</li>
                <li>Vue Router</li>
                <li>TipTap 编辑器</li>
                <li>Markdown (marked + highlight.js)</li>
                <li>IndexedDB (idb)</li>
                <li>PWA 支持</li>
              </ul>
            </div>
            <div class="tech-column">
              <h3>后端</h3>
              <ul>
                <li>Go (Gin 框架)</li>
                <li>JWT 认证</li>
                <li>JSON 文件存储</li>
                <li>Goldmark Markdown 渲染</li>
                <li>Imaging 图片处理</li>
                <li>Cron 定时任务</li>
                <li>CSRF 防护</li>
                <li>限流防刷</li>
              </ul>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
</script>

<style lang="scss" scoped>
.about-page {
  padding: var(--spacing-2xl) 0;
  min-height: calc(100vh - var(--header-height) - 200px);
}

.about-card {
  max-width: 900px;
  margin: 0 auto;
  background: var(--bg-card);
  border-radius: var(--radius-xl);
  padding: var(--spacing-2xl);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
  
  h1 {
    font-size: var(--font-3xl);
    text-align: center;
    margin-bottom: var(--spacing-2xl);
    color: var(--primary-color);
  }
}

.section {
  margin-bottom: var(--spacing-2xl);
  
  h2 {
    font-size: var(--font-xl);
    margin-bottom: var(--spacing-md);
    padding-bottom: var(--spacing-sm);
    border-bottom: 2px solid var(--primary-color);
    display: inline-block;
  }
  
  p {
    color: var(--text-secondary);
    line-height: 1.8;
    margin-bottom: var(--spacing-sm);
  }
  
  strong {
    color: var(--text-primary);
  }
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: var(--spacing-md);
  margin-top: var(--spacing-md);
}

.feature {
  padding: var(--spacing-md);
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
  text-align: center;
  transition: all var(--transition-fast);
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
  }
}

.feature-icon {
  font-size: 2.5rem;
  margin-bottom: var(--spacing-sm);
  display: block;
}

.feature h3 {
  font-size: var(--font-base);
  margin-bottom: var(--spacing-xs);
  color: var(--text-primary);
}

.feature p {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  line-height: 1.5;
}

.tech-stack {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-xl);
  
  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }
}

.tech-column {
  h3 {
    font-size: var(--font-base);
    margin-bottom: var(--spacing-sm);
    color: var(--primary-color);
  }
  
  ul {
    list-style: none;
    
    li {
      padding: var(--spacing-xs) 0;
      padding-left: var(--spacing-md);
      position: relative;
      font-size: var(--font-sm);
      color: var(--text-secondary);
      
      &::before {
        content: '✓';
        position: absolute;
        left: 0;
        color: var(--success-color);
        font-weight: bold;
      }
    }
  }
}
</style>
