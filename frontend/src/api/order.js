import request from '../utils/request'

export function createOrder(params) {
    return request.post('/order/create', params)
}

export function orderList(page, size) {
    return request.post('/order/list', { page, size })
}

export function orderDetail(order_no) {
    return request.post('/order/detail', { order_no })
}

export function payOrder(order_no) {
    return request.post('/order/pay', { order_no })
}
