<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">我的订单</h2>
            <div class="card">
                <div class="card-body">
                    <table class="table">
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
                            <tr v-for="order in orders" :key="order.order_no">
                                <td>{{ order.order_no }}</td>
                                <td>{{ order.product_name || '-' }}</td>
                                <td>{{ order.quantity }}</td>
                                <td>￥{{ order.amount }}</td>
                                <td>
                                    <span v-if="order.status === 1" class="badge-yellow">待支付</span>
                                    <span v-else-if="order.status === 2" class="badge-green">已完成</span>
                                    <span v-else-if="order.status === 3" class="badge-red">已取消</span>
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
            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center space-x-2">
                <button @click="prevPage" :disabled="page <= 1" class="btn-sm btn-secondary">上一页</button>
                <span class="text-sm text-gray-600">第 {{ page }} / {{ totalPages }} 页</span>
                <button @click="nextPage" :disabled="page >= totalPages" class="btn-sm btn-secondary">下一页</button>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { orderList } from '../api/order'

const orders = ref([])
const page = ref(1)
const pageSize = 20
const pagination = ref({ total: 0, size: 20 })

const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.size) || 1)

function loadOrders() {
    orderList(page.value, pageSize).then((res) => {
        orders.value = res.data && res.data.list ? res.data.list : []
        pagination.value.total = (res.data && res.data.total) ? res.data.total : 0
    }).catch(() => {})
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
