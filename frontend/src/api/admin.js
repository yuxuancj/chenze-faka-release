import request from '../utils/request'

export function adminDashboard() {
    return request.post('/admin/dashboard')
}

export function adminProductList(page, size, category_id, keyword) {
    return request.post('/admin/product/list', { page, size, category_id, keyword })
}

export function adminProductCreate(data) {
    return request.post('/admin/product/create', data)
}

export function adminProductUpdate(id, data) {
    return request.post('/admin/product/update', { id, ...data })
}

export function adminProductDelete(id) {
    return request.post('/admin/product/delete', { id })
}

export function adminCardList(page, size, product_id) {
    return request.post('/admin/card/list', { page, size, product_id })
}

export function adminCardImport(product_id, cards) {
    const formData = new FormData()
    formData.append('product_id', product_id)
    formData.append('cards', cards)
    return request.post('/admin/card/import', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
    })
}

export function adminCategoryList(page, size) {
    return request.post('/admin/category/list', { page, size })
}

export function adminCategoryCreate(data) {
    return request.post('/admin/category/create', data)
}

export function adminCategoryUpdate(id, data) {
    return request.post('/admin/category/update', { id, ...data })
}

export function adminCategoryDelete(id) {
    return request.post('/admin/category/delete', { id })
}

export function adminOrderList(page, size) {
    return request.post('/admin/order/list', { page, size })
}

export function adminOrderDetail(id) {
    return request.post('/admin/order/detail', { id })
}

export function adminUserList(page, size) {
    return request.post('/admin/user/list', { page, size })
}

export function adminUserUpdate(id, data) {
    return request.post('/admin/user/update', { id, ...data })
}

export function adminSettingsGet() {
    return request.post('/admin/settings/get')
}

export function adminSettingsSet(settings) {
    return request.post('/admin/settings/set', { settings })
}
