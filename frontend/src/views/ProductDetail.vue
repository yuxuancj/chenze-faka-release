<template>
    <Layout>
        <div v-if="loading" class="card p-8 text-center text-gray-500">
            加载中...
        </div>
        <div v-else-if="!product.id" class="card p-8 text-center text-gray-500">
            商品不存在
            <div class="mt-4">
                <router-link to="/products" class="btn-primary">返回商品列表</router-link>
            </div>
        </div>
        <div v-else class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="card">
                    <div class="card-body">
                        <div class="h-96 bg-gradient-to-br from-gray-100 to-gray-200 rounded-md flex items-center justify-center text-gray-500 text-lg">
                            {{ product.name }}
                        </div>
                    </div>
                </div>
                <div class="card">
                    <div class="card-body space-y-4">
                        <h1 class="text-2xl font-bold text-gray-800">{{ product.name }}</h1>
                        <div class="text-3xl font-bold text-blue-600">￥{{ product.price }}</div>
                        <div class="flex items-center gap-4 text-sm">
                            <span class="badge-blue">库存: {{ product.stock || 0 }}</span>
                            <span class="badge-gray">已售: {{ product.sales || 0 }}</span>
                        </div>
                        <p class="text-gray-700 whitespace-pre-wrap">{{ product.description || '暂无商品描述' }}</p>
                        <div class="flex items-center gap-2 pt-2">
                            <label class="text-sm text-gray-700">数量</label>
                            <input v-model.number="quantity" type="number" min="1" class="form-input w-24">
                        </div>
                        <div class="flex items-center gap-3 pt-2">
                            <button @click="buyNow" class="btn-primary flex-1">立即购买</button>
                            <button @click="addToCart" class="btn-secondary flex-1">加入购物车</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'
import { useCartStore } from '../stores/cart'
import { productDetail } from '../api/product'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()

const product = ref({})
const quantity = ref(1)
const loading = ref(true)

function loadDetail() {
    loading.value = true
    productDetail(route.params.id).then((data) => {
        product.value = data || {}
    }).catch(() => {
        product.value = {}
    }).finally(() => {
        loading.value = false
    })
}

function addToCart() {
    if (!product.value || !product.value.id) return
    cartStore.addItem(product.value, quantity.value)
    alert('已添加到购物车')
}

function buyNow() {
    if (!product.value || !product.value.id) return
    router.push('/checkout?product_id=' + product.value.id + '&quantity=' + quantity.value)
}

onMounted(loadDetail)
</script>
