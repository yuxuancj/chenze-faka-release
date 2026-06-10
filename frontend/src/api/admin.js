import request from '../utils/request'

export function adminDashboard() {
    return request.get('/admin/api/dashboard')
}

export function adminProductList(page, size, keyword) {
    const params = { page: page || 1, size: size || 20, keyword: keyword || '' }
    return request.get('/admin/api/products', { params })
}

export function adminProductCreate(data) {
    return request.post('/admin/api/products', data)
}

export function adminProductUpdate(id, data) {
    return request.put('/admin/api/products/' + id, data)
}

export function adminProductDelete(id) {
    return request.delete('/admin/api/products/' + id)
}

export function adminCardList(page, size, product_id) {
    const params = { page: page || 1, size: size || 20, product_id: product_id || 0 }
    return request.get('/admin/api/cards', { params })
}

export function adminCardImport(product_id, cards) {
    const formData = new FormData()
    formData.append('product_id', product_id)
    if (cards instanceof File) {
        formData.append('cards', cards)
    } else {
        formData.append('cards', cards)
    }
    return request.post('/admin/api/cards/import', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
    })
}

export function adminCategoryList() {
    return request.get('/admin/api/categories')
}

export function adminCategoryCreate(data) {
    return request.post('/admin/api/categories', data)
}

export function adminCategoryUpdate(id, data) {
    return request.put('/admin/api/categories/' + id, data)
}

export function adminCategoryDelete(id) {
    return request.delete('/admin/api/categories/' + id)
}

export function adminOrderList(page, size, keyword) {
    const params = { page: page || 1, size: size || 20, keyword: keyword || '' }
    return request.get('/admin/api/orders', { params })
}

export function adminOrderDetail(id) {
    return request.get('/admin/api/orders/' + id)
}

export function adminUserList(page, size) {
    const params = { page: page || 1, size: size || 20 }
    return request.get('/admin/api/users', { params })
}

export function adminUserUpdate(id, data) {
    return request.put('/admin/api/users/' + id, data)
}

export function adminSettingsGet() {
    return request.get('/admin/api/settings')
}

export function adminSettingsSet(settings) {
    return request.post('/admin/api/settings', settings)
}
