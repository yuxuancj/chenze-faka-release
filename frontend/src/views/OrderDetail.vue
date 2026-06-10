<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">订单详情</h2>
            <div v-if="loading" class="card p-8 text-center text-gray-500">
                加载中...
            </div>
            <div v-else-if="!order.order_no" class="card p-8 text-center text-gray-500">
                订单不存在
            </div>
            <div v-else class="space-y-4">
                <div class="card">
                    <div class="card-header font-semibold">订单信息</div>
                    <div class="card-body overflow-x-auto">
                        <table class="table w-full">
                            <tbody>
                                <tr><td class="w-32 text-gray-500">订单号</td><td>{{ order.order_no }}</td></tr>
                                <tr><td class="text-gray-500">商品名称</td><td>{{ order.product_snapshot || order.product_name || '-' }}</td></tr>
                                <tr><td class="text-gray-500">购买数量</td><td>{{ order.quantity || 1 }}</td></tr>
                                <tr><td class="text-gray-500">订单金额</td><td>￥{{ order.amount }}</td></tr>
                                <tr><td class="text-gray-500">状态</td>
                                    <td>
                                        <span v-if="order.status === 0" class="badge-yellow">待支付</span>
                                        <span v-else-if="order.status === 1" class="badge-green">已支付</span>
                                        <span v-else-if="order.status === 2" class="badge-blue">已完成</span>
                                        <span v-else-if="order.status === 3" class="badge-red">已关闭</span>
                                        <span v-else class="badge-gray">未知</span>
                                    </td>
                                </tr>
                                <tr><td class="text-gray-500">邮箱</td><td>{{ order.email }}</td></tr>
                                <tr><td class="text-gray-500">备注</td><td>{{ order.remark || '-' }}</td></tr>
                                <tr><td class="text-gray-500">创建时间</td><td>{{ order.created_at }}</td></tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div v-if="order.status === 0" class="card">
                    <div class="card-body flex items-center justify-between">
                        <span>待支付金额: <span class="text-blue-600 font-bold text-xl">￥{{ order.amount }}</span></span>
                        <button @click="payNow" :disabled="paying" class="btn-primary">
                            {{ paying ? '支付中...' : '立即支付' }}
                        </button>
                    </div>
                </div>
                <div v-if="cards && cards.length > 0" class="card">
                    <div class="card-header font-semibold">卡密信息</div>
                    <div class="card-body overflow-x-auto">
                        <table class="table w-full">
                            <thead>
                                <tr>
                                    <th>卡密</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(card, idx) in cards" :key="idx">
                                    <td class="font-mono">{{ card.card_data || card.card_no || card }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Layout from '../components/Layout.vue'
import { orderDetail, payOrder } from '../api/order'

const route = useRoute()
const order = ref({})
const cards = ref([])
const loading = ref(true)
const paying = ref(false)

function loadDetail() {
    loading.value = true
    orderDetail(route.params.order_no).then((data) => {
        if (data && data.order) {
            order.value = data.order
            cards.value = data.cards || []
        } else if (data) {
            order.value = data
            cards.value = data.cards || []
        }
    }).catch(() => {
        order.value = {}
        cards.value = []
    }).finally(() => {
        loading.value = false
    })
}

function payNow() {
    paying.value = true
    payOrder(order.value.order_no).then((data) => {
        if (data && data.pay_url) {
            window.location.href = data.pay_url
        } else {
            alert('支付成功')
            loadDetail()
        }
    }).catch(() => {}).finally(() => {
        paying.value = false
    })
}

onMounted(loadDetail)
</script>
