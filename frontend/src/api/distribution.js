import request from '../utils/request'

export function distributionSummary() {
    return request.get('/api/v1/distribution/summary')
}

export function distributionTeam(page, size) {
    const params = { page: page || 1, size: size || 20 }
    return request.get('/api/v1/distribution/team', { params })
}

export function distributionCommissions(page, size) {
    const params = { page: page || 1, size: size || 20 }
    return request.get('/api/v1/user/commissions', { params })
}

export function distributionPoster() {
    return request.get('/api/v1/distribution/poster')
}

export function withdrawApply(amount, account_type, account) {
    return request.post('/api/v1/withdraw/apply', { amount, account_type, account })
}

// Admin
export function adminDistributionGet() {
    return request.get('/admin/api/distribution/settings')
}

export function adminDistributionSet(settings) {
    return request.post('/admin/api/distribution/settings', settings)
}
