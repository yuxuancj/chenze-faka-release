import request from '../utils/request'

export function couponList() {
    return request.get('/api/v1/coupons')
}

export function couponRedeem(code) {
    return request.post('/api/v1/coupon/redeem', { code })
}

export function userCoupons() {
    return request.get('/api/v1/user/coupons')
}

// Admin
export function adminCouponList(page, size) {
    const params = { page: page || 1, size: size || 20 }
    return request.get('/admin/api/coupons', { params })
}

export function adminCouponCreate(data) {
    return request.post('/admin/api/coupons', data)
}

export function adminCouponUpdate(id, data) {
    return request.put('/admin/api/coupons/' + id, data)
}

export function adminCouponDelete(id) {
    return request.delete('/admin/api/coupons/' + id)
}

export function adminCouponGenerateCodes(id, count) {
    return request.post('/admin/api/coupons/' + id + '/codes', { count })
}
