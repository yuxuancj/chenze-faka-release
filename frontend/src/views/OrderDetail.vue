<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">订单详情</h2>
            <el-card v-if="loading" class="text-center" shadow="never">
                <el-skeleton active :rows="5" />
            </el-card>
            <el-card v-else-if="!order.order_no" class="text-center" shadow="never">
                <el-empty description="订单不存在" />
            </el-card>
            <template v-else>
                <el-card shadow="never">
                    <template #header>
                        <span class="font-semibold">订单信息</span>
                    </template>
                    <el-descriptions :column="1" border>
                        <el-descriptions-item label="订单号">{{ order.order_no }}</el-descriptions-item>
                        <el-descriptions-item label="商品名称">{{ order.product_snapshot || order.product_name || '-' }}</el-descriptions-item>
                        <el-descriptions-item label="购买数量">{{ order.quantity || 1 }}</el-descriptions-item>
                        <el-descriptions-item label="订单金额">￥{{ order.amount }}</el-descriptions-item>
                        <el-descriptions-item label="状态">
                            <el-tag :type="statusType(order.status)" effect="light">
                                {{ statusText(order.status) }}
                            </el-tag>
                        </el-descriptions-item>
                        <el-descriptions-item label="邮箱">{{ order.email }}</el-descriptions-item>
                        <el-descriptions-item label="备注">{{ order.remark || '-' }}</el-descriptions-item>
                        <el-descriptions-item label="创建时间">{{ order.created_at }}</el-descriptions-item>
                    </el-descriptions>
                </el-card>
                <el-card v-if="order.status === 0" shadow="never">
                    <div class="flex items-center justify-between">
                        <span>待支付金额: <span class="text-blue-600 font-bold text-xl">￥{{ order.amount }}</span></span>
                        <el-button type="primary" :loading="paying" @click="payNow">立即支付</el-button>
                    </div>
                </el-card>
                <el-card v-if="cards && cards.length > 0" shadow="never">
                    <template #header>
                        <span class="font-semibold">卡密信息</span>
                    </template>
                    <el-table :data="cards" stripe>
                        <el-table-column label="卡号">
                            <template #default="scope">
                                <span class="font-mono">{{ scope.row.card_data || scope.row.card_no || scope.row }}</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-card>
            </template>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Layout from '../components/Layout.vue'
import { orderDetail, payOrder } from '../api/order'
import { ElMessage } from 'element-plus'

const route = useRoute()
const order = ref({})
const cards = ref([])
const loading = ref(true)
const paying = ref(false)

function statusType(status) {
    if (status === 0) return 'warning'
    if (status === 1) return 'success'
    if (status === 2) return 'primary'
    if (status === 3) return 'danger'
    return 'info'
}

function statusText(status) {
    if (status === 0) return '待支付'
    if (status === 1) return '已支付'
    if (status === 2) return '已完成'
    if (status === 3) return '已关闭'
    return '未知'
}

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
            ElMessage.success('支付成功')
            loadDetail()
        }
    }).catch(() => {}).finally(() => {
        paying.value = false
    })
}

onMounted(loadDetail)
</script>
