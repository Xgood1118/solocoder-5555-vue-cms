import { openDB } from 'idb'

const DB_NAME = 'cms-drafts'
const DB_VERSION = 1
const STORE_NAME = 'drafts'

let db = null

async function getDB() {
  if (db) return db
  
  db = await openDB(DB_NAME, DB_VERSION, {
    upgrade(db) {
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        const store = db.createObjectStore(STORE_NAME, { keyPath: 'id' })
        store.createIndex('updatedAt', 'updatedAt')
      }
    }
  })
  
  return db
}

export async function saveDraft(articleId, data) {
  const database = await getDB()
  const draft = {
    id: articleId || 'new',
    ...data,
    updatedAt: new Date().toISOString()
  }
  await database.put(STORE_NAME, draft)
  return draft
}

export async function getDraft(articleId) {
  const database = await getDB()
  return await database.get(STORE_NAME, articleId || 'new')
}

export async function getAllDrafts() {
  const database = await getDB()
  return await database.getAllFromIndex(STORE_NAME, 'updatedAt')
}

export async function deleteDraft(articleId) {
  const database = await getDB()
  await database.delete(STORE_NAME, articleId || 'new')
}

export async function clearDrafts() {
  const database = await getDB()
  await database.clear(STORE_NAME)
}
