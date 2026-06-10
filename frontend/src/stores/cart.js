import { defineStore } from 'pinia'

const STORAGE_KEY = 'cart_items'

function loadFromStorage() {
    try {
        const raw = localStorage.getItem(STORAGE_KEY)
        return raw ? JSON.parse(raw) : []
    } catch (e) {
        return []
    }
}

function saveToStorage(items) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(items))
}

export const useCartStore = defineStore('cart', {
    state: () => ({
        items: loadFromStorage()
    }),
    getters: {
        totalCount: (state) => {
            return state.items.reduce((sum, item) => sum + item.quantity, 0)
        },
        totalPrice: (state) => {
            return state.items.reduce((sum, item) => sum + item.quantity * item.price, 0)
        }
    },
    actions: {
        addItem(product, quantity) {
            const existing = this.items.find((item) => item.product_id === product.id)
            if (existing) {
                existing.quantity += quantity || 1
            } else {
                this.items.push({
                    product_id: product.id,
                    name: product.name,
                    price: product.price,
                    image: product.image,
                    quantity: quantity || 1
                })
            }
            saveToStorage(this.items)
        },
        removeItem(productId) {
            this.items = this.items.filter((item) => item.product_id !== productId)
            saveToStorage(this.items)
        },
        updateQuantity(productId, quantity) {
            const item = this.items.find((i) => i.product_id === productId)
            if (item) {
                item.quantity = Math.max(1, quantity)
            }
            saveToStorage(this.items)
        },
        clearCart() {
            this.items = []
            saveToStorage([])
        }
    }
})
