<template>
    <AdminLayout page-title="积分设置">
        <div class="space-y-4">
            <div class="card">
                <div class="card-header font-semibold">积分规则</div>
                <div class="card-body space-y-3">
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                        <div>
                            <label class="form-label">消费1元获得积分</label>
                            <input v-model.number="settings.points_per_yuan" type="number" step="1" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">积分抵扣比例(100积分=多少元)</label>
                            <input v-model.number="settings.yuan_per_hundred" type="number" step="0.01" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">每日签到奖励积分</label>
                            <input v-model.number="settings.signin_reward" type="number" step="1" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">连续签到奖励积分</label>
                            <input v-model.number="settings.continuous_reward" type="number" step="1" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">最高抵扣比例(%)</label>
                            <input v-model.number="settings.max_discount_rate" type="number" step="1" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">功能开关</label>
                            <select v-model.number="settings.enabled" class="form-input">
                                <option :value="1">启用积分</option>
                                <option :value="0">禁用积分</option>
                            </select>
                        </div>
                    </div>
                    <button @click="save" :disabled="saving" class="btn-primary btn-sm">
                        {{ saving ? '保存中...' : '保存设置' }}
                    </button>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminPointsGet, adminPointsSet } from '../api/points'

const settings = reactive({
    points_per_yuan: 1, yuan_per_hundred: 1, signin_reward: 10,
    continuous_reward: 5, max_discount_rate: 50, enabled: 1
})
const saving = ref(false)
const loading = ref(true)

function load() {
    loading.value = true
    adminPointsGet().then((d) => {
        if (d) {
            settings.points_per_yuan = d.points_per_yuan || d.points_per_rmb || 1
            settings.yuan_per_hundred = d.yuan_per_hundred || d.deduct_rate || 1
            settings.signin_reward = d.signin_reward || d.sign_in_points || 10
            settings.continuous_reward = d.continuous_reward || 5
            settings.max_discount_rate = d.max_discount_rate || d.max_rate || 50
            settings.enabled = d.enabled === undefined ? 1 : d.enabled
        }
    }).catch(() => {}).finally(() => { loading.value = false })
}

function save() {
    saving.value = true
    adminPointsSet({ ...settings }).then(() => {
        alert('保存成功')
    }).catch(() => {}).finally(() => { saving.value = false })
}

onMounted(load)
</script>
