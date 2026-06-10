import request from '../utils/request'

export function register(data) {
    return request.post('/api/v1/user/register', data)
}

export function login(email, password) {
    return request.post('/api/v1/user/login', { email, password })
}

export function profile() {
    return request.get('/api/v1/user/profile')
}

export function updateProfile(nickname) {
    return request.post('/api/v1/user/profile', { nickname })
}

export function changePassword(old_password, new_password) {
    return request.post('/api/v1/user/password', { old_password, new_password })
}
