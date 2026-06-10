<template>
    <Layout>
        <div class="max-w-md mx-auto space-y-6">
            <div class="bg-gradient-to-r from-green-500 to-emerald-600 rounded-lg overflow-hidden">
                <div class="px-6 py-8 text-center">
                    <h1 class="text-2xl font-bold text-white mb-2">每日签到</h1>
                    <p class="text-green-100">连续签到领取额外积分奖励</p>
                </div>
            </div>

            <div v-if="loading" class="card p-8 text-center text-gray-500">
                <div class="animate-pulse">
                    <div class="h-4 bg-gray-200 rounded w-1/2 mx-auto mb-4"></div>
                    <div class="h-4 bg-gray-200 rounded w-3/4 mx-auto"></div>
                </div>
            </div>
            <template v-else>
                <div class="card">
                    <div class="card-body text-center space-y-4">
                        <div class="text-7xl">{{ canSignin ? '🎯' : '✅' }}</div>
                        <div v-if="canSignin" class="space-y-4">
                            <p class="text-gray-600">
                                今日签到可获得 <span class="text-green-600 font-bold text-xl">{{ signinReward }}</span> 积分
                            </p>
                            <p v-if="continuousDays > 0" class="text-sm text-orange-500">
                                继续坚持！已连续签到 <strong>{{ continuousDays }}</strong> 天
                            </p>
                            <p class="text-xs text-gray-400">
                                连续3天额外+{{ continuousReward }}积分，连续7天双倍奖励
                            </p>
                            <button @click="doSignin" :disabled="signing"
                                class="btn-primary w-full py-4 text-lg rounded-xl shadow-lg hover:shadow-xl transition-all">
                                <span v-if="signing" class="flex items-center justify-center gap-2">
                                    <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
                                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
                                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                    </svg>
                                    签到中...
                                </span>
                                <span v-else>立即签到</span>
                            </button>
                        </div>
                        <div v-else class="space-y-3">
                            <p class="text-gray-600 text-lg">今日已签到</p>
                            <div class="inline-flex items-center gap-2 bg-green-50 text-green-700 px-4 py-2 rounded-full">
                                <span>🔥</span>
                                <span class="font-medium">连续签到 {{ continuousDays }} 天</span>
                            </div>
                            <p class="text-green-600 font-medium">明日再来领取更多积分</p>
                        </div>
                    </div>
                </div>

                <div class="card">
                    <div class="card-header font-semibold">签到规则</div>
                    <div class="card-body space-y-2 text-sm text-gray-600">
                        <div class="flex items-start gap-2">
                            <span class="text-green-500">✓</span>
                            <span>每日签到可获得 <strong>{{ signinReward }}</strong> 积分</span>
                        </div>
                        <div class="flex items-start gap-2">
                            <span class="text-orange-500">🔥</span>
                            <span>连续签到{{ continuousRewardDays }}天额外奖励 <strong>+{{ continuousReward }}</strong> 积分</span>
                        </div>
                        <div class="flex items-start gap-2">
                            <span class="text-purple-500">⭐</span>
                            <span>连续签到{{ doubleRewardDays }}天奖励<strong>翻倍</strong></span>
                        </div>
                        <div class="flex items-start gap-2">
                            <span class="text-blue-500">💰</span>
                            <span>积分可抵扣订单金额</span>
                        </div>
                    </div>
                </div>

                <div class="card">
                    <div class="card-header font-semibold flex justify-between items-center">
                        <span>我的积分</span>
                        <router-link to="/user/points" class="text-blue-500 text-sm font-normal">查看全部</router-link>
                    </div>
                    <div class="card-body">
                        <div class="text-center py-4">
                            <span class="text-5xl font-bold text-green-600">{{ balance || 0 }}</span>
                            <p class="text-sm text-gray-500 mt-1">可用积分</p>
                        </div>
                        <div class="border-t pt-4 space-y-2">
                            <p class="text-sm text-gray-500 mb-2">最近积分变动</p>
                            <div v-if="logsLoading" class="text-center text-gray-400 py-4 text-sm">加载中...</div>
                            <div v-else-if="logs.length === 0" class="text-center text-gray-400 py-4 text-sm">暂无记录</div>
                            <div v-else class="space-y-2">
                                <div v-for="log in logs.slice(0, 5)" :key="log.id"
                                    class="flex justify-between items-center text-sm py-2 border-b border-gray-100 last:border-0">
                                    <div>
                                        <span class="text-gray-700">{{ log.description || getTypeDesc(log.type) }}</span>
                                        <span v-if="log.created_at" class="text-xs text-gray-400 block">
                                            {{ formatTime(log.created_at) }}
                                        </span>
                                    </div>
                                    <span :class="log.amount > 0 ? 'text-green-600 font-medium' : 'text-red-600'">
                                        {{ log.amount > 0 ? '+' : '' }}{{ log.amount }}
                                    </span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </template>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { signIn } from '../api/points'
import Toast from '../utils/toast'

const loading = ref(true)
const logsLoading = ref(false)
const signing = ref(false)
const canSignin = ref(false)
const continuousDays = ref(0)
const balance = ref(0)
const logs = ref([])
const signinReward = ref(10)
const continuousReward = ref(5)
const continuousRewardDays = ref(3)
const doubleRewardDays = ref(7)

function getTypeDesc(type) {
    const map = {
        'signin': '每日签到',
        'order': '购物返积分',
        'deduct': '积分抵扣',
        'refund': '退款返还'
    }
    return map[type] || '积分变动'
}

function formatTime(timeStr) {
    if (!timeStr) return ''
    const d = new Date(timeStr)
    const now = new Date()
    const diff = now - d
    if (diff < 60000) return '刚刚'
    if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
    if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
    return d.toLocaleDateString()
}

async function loadStatus() {
    loading.value = true
    try {
        const resp = await fetch('/api/v1/user/signin/status')
        const d = await resp.json()
        if (d.code === 0) {
            canSignin.value = !!d.data?.can_signin
            continuousDays.value = d.data?.continuous_days || 0
            balance.value = d.data?.balance || 0
        } else {
            canSignin.value = true
        }
    } catch (e) {
        console.error('加载签到状态失败:', e)
        canSignin.value = true
    } finally {
        loading.value = false
    }
}

async function loadLogs() {
    logsLoading.value = true
    try {
        const resp = await fetch('/api/v1/user/points/logs?page=1&size=10')
        const d = await resp.json()
        if (d.code === 0) {
            logs.value = (d.data?.list) ? d.data.list : []
        }
    } catch (e) {
        console.error('加载积分记录失败:', e)
    } finally {
        logsLoading.value = false
    }
}

async function doSignin() {
    signing.value = true
    try {
        const result = await signIn()
        const data = result || {}
        const earned = data.points || signinReward.value
        const days = data.continuous_days || (continuousDays.value + 1)

        Toast.success(`签到成功！获得 ${earned} 积分`)
        canSignin.value = false
        continuousDays.value = days
        balance.value += earned

        // 重新加载积分记录
        await loadLogs()
    } catch (e) {
        const msg = e.message || ''
        if (msg.includes('已签到')) {
            Toast.info('今日已签到')
            canSignin.value = false
        } else if (msg.includes('登录')) {
            Toast.error('请先登录')
            window.location.href = '/user/login'
        } else {
            Toast.error(msg || '签到失败，请重试')
        }
    } finally {
        signing.value = false
    }
}

onMounted(() => {
    loadStatus()
    loadLogs()
})
</script>
