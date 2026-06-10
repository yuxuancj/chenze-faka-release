<template>
    <Layout>
        <div class="space-y-4">
            <el-card class="gradient-banner" body-style="padding: 0">
                <div class="px-6 py-8 md:px-12 md:py-10">
                    <h1 class="text-2xl md:text-3xl font-bold text-white mb-2">限时秒杀</h1>
                    <p class="text-red-100">精选商品，限时低价，售完即止</p>
                </div>
            </el-card>

            <el-empty v-if="loading" description="加载中..." />
            <el-empty v-else-if="activities.length === 0" description="暂无秒杀活动" />
            <el-row v-else :gutter="16">
                <el-col v-for="item in activities" :key="item.id" :xs="24" :sm="12" :md="8">
                    <el-card shadow="hover" class="seckill-card">
                        <div class="h-40 bg-gradient-to-br from-red-50 to-orange-100 flex items-center justify-center text-gray-500 text-sm">
                            {{ item.product_name || '秒杀商品' }}
                        </div>
                        <h3 class="font-semibold text-gray-800 mt-3 truncate">{{ item.product_name || '秒杀商品' }}</h3>
                        <div class="flex items-baseline gap-2 mt-2">
                            <span class="text-2xl font-bold text-red-600">￥{{ item.seckill_price }}</span>
                            <span class="text-sm text-gray-400 line-through">￥{{ item.original_price }}</span>
                        </div>
                        <div class="flex items-center gap-4 text-xs text-gray-600 mt-2">
                            <span>每人限购 {{ item.limit_per_user }} 件</span>
                            <span>秒杀库存: {{ item.stock }}</span>
                        </div>
                        <div class="mt-2">
                            <el-tag v-if="item.status === 'ongoing'" type="success" size="small">
                                进行中 · 剩余: {{ countdownFor(item) }}
                            </el-tag>
                            <el-tag v-else-if="item.status === 'upcoming'" type="warning" size="small">
                                即将开始 · {{ formatDateTime(item.start_time) }}
                            </el-tag>
                            <el-tag v-else type="info" size="small">已结束</el-tag>
                        </div>
                        <el-progress :percentage="stockProgress(item)" :color="'#ef4444'" :stroke-width="6" class="mt-3" />
                        <el-button v-if="item.status === 'ongoing'" type="danger" @click="handleBuy(item)" :loading="ordering" class="w-full mt-3">
                            立即抢购
                        </el-button>
                        <el-button v-else disabled class="w-full mt-3">
                            {{ item.status === 'upcoming' ? '尚未开始' : '已结束' }}
                        </el-button>
                    </el-card>
                </el-col>
            </el-row>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'
import { seckillActivities, seckillOrder } from '../api/seckill'
import { ElMessage } from 'element-plus'

const router = useRouter()
const activities = ref([])
const loading = ref(true)
const ordering = ref(false)
const now = ref(Date.now())
let timer = null

function formatDateTime(dateStr) {
    if (!dateStr) return '--'
    try {
        const d = new Date(dateStr)
        if (isNaN(d.getTime())) return '--'
        return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
    } catch (e) {
        return '--'
    }
}

function countdownFor(item) {
    if (!item.end_time) return '已结束'
    try {
        const end = new Date(item.end_time).getTime()
        if (isNaN(end)) return '已结束'
        const diff = Math.max(0, end - now.value)
        const h = Math.floor(diff / 3600000)
        const m = Math.floor((diff % 3600000) / 60000)
        const s = Math.floor((diff % 60000) / 1000)
        return h + '时' + String(m).padStart(2, '0') + '分' + String(s).padStart(2, '0') + '秒'
    } catch (e) {
        return '已结束'
    }
}

function stockProgress(item) {
    const total = (item.original_stock || item.stock || 0) + (item.sold || 0)
    if (total <= 0) return 0
    return Math.floor(((item.sold || 0) / total) * 100)
}

function load() {
    loading.value = true
    seckillActivities().then((data) => {
        activities.value = (data && data.list) ? data.list : (Array.isArray(data) ? data : [])
    }).catch(() => {
        activities.value = []
    }).finally(() => {
        loading.value = false
    })
}

function handleBuy(item) {
    if (!item || ordering.value) return
    ordering.value = true
    seckillOrder(item.id).then((data) => {
        const orderNo = (data && (data.order_no || (data.data && data.data.order_no) || (data.order && data.order.order_no))) || null
        if (orderNo) {
            router.push('/order/' + orderNo)
        } else {
            ElMessage.success('秒杀下单成功')
            router.push('/user/orders')
        }
    }).catch((err) => {
        ElMessage.error(err.response?.data?.msg || '秒杀失败，请稍后重试')
    }).finally(() => {
        ordering.value = false
        load()
    })
}

onMounted(() => {
    load()
    timer = setInterval(() => {
        now.value = Date.now()
    }, 1000)
})

onUnmounted(() => {
    if (timer) clearInterval(timer)
})
</script>

<style scoped>
.gradient-banner :deep(.el-card__body) {
    padding: 0;
}
.gradient-banner {
    background: linear-gradient(to right, #dc2626, #f97316);
    border: none;
}
.seckill-card :deep(.el-card__body) {
    padding: 16px;
}
</style>
