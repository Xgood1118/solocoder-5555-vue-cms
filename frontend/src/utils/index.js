import { pinyin } from 'pinyin-pro'

export function slugify(title) {
  if (!title) return ''
  
  const py = pinyin(title, { toneType: 'none', type: 'array' })
  let slug = py.join('-')
  
  slug = slug.toLowerCase()
  slug = slug.replace(/\s+/g, '-')
  slug = slug.replace(/[^a-z0-9-]/g, '')
  slug = slug.replace(/-+/g, '-')
  slug = slug.replace(/^-+|-+$/g, '')
  
  return slug || 'article'
}

export function truncateText(text, maxLen = 200) {
  if (!text) return ''
  
  let clean = text.replace(/\n/g, ' ').replace(/\s+/g, ' ').trim()
  
  if (clean.length <= maxLen) return clean
  
  return clean.substring(0, maxLen) + '...'
}

export function extractFirstParagraph(content) {
  if (!content) return ''
  
  if (content.startsWith('<')) {
    const match = content.match(/<p[^>]*>(.*?)<\/p>/i)
    if (match) {
      return match[1].replace(/<[^>]*>/g, '')
    }
    return content.replace(/<[^>]*>/g, '').trim()
  }
  
  const lines = content.split('\n')
  for (const line of lines) {
    const trimmed = line.trim()
    if (trimmed && !trimmed.startsWith('#') && !trimmed.startsWith('```')) {
      return trimmed
    }
  }
  
  return content
}

export function formatDate(date, format = 'YYYY-MM-DD HH:mm') {
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
}

export function highlightText(text, keyword) {
  if (!keyword || !text) return text
  
  const regex = new RegExp(`(${escapeRegExp(keyword)})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}

function escapeRegExp(string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

export function copyToClipboard(text) {
  if (navigator.clipboard) {
    return navigator.clipboard.writeText(text)
  }
  
  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  document.body.appendChild(textarea)
  textarea.select()
  document.execCommand('copy')
  document.body.removeChild(textarea)
  
  return Promise.resolve()
}

export function debounce(fn, delay = 300) {
  let timer = null
  return function(...args) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => fn.apply(this, args), delay)
  }
}

export function throttle(fn, delay = 300) {
  let last = 0
  return function(...args) {
    const now = Date.now()
    if (now - last >= delay) {
      last = now
      fn.apply(this, args)
    }
  }
}
