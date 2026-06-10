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
                加载中...
            </div>
            <template v-else>
                <div class="card">
                    <div class="card-body text-center space-y-4">
                        <div class="text-6xl">{{ canSignin ? '🎯' : '✅' }}</div>
                        <div v-if="canSignin" class="space-y-4">
                            <p class="text-gray-600">今日签到可获得 <span class="text-green-600 font-bold">{{ signinReward }}</span> 积分</p>
                            <p v-if="continuousDays > 0" class="text-sm text-orange-600">
                                连续签到 {{ continuousDays }} 天！
                            </p>
                            <p class="text-xs text-gray-500">
                                连续3天额外+{{ continuousReward }}积分，连续7天双倍奖励
                            </p>
                            <button @click="doSignin" :disabled="signing" class="btn-primary w-full py-3 text-lg">
                                {{ signing ? '签到中...' : '立即签到' }}
                            </button>
                        </div>
                        <div v-else class="space-y-2">
                            <p class="text-gray-600">今日已签到</p>
                            <p class="text-sm text-gray-500">连续签到 {{ continuousDays }} 天</p>
                            <p class="text-green-600 font-medium">明日再来领取更多积分</p>
                        </div>
                    </div>
                </div>

                <div class="card">
                    <div class="card-header font-semibold">签到规则</div>
                    <div class="card-body space-y-2 text-sm text-gray-600">
                        <p>• 每日签到可获得 {{ signinReward }} 积分</p>
                        <p>• 连续签到{{ continuousRewardDays }}天额外奖励 {{ continuousReward }} 积分</p>
                        <p>• 连续签到{{ doubleRewardDays }}天奖励翻倍</p>
                        <p>• 积分可抵扣订单金额</p>
                    </div>
                </div>

                <div class="card">
                    <div class="card-header font-semibold">我的积分</div>
                    <div class="card-body">
                        <div class="text-center">
                            <span class="text-4xl font-bold text-green-600">{{ balance || 0 }}</span>
                            <p class="text-sm text-gray-500 mt-1">可用积分</p>
                        </div>
                        <div class="mt-4 space-y-2">
                            <div class="flex justify-between text-sm">
                                <span class="text-gray-500">最近积分变动</span>
                                <router-link to="/user/points" class="text-blue-500 text-sm">查看全部</router-link>
                            </div>
                            <div v-if="logs.length === 0" class="text-center text-gray-400 py-4">暂无记录</div>
                            <div v-else class="space-y-2">
                                <div v-for="log in logs.slice(0, 5)" :key="log.id" class="flex justify-between text-sm">
                                    <span class="text-gray-600">{{ log.description || getTypeDesc(log.type) }}</span>
                                    <span :class="log.amount > 0 ? 'text-green-600' : 'text-red-600'">
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
import { signIn, pointsLogs } from '../api/points'
import Toast from '../utils/toast'

const loading = ref(true)
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

function loadStatus() {
    loading.value = true
    fetch('/api/v1/user/signin/status').then(r => r.json()).then((d) => {
        if (d.code === 0) {
            canSignin.value = d.data.can_signin
            continuousDays.value = d.data.continuous_days || 0
            balance.value = d.data.balance || 0
        }
        loadLogs()
    }).catch(() => {
        canSignin.value = true
    }).finally(() => {
        loading.value = false
    })
}

function loadLogs() {
    fetch('/api/v1/user/points/logs').then(r => r.json()).then((d) => {
        if (d.code === 0) {
            logs.value = (d.data && d.data.list) ? d.data.list : []
        }
    }).catch(() => {})
}

function doSignin() {
    signing.value = true
    signIn().then((d) => {
        const data = d || {}
        const earned = data.points || signinReward.value
        const days = data.continuous_days || (continuousDays.value + 1)
        Toast.success(`签到成功！获得 ${earned} 积分，连续签到 ${days} 天`)
        canSignin.value = false
        continuousDays.value = days
        balance.value += earned
        loadLogs()
    }).catch((e) => {
        const msg = e.message || ''
        if (msg.includes('已签到')) {
            Toast.info('今日已签到')
            canSignin.value = false
        } else {
            Toast.error(msg || '签到失败，请重试')
        }
    }).finally(() => {
        signing.value = false
    })
}

onMounted(() => {
    loadStatus()
})
</script>
