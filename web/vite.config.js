import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import styleImport, {ElementPlusResolve} from 'vite-plugin-style-import'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import WindiCSS from 'vite-plugin-windicss'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': resolve(__dirname, '/src'),
      'views': resolve(__dirname, '/src/views'),
      'components': resolve(__dirname, '/src/components'),
    }
  },
  plugins: [
    vue(),
    WindiCSS(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [
        ElementPlusResolver(),
      ],
    }),
    styleImport({
      resolves: [ElementPlusResolve()],
    })
  ],
  server: {
    host: "127.0.0.1",
    port: 8081,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        ws: true,
        rewrite: (path) => path.replace(/\/api/, '/api')
      }
    },
  }
})
