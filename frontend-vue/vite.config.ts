import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  base: '/',
  server: {
    proxy: {
      '/auth': 'http://frontend-svc:8080',
      '/api': 'http://frontend-svc:8080',
      '/product': 'http://frontend-svc:8080',
      '/category': 'http://frontend-svc:8080',
      '/cart': 'http://frontend-svc:8080',
      '/search': 'http://frontend-svc:8080',
      '/order': 'http://frontend-svc:8080',
      '/checkout': 'http://frontend-svc:8080',
      '/ai': 'http://frontend-svc:8080'
    },
    fs: {
      strict: true
    }
  }
})
