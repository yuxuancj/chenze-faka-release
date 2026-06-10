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
                <div>
                    <label class="form-label">联系方式</label>
                    <input v-model="settings.contact" type="text" class="form-input" placeholder="联系方式">
                </div>
                <div class="flex items-center space-x-2">
                    <button @click="save" :disabled="loading" class="btn-primary">
                        {{ loading ? '保存中...' : '保存设置' }}
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

const loading = ref(false)
const settings = reactive({
    site_name: '',
    site_desc: '',
    contact: ''
})

function load() {
    adminSettingsGet().then((res) => {
        const data = res.data || {}
        settings.site_name = data.site_name || ''
        settings.site_desc = data.site_desc || ''
        settings.contact = data.contact || ''
    }).catch(() => {})
}

function save() {
    loading.value = true
    adminSettingsSet({
        site_name: settings.site_name,
        site_desc: settings.site_desc,
        contact: settings.contact
    }).then(() => {
        alert('保存成功')
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

onMounted(load)
</script>
