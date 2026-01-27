---
name: vue3-nuxt
description: Vue 3 and Nuxt 3 development specialist. Handles component patterns, composables, state management, and modern frontend best practices.
tools: Read, Write, Edit, Bash, Grep, Glob
model: sonnet
---

# Vue 3 & Nuxt 3 Development Patterns

Vue 3 和 Nuxt 3 應用開發的架構模式和最佳實踐。

## Composition API Basics

使用 `<script setup>` 語法：

```vue
<template>
  <div class="market-card">
    <h3>{{ market.name }}</h3>
    <p>{{ market.description }}</p>
    <button @click="handleTrade" :disabled="isLoading">
      {{ isLoading ? 'Loading...' : 'Trade' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface Market {
  id: string
  name: string
  description: string
  volume: number
}

interface Props {
  market: Market
}

const props = defineProps<Props>()
const emit = defineEmits<{
  trade: [amount: number]
}>()

const isLoading = ref(false)
const tradeAmount = ref(0)

const isTradeDisabled = computed(() => {
  return isLoading.value || tradeAmount.value <= 0
})

const handleTrade = async () => {
  isLoading.value = true
  try {
    emit('trade', tradeAmount.value)
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.market-card {
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
}
</style>
```

## Custom Composables

可重用邏輯：

```typescript
// composables/useMarketSearch.ts
import { ref, computed } from 'vue'
import { useFetch } from '#app'

export function useMarketSearch() {
  const query = ref('')
  const results = ref<Market[]>([])
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  const debouncedQuery = useDebounce(query, 500)

  const search = async () => {
    if (!debouncedQuery.value) {
      results.value = []
      return
    }

    isLoading.value = true
    error.value = null

    try {
      const { data } = await useFetch('/api/markets/search', {
        query: { q: debouncedQuery.value }
      })
      results.value = data.value?.data || []
    } catch (err) {
      error.value = err as Error
      results.value = []
    } finally {
      isLoading.value = false
    }
  }

  watch(debouncedQuery, search)

  return {
    query,
    results,
    isLoading,
    error,
    search: () => search()
  }
}

// composables/useDebounce.ts
export function useDebounce<T>(value: Ref<T>, delay: number): Ref<T> {
  const debounced = ref(value.value) as Ref<T>

  const timeout = setTimeout(() => {
    debounced.value = value.value
  }, delay)

  watch(value, () => {
    clearTimeout(timeout)
    setTimeout(() => {
      debounced.value = value.value
    }, delay)
  })

  return debounced
}
```

## State Management with Pinia

全局狀態管理：

```typescript
// stores/marketStore.ts
import { defineStore } from 'pinia'

export const useMarketStore = defineStore('market', () => {
  // State
  const markets = ref<Market[]>([])
  const selectedMarket = ref<Market | null>(null)
  const isLoading = ref(false)
  const filters = ref({
    status: 'all',
    sortBy: 'volume'
  })

  // Getters
  const filteredMarkets = computed(() => {
    return markets.value.filter(m => {
      if (filters.value.status !== 'all') {
        return m.status === filters.value.status
      }
      return true
    }).sort((a, b) => {
      if (filters.value.sortBy === 'volume') {
        return b.volume - a.volume
      }
      return a.name.localeCompare(b.name)
    })
  })

  // Actions
  async function fetchMarkets() {
    isLoading.value = true
    try {
      const { data } = await $fetch('/api/markets')
      markets.value = data
    } finally {
      isLoading.value = false
    }
  }

  async function selectMarket(id: string) {
    selectedMarket.value = markets.value.find(m => m.id === id) || null
  }

  function setFilter(key: string, value: any) {
    filters.value[key] = value
  }

  return {
    markets,
    selectedMarket,
    isLoading,
    filters,
    filteredMarkets,
    fetchMarkets,
    selectMarket,
    setFilter
  }
})
```

## Component Patterns

### Form Component

```vue
<template>
  <form @submit.prevent="handleSubmit">
    <div class="form-group">
      <label for="name">Market Name</label>
      <input
        id="name"
        v-model="formData.name"
        type="text"
        :aria-invalid="!!errors.name"
        @blur="validateField('name')"
      />
      <span v-if="errors.name" class="error">{{ errors.name }}</span>
    </div>

    <div class="form-group">
      <label for="description">Description</label>
      <textarea
        id="description"
        v-model="formData.description"
        @blur="validateField('description')"
      />
      <span v-if="errors.description" class="error">{{ errors.description }}</span>
    </div>

    <button type="submit" :disabled="isSubmitting">
      {{ isSubmitting ? 'Creating...' : 'Create Market' }}
    </button>
  </form>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'

interface FormData {
  name: string
  description: string
}

interface Errors {
  name?: string
  description?: string
}

const formData = reactive<FormData>({
  name: '',
  description: ''
})

const errors = ref<Errors>({})
const isSubmitting = ref(false)

const validateField = (field: keyof FormData) => {
  if (!formData[field].trim()) {
    errors.value[field] = `${field} is required`
  } else {
    delete errors.value[field]
  }
}

const handleSubmit = async () => {
  // Validate all fields
  Object.keys(formData).forEach(field => {
    validateField(field as keyof FormData)
  })

  if (Object.keys(errors.value).length > 0) return

  isSubmitting.value = true
  try {
    await $fetch('/api/markets', {
      method: 'POST',
      body: formData
    })
    // Success handling
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.form-group {
  margin-bottom: 1rem;
}

input, textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
}

input[aria-invalid="true"] {
  border-color: #dc2626;
}

.error {
  color: #dc2626;
  font-size: 0.875rem;
  margin-top: 0.25rem;
  display: block;
}
</style>
```

### Data Table Component

```vue
<template>
  <div class="data-table">
    <div class="table-controls">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search..."
        class="search-input"
      />
      <select v-model="sortBy" class="sort-select">
        <option value="name">Sort by Name</option>
        <option value="volume">Sort by Volume</option>
        <option value="date">Sort by Date</option>
      </select>
    </div>

    <table v-if="filteredItems.length > 0">
      <thead>
        <tr>
          <th v-for="column in columns" :key="column.key">
            {{ column.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in paginatedItems" :key="item.id">
          <td v-for="column in columns" :key="column.key">
            {{ getNestedValue(item, column.key) }}
          </td>
        </tr>
      </tbody>
    </table>

    <div v-else class="no-data">No data available</div>

    <div v-if="totalPages > 1" class="pagination">
      <button
        :disabled="currentPage === 1"
        @click="currentPage--"
      >
        Previous
      </button>
      <span>Page {{ currentPage }} of {{ totalPages }}</span>
      <button
        :disabled="currentPage === totalPages"
        @click="currentPage++"
      >
        Next
      </button>
    </div>
  </div>
</template>

<script setup lang="ts" generic="T extends { id: string }">
import { ref, computed } from 'vue'

interface Column<K extends string> {
  key: K
  label: string
}

interface Props<K extends string> {
  items: T[]
  columns: Column<K>[]
  pageSize?: number
}

const props = withDefaults(defineProps<Props>(), {
  pageSize: 10
})

const searchQuery = ref('')
const sortBy = ref('name')
const currentPage = ref(1)

const filteredItems = computed(() => {
  return props.items.filter(item => {
    if (!searchQuery.value) return true
    return JSON.stringify(item)
      .toLowerCase()
      .includes(searchQuery.value.toLowerCase())
  })
})

const totalPages = computed(() => {
  return Math.ceil(filteredItems.value.length / props.pageSize)
})

const paginatedItems = computed(() => {
  const start = (currentPage.value - 1) * props.pageSize
  return filteredItems.value.slice(start, start + props.pageSize)
})

const getNestedValue = (obj: any, path: string) => {
  return path.split('.').reduce((acc, part) => acc?.[part], obj)
}
</script>

<style scoped>
.data-table {
  width: 100%;
}

.table-controls {
  margin-bottom: 1rem;
  display: flex;
  gap: 1rem;
}

.search-input,
.sort-select {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

th {
  background-color: #f5f5f5;
  font-weight: 600;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1rem;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
```

## Page Layouts

```vue
<!-- layouts/default.vue -->
<template>
  <div class="layout">
    <header class="header">
      <nav>
        <NuxtLink to="/">Home</NuxtLink>
        <NuxtLink to="/markets">Markets</NuxtLink>
      </nav>
    </header>

    <main class="main">
      <slot />
    </main>

    <footer class="footer">
      <p>&copy; 2024 Hiskio App</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'default'
})
</script>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  background-color: #1a202c;
  color: white;
  padding: 1rem;
}

.main {
  flex: 1;
  padding: 2rem;
}

.footer {
  background-color: #1a202c;
  color: white;
  padding: 1rem;
  text-align: center;
}
</style>
```

## Nuxt Auto Imports

Nuxt 自動導入常用函數（無需手動 import）：

```vue
<script setup lang="ts">
// ✅ 這些會自動導入
// ref, computed, watch, onMounted 等來自 'vue'
// definePageMeta, navigateTo 等來自 'nuxt'
// useRoute, useRouter 等來自 '#app'

const route = useRoute()
const router = useRouter()
const { $fetch } = useNuxtApp()
</script>
```

## Error Handling

```typescript
// middleware/error.ts
export default defineEventHandler((event) => {
  if (event.node.req.url?.includes('/api')) {
    return errorHandler(event)
  }
})

function errorHandler(event: any) {
  event.node.res.statusCode = 500
  return { error: 'Internal Server Error' }
}
```

## Performance Tips

### Lazy Loading Components

```vue
<template>
  <div>
    <Suspense>
      <template #default>
        <HeavyChart :data="data" />
      </template>
      <template #fallback>
        <div>Loading chart...</div>
      </template>
    </Suspense>
  </div>
</template>

<script setup lang="ts">
const HeavyChart = defineAsyncComponent(
  () => import('./HeavyChart.vue')
)
</script>
```

### Image Optimization

```vue
<template>
  <NuxtImg
    src="/images/market.png"
    alt="Market image"
    width="800"
    height="600"
    format="webp"
    quality="80"
  />
</template>
```

## Testing Pattern

```typescript
// tests/unit/composables/useMarketSearch.test.ts
import { describe, it, expect, vi } from 'vitest'
import { useMarketSearch } from '~/composables/useMarketSearch'

describe('useMarketSearch', () => {
  it('fetches markets when query changes', async () => {
    const { query, results, search } = useMarketSearch()

    query.value = 'election'
    await search()

    expect(results.value.length).toBeGreaterThan(0)
  })

  it('debounces search query', async () => {
    const { query } = useMarketSearch()
    const fetchSpy = vi.spyOn(global, 'fetch')

    query.value = 'test'
    query.value = 'test2'

    await vi.waitFor(() => {
      expect(fetchSpy).toHaveBeenCalledTimes(1)
    })
  })
})
```

---

**Remember**: Vue 3's Composition API enables clean, composable, and testable code. Combine it with Nuxt's powerful features for a fantastic development experience.
