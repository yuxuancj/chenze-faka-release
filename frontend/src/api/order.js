import request from '../utils/request'

export function createOrder(params) {
    return request.post('/api/v1/orders', params)
}

export function orderList(page, size) {
    const params = { page: page || 1, size: size || 20 }
    return request.get('/api/v1/orders', { params })
}

export function orderDetail(order_no) {
    return request.get('/api/v1/orders/' + order_no)
}

export function payOrder(order_no) {
    return request.post('/api/v1/pay', { order_no })
}
