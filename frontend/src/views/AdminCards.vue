<template>
    <AdminLayout page-title="卡密管理">
        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <select v-model="product_id" class="form-input w-48">
                        <option :value="0">选择商品</option>
                        <option v-for="p in products" :key="p.id" :value="p.id">{{ p.name }}</option>
                    </select>
                    <button @click="load" class="btn-primary btn-sm">加载</button>
                </div>
                <div class="flex items-center space-x-2">
                    <input v-model="product_id" type="number" placeholder="商品ID" class="form-input w-24">
                    <input type="file" ref="fileInput" @change="onFileChange" class="hidden">
                    <button @click="$refs.fileInput.click()" class="btn-primary">批量导入卡密</button>
                    <button @click="showTextModal = true" class="btn-secondary">文本导入</button>
                </div>
            </div>

            <div class="card">
                <div class="card-body">
                    <table class="table">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>卡号</th>
                                <th>密码</th>
                                <th>商品</th>
                                <th>状态</th>
                                <th>创建时间</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="cards.length === 0">
                                <td colspan="6" class="text-center text-gray-500 py-8">请选择商品加载卡密</td>
                            </tr>
                            <tr v-for="item in cards" :key="item.id">
                                <td>{{ item.id }}</td>
                                <td>{{ item.card_no }}</td>
                                <td>{{ item.card_pwd }}</td>
                                <td>{{ item.product_name || '-' }}</td>
                                <td>
                                    <span v-if="item.status === 1" class="badge-green">已使用</span>
                                    <span v-else class="badge-gray">未使用</span>
                                </td>
                                <td>{{ item.created_at }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center space-x-2">
                <button @click="prevPage" :disabled="page <= 1" class="btn-sm btn-secondary">上一页</button>
                <span class="text-sm text-gray-600">第 {{ page }} / {{ totalPages }} 页</span>
                <button @click="nextPage" :disabled="page >= totalPages" class="btn-sm btn-secondary">下一页</button>
            </div>

            <div v-if="showTextModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-20">
                <div class="card w-full max-w-2xl">
                    <div class="card-header flex items-center justify-between">
                        <span class="font-semibold">文本导入卡密（每行一条，卡号和密码用逗号或空格分隔）</span>
                        <button @click="showTextModal = false" class="btn-sm btn-secondary">关闭</button>
                    </div>
                    <div class="card-body space-y-4">
                        <textarea v-model="importText" class="form-input" rows="10" placeholder="卡号,密码&#10;卡号 密码"></textarea>
                        <div class="flex items-center justify-between">
                            <span class="text-sm text-gray-500">选择商品ID: {{ product_id }}</span>
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
const product_id = ref(0)
const page = ref(1)
const pageSize = 20
const pagination = ref({ total: 0, size: 20 })
const showTextModal = ref(false)
const importText = ref('')
const importing = ref(false)
const fileInput = ref(null)

const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.size) || 1)

function loadProducts() {
    adminProductList(1, 100, 0, '').then((res) => {
        products.value = res.data && res.data.list ? res.data.list : []
    }).catch(() => {})
}

function load() {
    adminCardList(page.value, pageSize, product_id.value || 0).then((res) => {
        cards.value = res.data && res.data.list ? res.data.list : []
        pagination.value.total = (res.data && res.data.total) ? res.data.total : 0
    }).catch(() => {})
}

function onFileChange(e) {
    const file = e.target.files[0]
    if (!file) return
    if (!product_id.value) {
        alert('请先选择商品')
        return
    }
    adminCardImport(product_id.value, file).then(() => {
        alert('导入成功')
        load()
    }).catch(() => {})
}

function importTextCards() {
    if (!importText.value.trim()) {
        alert('请输入卡密内容')
        return
    }
    if (!product_id.value) {
        alert('请选择商品')
        return
    }
    importing.value = true
    const blob = new Blob([importText.value], { type: 'text/plain' })
    const file = new File([blob], 'cards.txt', { type: 'text/plain' })
    adminCardImport(product_id.value, file).then(() => {
        alert('导入成功')
        showTextModal.value = false
        importText.value = ''
        load()
    }).catch(() => {}).finally(() => {
        importing.value = false
    })
}

function prevPage() {
    if (page.value > 1) {
        page.value--
        load()
    }
}

function nextPage() {
    if (page.value < totalPages.value) {
        page.value++
        load()
    }
}

onMounted(loadProducts)
</script>
