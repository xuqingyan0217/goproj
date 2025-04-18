import { createRouter, createWebHashHistory } from 'vue-router'

import { hasToken } from '../utils/token'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue')
    },
    {
      path: '/sign-in',
      name: 'sign-in',
      component: () => import('../views/SignIn.vue')
    },
    {
      path: '/sign-up',
      name: 'sign-up',
      component: () => import('../views/SignUp.vue')
    },
    {
      path: '/product',
      name: 'product',
      component: () => import('../views/Product.vue')
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/About.vue')
    },
    {
      path: '/category/:category',
      name: 'category',
      component: () => import('../views/Category.vue')
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('../views/Search.vue')
    },
    {
      path: '/cart',
      name: 'cart',
      component: () => import('../views/Cart.vue')
    },
    {
      path: '/order',
      name: 'order',
      component: () => import('../views/OrderView.vue')
    },
    {
      path: '/ai',
      name: 'ai',
      component: () => import('../views/AIOrder.vue')
    },
    {
      path: '/checkout',
      name: 'checkout',
      component: () => import('../views/Checkout.vue')
    },
    {
      path: '/checkout/waiting',
      name: 'checkout-waiting',
      component: () => import('../views/Waiting.vue')
    },
    {
      path: '/checkout/result',
      name: 'checkout-result',
      component: () => import('../views/Result.vue')
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/Admin.vue'),
      children: [
        {
          path: 'product/create',
          name: 'create-product',
          component: () => import('../views/admin/CreateProduct.vue')
        },
        {
          path: 'product/delete',
          name: 'delete-product',
          component: () => import('../views/admin/DeleteProduct.vue')
        },
        {
          path: 'product/update',
          name: 'update-product',
          component: () => import('../views/admin/UpdateProduct.vue')
        }
      ]
    }
  ]
})

// 全局导航守卫
router.beforeEach((to, _from, next) => {
  // 检查是否是需要登录的页面
  if (to.path.startsWith('/admin') || to.path === '/cart' || to.path === '/order' || to.path === '/checkout') {
    if (!hasToken()) {
      // 未登录时重定向到登录页面，并记录原始目标路由
      next({ path: '/sign-in', query: { next: to.fullPath } })
      return
    }
  }
  next()
})

export default router