<template>
    <Layout>
        <div v-if="loading" class="text-center py-16">
            <el-icon class="is-loading text-blue-600 text-4xl"><Loading /></el-icon>
            <p class="text-gray-500 mt-4">加载中...</p>
        </div>
        <div v-else-if="!product.id" class="text-center py-16 text-gray-500">
            商品不存在
            <div class="mt-4">
                <router-link to="/products">
                    <el-button type="primary">返回商品列表</el-button>
                </router-link>
            </div>
        </div>
        <div v-else>
            <el-row :gutter="20">
                <el-col :xs="24" :md="12">
                    <el-card shadow="never">
                        <div class="product-image-area flex items-center justify-center text-gray-500 text-lg">
                            {{ product.name }}
                        </div>
                    </el-card>
                </el-col>
                <el-col :xs="24" :md="12">
                    <el-card shadow="never">
                        <template #header>
                            <h1 class="text-2xl font-bold text-gray-800">{{ product.name }}</h1>
                        </template>
                        <div class="space-y-4">
                            <div class="text-3xl font-bold text-blue-600">￥{{ product.price }}</div>
                            <div class="flex items-center gap-3">
                                <el-tag type="success">库存: {{ product.stock || 0 }}</el-tag>
                                <el-tag>已售: {{ product.sales || 0 }}</el-tag>
                            </div>
                            <el-input
                                v-model="product.description"
                                type="textarea"
                                :rows="4"
                                placeholder="暂无商品描述"
                                readonly
                                style="margin-bottom: 8px;"
                            />
                            <div class="flex items-center gap-3">
                                <span class="text-gray-700">数量</span>
                                <el-input-number v-model="quantity" :min="1" :max="product.stock || 999" />
                            </div>
                            <div class="flex items-center gap-3 pt-2">
                                <el-button type="primary" size="large" @click="buyNow">立即购买</el-button>
                                <el-button size="large" @click="addToCart">加入购物车</el-button>
                            </div>
                        </div>
                    </el-card>
                </el-col>
            </el-row>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
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
    ElMessage.success('已添加到购物车')
}

function buyNow() {
    if (!product.value || !product.value.id) return
    router.push('/checkout?product_id=' + product.value.id + '&quantity=' + quantity.value)
}

onMounted(loadDetail)
</script>

<style scoped>
.product-image-area {
    height: 400px;
    background: linear-gradient(135deg, #f7f7f7 0%, #e8e8e8 100%);
    border-radius: 4px;
}
</style>
