import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import ElementPlus from 'unplugin-element-plus/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
    ElementPlus({
      useSource: true,
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    host: 'localhost',
    port: 80,
    open: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        configure: (proxy, _options) => {
          proxy.on('error', (err, _req, _res) => {
            console.log('proxy error', err);
          });
          proxy.on('proxyReq', (proxyReq, req, _res) => {
            console.log('Sending Request to the Target:', req.method, req.url);
          });
          proxy.on('proxyRes', (proxyRes, req, _res) => {
            console.log('Received Response from the Target:', proxyRes.statusCode, req.url);
          });
        },
      }
     }
  },
  build: {
    // 生产构建优化
    target: 'es2015',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
      },
    } as any,
    // 代码分割
    rollupOptions: {
      output: {
        manualChunks: {
          // 将 Vue 相关库打包到一个 chunk
          vue: ['vue', 'vue-router', 'pinia'],
          // 将 Element Plus 打包到一个 chunk
          'element-plus': ['element-plus'],
          // 将编辑器相关库打包到一个 chunk
          editor: ['@toast-ui/editor', '@toast-ui/editor-plugin-code-syntax-highlight'],
          // 将工具库打包到一个 chunk
          utils: ['axios', 'prismjs'],
        },
      },
    },
    // 启用 gzip 压缩大小报告
    reportCompressedSize: true,
    // 设置 chunk 大小警告限制
    chunkSizeWarningLimit: 1000,
  },
  // CSS 优化
  css: {
    preprocessorOptions: {
      scss: {
        // 移除 additionalData 以避免与 @use 规则冲突
      },
    },
  },
})
