<template>
    <Layout>
        <div v-if="product.id" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="card">
                    <div class="card-body">
                        <div class="bg-gray-100 h-96 rounded-md flex items-center justify-center text-gray-400 text-lg">
                            {{ product.image || '商品图片' }}
                        </div>
                    </div>
                </div>
                <div class="card">
                    <div class="card-body space-y-4">
                        <h1 class="text-2xl font-bold text-gray-800">{{ product.name }}</h1>
                        <div class="text-3xl font-bold text-blue-600">￥{{ product.price }}</div>
                        <div class="flex items-center space-x-2 text-sm text-gray-500">
                            <span class="badge-blue">库存: {{ product.stock || 0 }}</span>
                            <span class="badge-gray">已售: {{ product.sold || 0 }}</span>
                        </div>
                        <p class="text-gray-700 whitespace-pre-wrap">{{ product.description || '暂无商品描述' }}</p>
                        <div class="flex items-center space-x-2">
                            <label class="form-label">购买数量</label>
                            <input v-model.number="quantity" type="number" min="1" class="form-input w-24">
                        </div>
                        <div class="flex items-center space-x-3">
                            <button @click="addToCart" class="btn-primary">加入购物车</button>
                            <button @click="buyNow" class="btn-secondary">立即购买</button>
                        </div>
                    </div>
                </div>
            </div>

            <div class="card">
                <div class="card-header">
                    <h3 class="font-semibold">商品详情</h3>
                </div>
                <div class="card-body">
                    <table class="table">
                        <tbody>
                            <tr><td class="w-32 text-gray-500">商品编号</td><td>{{ product.id }}</td></tr>
                            <tr><td class="text-gray-500">商品名称</td><td>{{ product.name }}</td></tr>
                            <tr><td class="text-gray-500">商品价格</td><td>￥{{ product.price }}</td></tr>
                            <tr><td class="text-gray-500">商品库存</td><td>{{ product.stock || 0 }}</td></tr>
                            <tr><td class="text-gray-500">分类</td><td>{{ product.category_name || '未分类' }}</td></tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
        <div v-else class="card p-8 text-center text-gray-500">
            加载中...
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

function loadDetail() {
    productDetail(route.params.id).then((res) => {
        product.value = res.data || {}
    }).catch(() => {})
}

function addToCart() {
    if (!product.value || !product.value.id) return
    cartStore.addItem(product.value, quantity.value)
    alert('已加入购物车')
}

function buyNow() {
    if (!product.value || !product.value.id) return
    cartStore.clearCart()
    cartStore.addItem(product.value, quantity.value)
    router.push('/checkout')
}

onMounted(loadDetail)
</script>
