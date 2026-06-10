import axios from 'axios'

const request = axios.create({
    baseURL: window.location.origin + '/api/v1',
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
    if (data.code === 0) {
        return data
    }
    if (data.code === 1004) {
        localStorage.removeItem('token')
        localStorage.removeItem('is_admin')
        window.location.href = '/user/login'
    }
    alert(data.msg || '请求失败')
    return Promise.reject(data)
}, (error) => {
    if (error.response) {
        const data = error.response.data
        if (data && data.code === 1004) {
            localStorage.removeItem('token')
            localStorage.removeItem('is_admin')
            window.location.href = '/user/login'
        }
        alert(data && data.msg ? data.msg : '网络请求失败')
    } else {
        alert('网络请求失败')
    }
    return Promise.reject(error)
})

export default request
