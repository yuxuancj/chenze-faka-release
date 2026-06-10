import request from '../utils/request'

export function productList(page, size, keyword) {
    const params = { page: page || 1, size: size || 20, keyword: keyword || '' }
    return request.get('/api/v1/products', { params })
}

export function productDetail(id) {
    return request.get('/api/v1/products/' + id)
}

export function categoryList() {
    return request.get('/api/v1/categories')
}
