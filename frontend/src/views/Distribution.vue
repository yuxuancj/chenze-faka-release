<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">分销中心</h2>

            <div v-if="loading" class="card p-8 text-center text-gray-500">
                加载中...
            </div>
            <template v-else>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div class="card">
                        <div class="card-body">
                            <div class="text-sm text-gray-500">累计佣金</div>
                            <div class="text-2xl font-bold text-gray-800 mt-2">￥{{ summary.total_commission || 0 }}</div>
                        </div>
                    </div>
                    <div class="card">
                        <div class="card-body">
                            <div class="text-sm text-gray-500">可提现金额</div>
                            <div class="text-2xl font-bold text-blue-600 mt-2">￥{{ summary.available_commission || 0 }}</div>
                        </div>
                    </div>
                    <div class="card">
                        <div class="card-body">
                            <div class="text-sm text-gray-500">下级人数</div>
                            <div class="text-2xl font-bold text-gray-800 mt-2">{{ summary.team_count || 0 }}</div>
                        </div>
                    </div>
                </div>

                <div class="card">
                    <div class="card-header font-semibold">推广链接</div>
                    <div class="card-body space-y-3">
                        <div class="flex items-center gap-2">
                            <input :value="promoLink" readonly class="form-input flex-1 font-mono text-xs">
                            <button @click="copyLink" class="btn-primary btn-sm">复制</button>
                        </div>
                        <div v-if="posterUrl" class="mt-3">
                            <div class="text-sm text-gray-600 mb-2">推广海报</div>
                            <img :src="posterUrl" class="max-w-xs w-full border border-gray-200 rounded-md">
                        </div>
                        <button v-if="!posterLoading && !posterUrl" @click="loadPoster" class="btn-secondary btn-sm">
                            生成海报
                        </button>
                        <div v-if="posterLoading" class="text-sm text-gray-500">生成中...</div>
                    </div>
                </div>

                <div class="card">
                    <div class="card-header font-semibold">佣金明细</div>
                    <div class="card-body overflow-x-auto">
                        <table class="table w-full">
                            <thead>
                                <tr>
                                    <th>来源</th>
                                    <th>金额</th>
                                    <th>等级</th>
                                    <th>状态</th>
                                    <th>时间</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-if="commissions.length === 0">
                                    <td colspan="5" class="text-center text-gray-500 py-8">暂无佣金记录</td>
                                </tr>
                                <tr v-for="c in commissions" :key="c.id">
                                    <td>{{ c.from_user_nickname || ('用户' + (c.from_user_id || '')) }}</td>
                                    <td class="text-blue-600 font-semibold">￥{{ c.amount }}</td>
                                    <td>{{ c.level }}</td>
                                    <td>
                                        <span v-if="c.status === 1" class="badge-green">已结算</span>
                                        <span v-else class="badge-yellow">待结算</span>
                                    </td>
                                    <td>{{ c.created_at || c.settled_at || '-' }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>

                <div class="card">
                    <div class="card-header font-semibold">申请提现</div>
                    <div class="card-body space-y-3">
                        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                            <div>
                                <label class="form-label">提现金额</label>
                                <input v-model.number="withdrawForm.amount" type="number" step="0.01" class="form-input" placeholder="可提现: ￥{{ summary.available_commission || 0 }}">
                            </div>
                            <div>
                                <label class="form-label">收款方式</label>
                                <select v-model="withdrawForm.account_type" class="form-input">
                                    <option value="alipay">支付宝</option>
                                    <option value="wechat">微信</option>
                                    <option value="bank">银行卡</option>
                                </select>
                            </div>
                            <div>
                                <label class="form-label">收款账号</label>
                                <input v-model="withdrawForm.account" type="text" class="form-input" placeholder="请输入收款账号">
                            </div>
                        </div>
                        <button @click="applyWithdraw" :disabled="withdrawing" class="btn-primary">
                            {{ withdrawing ? '提交中...' : '提交申请' }}
                        </button>
                    </div>
                </div>
            </template>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import Layout from '../components/Layout.vue'
import { distributionSummary, distributionCommissions, distributionPoster, withdrawApply } from '../api/distribution'

const summary = ref({})
const commissions = ref([])
const loading = ref(true)
const withdrawing = ref(false)
const posterUrl = ref('')
const posterLoading = ref(false)
const withdrawForm = reactive({ amount: 0, account_type: 'alipay', account: '' })

const promoLink = computed(() => {
    const code = summary.value.invite_code || ''
    const base = window.location.origin + '/register?invite='
    return base + code
})

function loadAll() {
    loading.value = true
    Promise.all([
        distributionSummary().then((d) => { summary.value = d || {} }).catch(() => {}),
        distributionCommissions(1, 20).then((d) => {
            commissions.value = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
        }).catch(() => {})
    ]).finally(() => {
        loading.value = false
    })
}

function loadPoster() {
    posterLoading.value = true
    distributionPoster().then((d) => {
        posterUrl.value = (d && (d.poster_url || d.url || d.image)) || ''
    }).catch(() => {}).finally(() => {
        posterLoading.value = false
    })
}

function copyLink() {
    const ta = document.createElement('textarea')
    ta.value = promoLink.value
    document.body.appendChild(ta)
    ta.select()
    try {
        document.execCommand('copy')
        alert('已复制推广链接')
    } catch (e) {
        alert('复制失败，请手动复制')
    }
    document.body.removeChild(ta)
}

function applyWithdraw() {
    if (!withdrawForm.amount || withdrawForm.amount <= 0) {
        alert('请输入正确的提现金额')
        return
    }
    if (!withdrawForm.account) {
        alert('请输入收款账号')
        return
    }
    withdrawing.value = true
    withdrawApply(withdrawForm.amount, withdrawForm.account_type, withdrawForm.account).then(() => {
        alert('申请已提交，等待审核')
        withdrawForm.amount = 0
        withdrawForm.account = ''
    }).catch(() => {}).finally(() => {
        withdrawing.value = false
    })
}

onMounted(loadAll)
</script>
