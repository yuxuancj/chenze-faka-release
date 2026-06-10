<template>
    <AdminLayout page-title="系统设置">
        <div class="card max-w-2xl">
            <div class="card-header font-semibold">系统配置</div>
            <div class="card-body space-y-4">
                <div>
                    <label class="form-label">站点名称</label>
                    <input v-model="settings.site_name" type="text" class="form-input" placeholder="站点名称">
                </div>
                <div>
                    <label class="form-label">站点描述</label>
                    <textarea v-model="settings.site_desc" class="form-input" rows="3" placeholder="站点描述"></textarea>
                </div>
                <div class="flex items-center gap-2 pt-2">
                    <button @click="save" :disabled="saving" class="btn-primary">
                        {{ saving ? '保存中...' : '保存设置' }}
                    </button>
                    <button @click="load" class="btn-secondary">重新加载</button>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminSettingsGet, adminSettingsSet } from '../api/admin'

const saving = ref(false)
const settings = reactive({
    site_name: '',
    site_desc: ''
})

function load() {
    adminSettingsGet().then((data) => {
        if (data) {
            settings.site_name = data.site_name || ''
            settings.site_desc = data.site_desc || ''
        }
    }).catch(() => {})
}

function save() {
    saving.value = true
    adminSettingsSet({
        site_name: settings.site_name,
        site_desc: settings.site_desc
    }).then(() => {
        alert('保存成功')
    }).catch(() => {}).finally(() => {
        saving.value = false
    })
}

onMounted(load)
</script>
