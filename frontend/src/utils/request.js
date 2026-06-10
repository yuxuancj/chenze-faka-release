import axios from 'axios'

const request = axios.create({
    baseURL: '',
    timeout: 15000
})

request.interceptors.request.use((config) => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = 'Bearer ' + token
    }
    return config
}, (error) => {
    return Promise.reject(error)
})

request.interceptors.response.use((response) => {
    const data = response.data
    if (data && typeof data.code !== 'undefined') {
        if (data.code === 0) {
            return data.data !== undefined ? data.data : true
        }
        if (data.code === 1004) {
            localStorage.removeItem('token')
            localStorage.removeItem('is_admin')
            window.location.href = '/user/login'
            return Promise.reject(new Error(data.msg || '登录状态已过期'))
        }
        alert(data.msg || '请求失败')
        return Promise.reject(new Error(data.msg || '请求失败'))
    }
    return data
}, (error) => {
    if (error.response) {
        const data = error.response.data
        // Handle 401 unauthorized
        if (error.response.status === 401 || (data && data.code === 1004)) {
            localStorage.removeItem('token')
            localStorage.removeItem('is_admin')
            window.location.href = '/user/login'
            return Promise.reject(new Error(data && data.msg ? data.msg : '登录状态已过期'))
        }
        alert(data && data.msg ? data.msg : '网络请求失败')
    } else {
        alert('网络请求失败')
    }
    return Promise.reject(error)
})

export default request
