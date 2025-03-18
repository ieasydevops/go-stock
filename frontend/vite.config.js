import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      '@wailsjs': path.resolve(__dirname, 'src/wailsjs'),
    },
  },
  server: {
    port: 5173,
    host: true,
    watch: {
      usePolling: true,
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor': ['vue', 'vue-router', 'naive-ui'],
        },
      },
    },
  },
})
