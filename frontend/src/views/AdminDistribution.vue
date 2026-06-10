<template>
    <AdminLayout page-title="分销设置">
        <div class="space-y-4">
            <div class="card">
                <div class="card-header font-semibold">佣金设置</div>
                <div class="card-body space-y-3">
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                        <div>
                            <label class="form-label">一级佣金比例（%）</label>
                            <input v-model.number="settings.level1_rate" type="number" step="0.01" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">二级佣金比例（%）</label>
                            <input v-model.number="settings.level2_rate" type="number" step="0.01" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">三级佣金比例（%）</label>
                            <input v-model.number="settings.level3_rate" type="number" step="0.01" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">最低提现金额</label>
                            <input v-model.number="settings.min_withdraw" type="number" step="0.01" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">功能开关</label>
                            <select v-model.number="settings.enabled" class="form-input">
                                <option :value="1">启用分销</option>
                                <option :value="0">禁用分销</option>
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
import { adminDistributionGet, adminDistributionSet } from '../api/distribution'

const settings = reactive({
    level1_rate: 10, level2_rate: 5, level3_rate: 2,
    min_withdraw: 10, enabled: 1
})
const saving = ref(false)
const loading = ref(true)

function load() {
    loading.value = true
    adminDistributionGet().then((d) => {
        if (d) {
            settings.level1_rate = d.level1_rate || d.level_1 || 10
            settings.level2_rate = d.level2_rate || d.level_2 || 5
            settings.level3_rate = d.level3_rate || d.level_3 || 2
            settings.min_withdraw = d.min_withdraw || d.min_amount || 10
            settings.enabled = d.enabled === undefined ? 1 : d.enabled
        }
    }).catch(() => {}).finally(() => { loading.value = false })
}

function save() {
    saving.value = true
    adminDistributionSet({ ...settings }).then(() => {
        alert('保存成功')
    }).catch(() => {}).finally(() => { saving.value = false })
}

onMounted(load)
</script>
