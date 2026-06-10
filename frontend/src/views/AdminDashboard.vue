<template>
    <AdminLayout page-title="仪表盘">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            <div class="card">
                <div class="card-body">
                    <div class="text-sm text-gray-500">商品总数</div>
                    <div class="text-2xl font-bold text-blue-600 mt-2">{{ stats.product_total || 0 }}</div>
                </div>
            </div>
            <div class="card">
                <div class="card-body">
                    <div class="text-sm text-gray-500">订单总数</div>
                    <div class="text-2xl font-bold text-green-600 mt-2">{{ stats.order_total || 0 }}</div>
                </div>
            </div>
            <div class="card">
                <div class="card-body">
                    <div class="text-sm text-gray-500">用户总数</div>
                    <div class="text-2xl font-bold text-gray-800 mt-2">{{ stats.user_total || 0 }}</div>
                </div>
            </div>
            <div class="card">
                <div class="card-body">
                    <div class="text-sm text-gray-500">总收入</div>
                    <div class="text-2xl font-bold text-yellow-600 mt-2">￥{{ stats.amount_total || 0 }}</div>
                </div>
            </div>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="card">
                <div class="card-header font-semibold">最近订单</div>
                <div class="card-body">
                    <table class="table">
                        <thead>
                            <tr>
                                <th>订单号</th>
                                <th>金额</th>
                                <th>状态</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="order in recentOrders" :key="order.order_no">
                                <td>{{ order.order_no }}</td>
                                <td>￥{{ order.amount }}</td>
                                <td>
                                    <span v-if="order.status === 1" class="badge-yellow">待支付</span>
                                    <span v-else-if="order.status === 2" class="badge-green">已完成</span>
                                    <span v-else class="badge-red">已取消</span>
                                </td>
                            </tr>
                            <tr v-if="recentOrders.length === 0">
                                <td colspan="3" class="text-center text-gray-500">暂无数据</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <div class="card">
                <div class="card-header font-semibold">系统信息</div>
                <div class="card-body">
                    <table class="table">
                        <tbody>
                            <tr><td class="w-32 text-gray-500">版本</td><td>v1.0.0</td></tr>
                            <tr><td class="text-gray-500">时间</td><td>{{ currentTime }}</td></tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminDashboard } from '../api/admin'

const stats = ref({})
const recentOrders = ref([])
const currentTime = ref('')

function loadData() {
    adminDashboard().then((res) => {
        stats.value = res.data && res.data.stats ? res.data.stats : {}
        recentOrders.value = res.data && res.data.orders ? res.data.orders : []
    }).catch(() => {})
    const d = new Date()
    currentTime.value = d.toLocaleString('zh-CN')
}

onMounted(loadData)
</script>
