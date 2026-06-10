<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">我的订单</h2>
            <div v-if="loading" class="card p-8 text-center text-gray-500">
                加载中...
            </div>
            <div v-else class="card">
                <div class="card-body overflow-x-auto">
                    <table class="table w-full">
                        <thead>
                            <tr>
                                <th>订单号</th>
                                <th>商品</th>
                                <th>数量</th>
                                <th>金额</th>
                                <th>状态</th>
                                <th>创建时间</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="orders.length === 0">
                                <td colspan="7" class="text-center text-gray-500 py-8">暂无订单</td>
                            </tr>
                            <tr v-for="order in orders" :key="order.order_no || order.id">
                                <td>{{ order.order_no || '-' }}</td>
                                <td>{{ order.product_snapshot || order.product_name || '-' }}</td>
                                <td>{{ order.quantity || 1 }}</td>
                                <td>￥{{ order.amount }}</td>
                                <td>
                                    <span v-if="order.status === 0" class="badge-yellow">待支付</span>
                                    <span v-else-if="order.status === 1" class="badge-green">已支付</span>
                                    <span v-else-if="order.status === 2" class="badge-blue">已完成</span>
                                    <span v-else-if="order.status === 3" class="badge-red">已关闭</span>
                                    <span v-else class="badge-gray">未知</span>
                                </td>
                                <td>{{ order.created_at }}</td>
                                <td>
                                    <router-link :to="'/order/' + order.order_no" class="btn-sm btn-primary">查看</router-link>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center gap-2">
                <button @click="prevPage" :disabled="page <= 1" class="btn-sm btn-secondary">上一页</button>
                <span class="text-sm text-gray-600">第 {{ page }} / {{ totalPages }} 页</span>
                <button @click="nextPage" :disabled="page >= totalPages" class="btn-sm btn-secondary">下一页</button>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { orderList } from '../api/order'

const orders = ref([])
const page = ref(1)
const pageSize = 20
const pagination = reactive({ total: 0, size: 20 })
const loading = ref(false)

const totalPages = computed(() => Math.ceil(pagination.total / pagination.size) || 1)

function loadOrders() {
    loading.value = true
    orderList(page.value, pageSize).then((data) => {
        orders.value = (data && data.list) ? data.list : []
        pagination.total = (data && data.total) ? data.total : 0
        pagination.size = (data && data.size) ? data.size : pageSize
    }).catch(() => {
        orders.value = []
        pagination.total = 0
    }).finally(() => {
        loading.value = false
    })
}

function prevPage() {
    if (page.value > 1) {
        page.value--
        loadOrders()
    }
}

function nextPage() {
    if (page.value < totalPages.value) {
        page.value++
        loadOrders()
    }
}

onMounted(loadOrders)
</script>
