<template>
    <Layout>
        <div class="space-y-6">
            <div class="bg-gradient-to-r from-red-600 to-orange-500 rounded-lg overflow-hidden">
                <div class="px-6 py-8 md:px-12 md:py-10">
                    <h1 class="text-2xl md:text-3xl font-bold text-white mb-2">限时秒杀</h1>
                    <p class="text-red-100">精选商品，限时低价，售完即止</p>
                </div>
            </div>

            <div v-if="loading" class="card p-8 text-center text-gray-500">
                加载中...
            </div>
            <div v-else-if="activities.length === 0" class="card p-8 text-center text-gray-500">
                暂无秒杀活动
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <div v-for="item in activities" :key="item.id" class="card overflow-hidden">
                    <div class="h-40 bg-gradient-to-br from-red-50 to-orange-100 flex items-center justify-center text-gray-500 text-sm">
                        {{ item.product_name || '秒杀商品' }}
                    </div>
                    <div class="card-body space-y-3">
                        <h3 class="font-semibold text-gray-800 truncate">{{ item.product_name || '秒杀商品' }}</h3>
                        <div class="flex items-baseline gap-2">
                            <span class="text-2xl font-bold text-red-600">￥{{ item.seckill_price }}</span>
                            <span class="text-sm text-gray-400 line-through">￥{{ item.original_price }}</span>
                        </div>
                        <div class="flex items-center gap-4 text-xs text-gray-600">
                            <span>每人限购 {{ item.limit_per_user }} 件</span>
                            <span>秒杀库存: {{ item.stock }}</span>
                        </div>
                        <div v-if="item.status === 'ongoing'" class="text-xs text-green-600 font-semibold">
                            进行中 · 剩余: {{ countdownFor(item) }}
                        </div>
                        <div v-else-if="item.status === 'upcoming'" class="text-xs text-yellow-600 font-semibold">
                            即将开始 · {{ formatDateTime(item.start_time) }}
                        </div>
                        <div v-else-if="item.status === 'ended'" class="text-xs text-gray-500">
                            已结束
                        </div>
                        <div class="w-full bg-gray-200 rounded-full h-2">
                            <div class="bg-red-500 h-2 rounded-full" :style="{ width: stockProgress(item) + '%' }"></div>
                        </div>
                        <button v-if="item.status === 'ongoing'" @click="handleBuy(item)" :disabled="ordering" class="btn-primary w-full bg-red-600 hover:bg-red-700">
                            {{ ordering ? '提交中...' : '立即抢购' }}
                        </button>
                        <button v-else-if="item.status === 'upcoming'" disabled class="btn w-full bg-gray-300 text-gray-600 cursor-not-allowed">
                            尚未开始
                        </button>
                        <button v-else disabled class="btn w-full bg-gray-300 text-gray-600 cursor-not-allowed">
                            已结束
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'
import { seckillActivities, seckillOrder } from '../api/seckill'

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
            alert('秒杀下单成功')
            router.push('/user/orders')
        }
    }).catch((err) => {
        alert(err.response?.data?.msg || '秒杀失败，请稍后重试')
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
