import request from '../utils/request'

export function productList(page, size, category_id, keyword) {
    return request.post('/product/list', { page, size, category_id, keyword })
}

export function productDetail(id) {
    return request.post('/product/detail', { id })
}

export function categoryList() {
    return request.post('/category/list')
}
