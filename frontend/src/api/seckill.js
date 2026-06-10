import request from '../utils/request'

export function seckillActivities() {
    return request.get('/api/v1/seckill/activities')
}

export function seckillOrder(product_id, sku_id, quantity) {
    return request.post('/api/v1/seckill/order', { product_id, sku_id, quantity })
}

// Admin
export function adminSeckillList(page, size) {
    const params = { page: page || 1, size: size || 20 }
    return request.get('/admin/api/seckills', { params })
}

export function adminSeckillCreate(data) {
    return request.post('/admin/api/seckills', data)
}

export function adminSeckillUpdate(id, data) {
    return request.put('/admin/api/seckills/' + id, data)
}

export function adminSeckillDelete(id) {
    return request.delete('/admin/api/seckills/' + id)
}
