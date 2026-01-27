---
name: vue3-patterns
description: Vue 3 and Nuxt composition API patterns, best practices, and advanced techniques for modern frontend development.
---

# Vue 3 & Nuxt Advanced Patterns

Vue 3 Composition API 和 Nuxt 3 高級模式。

## Composition API 基礎

### ✅ GOOD: 使用 <script setup>

```vue
<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'

const count = ref(0)
const doubled = computed(() => count.value * 2)

watch(() => count.value, (newVal) => {
  console.log('Count changed to:', newVal)
})

onMounted(() => {
  console.log('Component mounted')
})

const increment = () => count.value++
</script>

<template>
  <div>
    <p>Count: {{ count }}</p>
    <p>Doubled: {{ doubled }}</p>
    <button @click="increment">Increment</button>
  </div>
</template>
```

### ❌ AVOID: 舊的 Options API

```vue
<script>
export default {
  data() {
    return { count: 0 }
  },
  computed: {
    doubled() { return this.count * 2 }
  },
  methods: {
    increment() { this.count++ }
  }
}
</script>
```

## 可組合函數 (Composables)

### 數據獲取 Composable

```typescript
// composables/useFetch.ts
import { ref, computed } from 'vue'

interface FetchOptions<T> {
  onSuccess?: (data: T) => void
  onError?: (error: Error) => void
  immediate?: boolean
}

export function useFetchData<T>(
  url: string,
  options: FetchOptions<T> = {}
) {
  const data = ref<T | null>(null)
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  const isFetched = computed(() => data.value !== null)
  const isError = computed(() => error.value !== null)

  const fetch = async () => {
    isLoading.value = true
    error.value = null

    try {
      const response = await $fetch<T>(url)
      data.value = response
      options.onSuccess?.(response)
    } catch (err) {
      error.value = err as Error
      options.onError?.(error.value)
    } finally {
      isLoading.value = false
    }
  }

  const refetch = () => fetch()

  if (options.immediate !== false) {
    fetch()
  }

  return {
    data,
    isLoading,
    error,
    isFetched,
    isError,
    fetch,
    refetch
  }
}
```

### 本地存儲 Composable

```typescript
// composables/useLocalStorage.ts
import { ref, watch } from 'vue'

export function useLocalStorage<T>(key: string, initialValue: T) {
  // 從本地存儲讀取初始值
  const stored = process.client ? localStorage.getItem(key) : null
  const data = ref<T>(stored ? JSON.parse(stored) : initialValue)

  // 監聽變化並保存到本地存儲
  watch(
    () => data.value,
    (newValue) => {
      if (process.client) {
        localStorage.setItem(key, JSON.stringify(newValue))
      }
    },
    { deep: true }
  )

  // 清除
  const clear = () => {
    data.value = initialValue
    if (process.client) {
      localStorage.removeItem(key)
    }
  }

  return { data, clear }
}

// 使用
const theme = useLocalStorage('theme', 'light')
theme.data.value = 'dark'  // 自動保存
```

### 防抖 Composable

```typescript
// composables/useDebounce.ts
import { ref, watch } from 'vue'

export function useDebounce<T>(value: Ref<T>, delay: number = 500) {
  const debounced = ref(value.value) as Ref<T>
  let timeout: NodeJS.Timeout | null = null

  watch(
    () => value.value,
    (newValue) => {
      if (timeout) clearTimeout(timeout)

      timeout = setTimeout(() => {
        debounced.value = newValue
      }, delay)
    }
  )

  return debounced
}

// 使用
const searchQuery = ref('')
const debouncedQuery = useDebounce(searchQuery, 300)

watch(debouncedQuery, async (query) => {
  if (query) {
    await performSearch(query)
  }
})
```

## 狀態管理 (Pinia)

### 完整 Store 示例

```typescript
// stores/market.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useMarketStore = defineStore('market', () => {
  // State
  const markets = ref<Market[]>([])
  const selectedMarketId = ref<string | null>(null)
  const filters = ref({
    status: 'all' as const,
    search: '',
    sortBy: 'name' as const
  })

  // Getters
  const selectedMarket = computed(() => {
    return markets.value.find(m => m.id === selectedMarketId.value) || null
  })

  const filteredAndSortedMarkets = computed(() => {
    let result = [...markets.value]

    // 篩選狀態
    if (filters.value.status !== 'all') {
      result = result.filter(m => m.status === filters.value.status)
    }

    // 篩選搜索
    if (filters.value.search) {
      const query = filters.value.search.toLowerCase()
      result = result.filter(m =>
        m.name.toLowerCase().includes(query) ||
        m.description.toLowerCase().includes(query)
      )
    }

    // 排序
    result.sort((a, b) => {
      switch (filters.value.sortBy) {
        case 'name':
          return a.name.localeCompare(b.name)
        case 'volume':
          return b.volume - a.volume
        case 'date':
          return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
        default:
          return 0
      }
    })

    return result
  })

  // Actions
  const fetchMarkets = async () => {
    try {
      const data = await $fetch<Market[]>('/api/markets')
      markets.value = data
    } catch (error) {
      console.error('Failed to fetch markets:', error)
    }
  }

  const selectMarket = (id: string) => {
    selectedMarketId.value = id
  }

  const updateFilter = (key: keyof typeof filters.value, value: any) => {
    filters.value[key] = value
  }

  const clearFilters = () => {
    filters.value = {
      status: 'all',
      search: '',
      sortBy: 'name'
    }
  }

  return {
    markets,
    selectedMarket,
    filters,
    filteredAndSortedMarkets,
    fetchMarkets,
    selectMarket,
    updateFilter,
    clearFilters
  }
})
```

## 高級組件模式

### 作用域插槽

```vue
<template>
  <div class="data-list">
    <div
      v-for="(item, index) in items"
      :key="item.id"
      class="list-item"
    >
      <!-- 將項目和索引暴露給父組件 -->
      <slot :item="item" :index="index" :isLast="index === items.length - 1">
        <!-- 默認內容 -->
        <span>{{ item.name }}</span>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props<T> {
  items: T[]
}

withDefaults(defineProps<Props>(), {})
</script>

<!-- 使用 -->
<template>
  <DataList :items="markets">
    <template #default="{ item, index, isLast }">
      <div class="market-item">
        <span>{{ index + 1 }}. {{ item.name }}</span>
        <span v-if="isLast" class="last-item">最後一項</span>
      </div>
    </template>
  </DataList>
</template>
```

### 動態組件

```vue
<template>
  <div>
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab"
        :class="{ active: activeTab === tab }"
        @click="activeTab = tab"
      >
        {{ tab }}
      </button>
    </div>

    <div class="tab-content">
      <!-- 動態切換組件 -->
      <component :is="tabComponents[activeTab]" :data="data" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Overview from './tabs/Overview.vue'
import Details from './tabs/Details.vue'
import Analytics from './tabs/Analytics.vue'

const activeTab = ref('overview')
const tabs = ['overview', 'details', 'analytics']

const tabComponents = {
  overview: Overview,
  details: Details,
  analytics: Analytics
}
</script>
```

## 自定義指令

```typescript
// directives/v-focus.ts
export const vFocus = {
  mounted(el: HTMLElement) {
    el.focus()
  }
}

// 使用
<template>
  <input v-focus type="text" />
</template>

// directives/v-click-outside.ts
export const vClickOutside = {
  mounted(el: HTMLElement, binding: DirectiveBinding<() => void>) {
    const clickHandler = (event: MouseEvent) => {
      if (!el.contains(event.target as Node)) {
        binding.value()
      }
    }

    document.addEventListener('click', clickHandler)

    el._clickOutsideListener = clickHandler as any
  },

  unmounted(el: HTMLElement) {
    document.removeEventListener('click', el._clickOutsideListener)
    delete el._clickOutsideListener
  }
}

// 使用
<template>
  <div v-click-outside="closeDropdown" class="dropdown">
    <!-- 內容 -->
  </div>
</template>
```

## 錯誤處理

```typescript
// composables/useAsyncState.ts
export function useAsyncState<T>(
  asyncFn: () => Promise<T>,
  initialValue: T
) {
  const state = ref(initialValue)
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  const execute = async () => {
    isLoading.value = true
    error.value = null

    try {
      state.value = await asyncFn()
    } catch (err) {
      error.value = err instanceof Error ? err : new Error(String(err))
    } finally {
      isLoading.value = false
    }
  }

  onMounted(() => execute())

  return { state, isLoading, error, execute }
}

// 使用
const { state: markets, isLoading, error, execute } = useAsyncState(
  () => $fetch('/api/markets'),
  [] as Market[]
)
```

## 性能優化

### 虛擬列表

```vue
<template>
  <div ref="container" class="virtual-list">
    <div
      v-for="item in visibleItems"
      :key="item.id"
      class="list-item"
      :style="{ transform: `translateY(${item.offset}px)` }"
    >
      {{ item.name }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const ITEM_HEIGHT = 50
const VISIBLE_COUNT = 10

const items = ref<Item[]>([])
const scrollTop = ref(0)

const visibleItems = computed(() => {
  const startIndex = Math.floor(scrollTop.value / ITEM_HEIGHT)
  const endIndex = startIndex + VISIBLE_COUNT

  return items.value.slice(startIndex, endIndex).map((item, index) => ({
    ...item,
    offset: (startIndex + index) * ITEM_HEIGHT
  }))
})

const handleScroll = (event: Event) => {
  scrollTop.value = (event.target as HTMLElement).scrollTop
}
</script>
```

## Nuxt 特定模式

### 中間件

```typescript
// middleware/auth.ts
export default defineRouteMiddleware((to, from) => {
  const user = useAuthStore().user

  if (!user && to.meta.requiresAuth) {
    return navigateTo('/login')
  }
})

// 在頁面中使用
definePageMeta({
  middleware: 'auth',
  requiresAuth: true
})
```

### 插件

```typescript
// plugins/globalProperties.ts
export default defineNuxtPlugin((nuxtApp) => {
  // 提供全局屬性
  nuxtApp.provide('formatCurrency', (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(value)
  })
})

// 使用
<template>
  <span>{{ $formatCurrency(100) }}</span>
</template>
```

---

**Remember**: Vue 3 的 Composition API 提供了強大的可組合性。正確使用可以創建可重用、可測試、易於維護的代碼。
