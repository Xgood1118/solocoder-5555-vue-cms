import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(false)

  const theme = computed(() => (isDark.value ? 'dark' : 'light'))

  function initTheme() {
    const saved = localStorage.getItem('theme')
    if (saved) {
      isDark.value = saved === 'dark'
    } else {
      isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    applyTheme()
  }

  function toggleTheme() {
    isDark.value = !isDark.value
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
    applyTheme()
  }

  function setTheme(dark) {
    isDark.value = dark
    localStorage.setItem('theme', dark ? 'dark' : 'light')
    applyTheme()
  }

  function applyTheme() {
    const html = document.documentElement
    if (isDark.value) {
      html.classList.add('theme-dark')
      html.classList.remove('theme-light')
    } else {
      html.classList.add('theme-light')
      html.classList.remove('theme-dark')
    }
  }

  watch(isDark, () => {
    const meta = document.querySelector('meta[name="theme-color"]')
    if (meta) {
      meta.setAttribute('content', isDark.value ? '#1f2937' : '#4f46e5')
    }
  })

  return {
    isDark,
    theme,
    initTheme,
    toggleTheme,
    setTheme
  }
})
