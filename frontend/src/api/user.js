import request from '../utils/request'

export function register(email, password, nickname) {
    return request.post('/user/register', { email, password, nickname })
}

export function login(email, password) {
    return request.post('/user/login', { email, password })
}

export function profile() {
    return request.post('/user/profile')
}

export function updateProfile(nickname) {
    return request.post('/user/update_profile', { nickname })
}

export function changePassword(old_password, new_password) {
    return request.post('/user/change_password', { old_password, new_password })
}
