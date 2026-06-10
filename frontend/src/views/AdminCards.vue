<template>
    <AdminLayout page-title="卡密管理">
        <div class="space-y-4">
            <div class="flex flex-col md:flex-row md:items-center gap-3">
                <div class="flex items-center gap-2">
                    <select v-model.number="filter_product_id" class="form-input w-48">
                        <option :value="0">选择商品</option>
                        <option v-for="p in products" :key="p.id" :value="p.id">{{ p.name }}</option>
                    </select>
                    <button @click="loadCards" class="btn-primary btn-sm">加载卡密</button>
                </div>
                <div class="flex items-center gap-2 ml-auto">
                    <input type="file" ref="fileInput" @change="onFileChange" class="hidden">
                    <button @click="$refs.fileInput.click()" class="btn-primary">批量导入卡密</button>
                    <button @click="showTextModal = true" class="btn-secondary">文本导入</button>
                </div>
            </div>

            <div v-if="loading" class="card p-8 text-center text-gray-500">
                加载中...
            </div>
            <div v-else class="card">
                <div class="card-body overflow-x-auto">
                    <table class="table w-full">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>卡密</th>
                                <th>状态</th>
                                <th>创建时间</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="cards.length === 0">
                                <td colspan="4" class="text-center text-gray-500 py-8">请选择商品加载卡密</td>
                            </tr>
                            <tr v-for="item in cards" :key="item.id">
                                <td>{{ item.id }}</td>
                                <td class="font-mono text-sm">{{ item.card_data }}</td>
                                <td>
                                    <span v-if="item.status === 1" class="badge-green">已使用</span>
                                    <span v-else class="badge-gray">未使用</span>
                                </td>
                                <td>{{ item.created_at || '-' }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center gap-2">
                <button @click="prevPage" :disabled="page <= 1" class="btn-sm btn-secondary">上一页</button>
                <span class="text-sm text-gray-600">第 {{ page }} / {{ totalPages }} 页</span>
                <button @click="nextPage" :disabled="page >= totalPages" class="btn-sm btn-secondary">下一页</button>
            </div>

            <div v-if="showTextModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-20 p-4">
                <div class="card w-full max-w-2xl">
                    <div class="card-header flex items-center justify-between">
                        <span class="font-semibold">文本导入卡密（每行一条）</span>
                        <button @click="showTextModal = false" class="btn-sm btn-secondary">关闭</button>
                    </div>
                    <div class="card-body space-y-4">
                        <textarea v-model="importText" class="form-input" rows="10" placeholder="每行一个卡密"></textarea>
                        <div class="flex items-center justify-between">
                            <span class="text-sm text-gray-500">选择商品ID: {{ filter_product_id }}</span>
                            <button @click="importTextCards" :disabled="importing" class="btn-primary">
                                {{ importing ? '导入中...' : '导入' }}
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminCardList, adminCardImport, adminProductList } from '../api/admin'

const cards = ref([])
const products = ref([])
const filter_product_id = ref(0)
const page = ref(1)
const pageSize = 20
const pagination = reactive({ total: 0, size: 20 })
const loading = ref(false)
const showTextModal = ref(false)
const importText = ref('')
const importing = ref(false)
const fileInput = ref(null)

const totalPages = computed(() => Math.ceil(pagination.total / pagination.size) || 1)

function loadProducts() {
    adminProductList(1, 100, '').then((data) => {
        products.value = (data && data.list) ? data.list : []
    }).catch(() => {})
}

function loadCards() {
    loading.value = true
    adminCardList(page.value, pageSize, filter_product_id.value).then((data) => {
        cards.value = (data && data.list) ? data.list : []
        pagination.total = (data && data.total) ? data.total : 0
        pagination.size = (data && data.size) ? data.size : pageSize
    }).catch(() => {
        cards.value = []
        pagination.total = 0
    }).finally(() => {
        loading.value = false
    })
}

function onFileChange(e) {
    const file = e.target.files[0]
    if (!file) return
    if (!filter_product_id.value) {
        alert('请先选择商品')
        return
    }
    adminCardImport(filter_product_id.value, file).then(() => {
        alert('导入成功')
        loadCards()
    }).catch(() => {})
    e.target.value = ''
}

function importTextCards() {
    if (!importText.value.trim()) {
        alert('请输入卡密内容')
        return
    }
    if (!filter_product_id.value) {
        alert('请先选择商品')
        return
    }
    importing.value = true
    const blob = new Blob([importText.value], { type: 'text/plain' })
    const file = new File([blob], 'cards.txt', { type: 'text/plain' })
    adminCardImport(filter_product_id.value, file).then(() => {
        alert('导入成功')
        showTextModal.value = false
        importText.value = ''
        loadCards()
    }).catch(() => {}).finally(() => {
        importing.value = false
    })
}

function prevPage() {
    if (page.value > 1) {
        page.value--
        loadCards()
    }
}

function nextPage() {
    if (page.value < totalPages.value) {
        page.value++
        loadCards()
    }
}

onMounted(() => {
    loadProducts()
    loadCards()
})
</script>
