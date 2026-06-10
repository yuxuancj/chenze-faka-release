<template>
    <Layout>
        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <h2 class="text-xl font-bold text-gray-800">商品列表</h2>
                <div class="flex items-center space-x-2">
                    <select v-model="filters.category_id" class="form-input w-32">
                        <option value="">全部分类</option>
                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                    </select>
                    <input v-model="filters.keyword" type="text" placeholder="搜索商品" class="form-input w-48" @keyup.enter="loadProducts">
                    <button @click="loadProducts" class="btn-primary">搜索</button>
                </div>
            </div>

            <div v-if="products.length === 0" class="card p-8 text-center text-gray-500">
                暂无商品
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                <div v-for="product in products" :key="product.id" class="card p-4">
                    <div class="bg-gray-100 h-48 rounded-md mb-3 flex items-center justify-center text-gray-400">
                        {{ product.image ? product.image : '商品图片' }}
                    </div>
                    <h3 class="font-semibold text-gray-800 truncate">{{ product.name }}</h3>
                    <p class="text-blue-600 font-bold mt-1">￥{{ product.price }}</p>
                    <p v-if="product.description" class="text-sm text-gray-500 mt-1 line-clamp-2">{{ product.description }}</p>
                    <div class="flex items-center justify-between mt-3">
                        <span class="text-sm text-gray-500">库存: {{ product.stock || 0 }}</span>
                        <router-link :to="'/product/' + product.id" class="btn-primary btn-sm">查看详情</router-link>
                    </div>
                </div>
            </div>

            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center space-x-2 mt-6">
                <button @click="prevPage" :disabled="page <= 1" class="btn-secondary btn-sm">上一页</button>
                <span class="text-sm text-gray-600">第 {{ page }} / {{ totalPages }} 页</span>
                <button @click="nextPage" :disabled="page >= totalPages" class="btn-secondary btn-sm">下一页</button>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { productList, categoryList } from '../api/product'

const products = ref([])
const categories = ref([])
const page = ref(1)
const pageSize = 20
const pagination = ref({ total: 0, size: 20 })
const filters = reactive({ category_id: '', keyword: '' })

const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.size) || 1)

function loadCategories() {
    categoryList().then((res) => {
        categories.value = res.data && res.data.list ? res.data.list : []
    }).catch(() => {})
}

function loadProducts() {
    productList(page.value, pageSize, filters.category_id || 0, filters.keyword || '').then((res) => {
        products.value = res.data && res.data.list ? res.data.list : []
        pagination.value.total = (res.data && res.data.total) ? res.data.total : 0
    }).catch(() => {})
}

function prevPage() {
    if (page.value > 1) {
        page.value--
        loadProducts()
    }
}

function nextPage() {
    if (page.value < totalPages.value) {
        page.value++
        loadProducts()
    }
}

onMounted(() => {
    loadCategories()
    loadProducts()
})
</script>
