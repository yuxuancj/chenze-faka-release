import request from '../utils/request'

export function pointsLogs(page, size) {
    const params = { page: page || 1, size: size || 20 }
    return request.get('/api/v1/user/points', { params })
}

export function signIn() {
    return request.post('/api/v1/user/signin')
}

// Admin
export function adminPointsGet() {
    return request.get('/admin/api/points/settings')
}

export function adminPointsSet(settings) {
    return request.post('/admin/api/points/settings', settings)
}
