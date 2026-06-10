<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">我的积分</h2>

            <div class="card bg-gradient-to-r from-purple-50 to-indigo-50">
                <div class="card-body">
                    <div class="flex items-center justify-between">
                        <div>
                            <div class="text-sm text-gray-500">当前积分</div>
                            <div class="text-4xl font-bold text-purple-700 mt-2">{{ currentPoints }}</div>
                        </div>
                        <button @click="handleSignin" :disabled="signingIn || signedToday" class="btn-primary">
                            <template v-if="signingIn">签到中...</template>
                            <template v-else-if="signedToday">今日已签到</template>
                            <template v-else>每日签到</template>
                        </button>
                    </div>
                    <div v-if="signinReward" class="text-sm text-green-600 mt-3">
                        签到奖励: +{{ signinReward }} 积分
                    </div>
                </div>
            </div>

            <div class="card">
                <div class="card-header font-semibold">积分明细</div>
                <div class="card-body overflow-x-auto">
                    <table class="table w-full">
                        <thead>
                            <tr>
                                <th>类型</th>
                                <th>变动</th>
                                <th>备注</th>
                                <th>时间</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="logs.length === 0">
                                <td colspan="4" class="text-center text-gray-500 py-8">暂无积分记录</td>
                            </tr>
                            <tr v-for="log in logs" :key="log.id">
                                <td>{{ log.type_name || (log.type === 1 ? '获取' : (log.type === 2 ? '消耗' : '其他')) }}</td>
                                <td :class="(log.change || log.points || 0) >= 0 ? 'text-green-600' : 'text-red-600'" class="font-semibold">
                                    {{ (log.change || log.points || 0) > 0 ? '+' : '' }}{{ log.change || log.points || 0 }}
                                </td>
                                <td class="text-gray-600">{{ log.remark || log.description || '-' }}</td>
                                <td>{{ log.created_at || '-' }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { pointsLogs, signIn } from '../api/points'
import { profile } from '../api/user'

const logs = ref([])
const currentPoints = ref(0)
const signedToday = ref(false)
const signinReward = ref(0)
const signingIn = ref(false)
const loading = ref(true)

function loadAll() {
    loading.value = true
    Promise.all([
        profile().then((d) => {
            currentPoints.value = (d && d.points) ? d.points : 0
        }).catch(() => {}),
        pointsLogs(1, 30).then((d) => {
            logs.value = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
        }).catch(() => {})
    ]).finally(() => {
        loading.value = false
    })
}

function handleSignin() {
    if (signedToday.value) return
    signingIn.value = true
    signIn().then((d) => {
        signinReward.value = (d && d.points) ? d.points : ((d && d.reward) ? d.reward : 10)
        signedToday.value = true
        currentPoints.value += signinReward.value
        loadAll()
    }).catch(() => {}).finally(() => {
        signingIn.value = false
    })
}

onMounted(loadAll)
</script>
