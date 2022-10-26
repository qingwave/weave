import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import prismjs from "vite-plugin-prismjs";
import fs from 'fs'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      'views': resolve(__dirname, 'src/views'),
      'components': resolve(__dirname, 'src/components'),
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@use "@/styles/element.scss" as *;`,
      },
    },
  },
  plugins: [
    vue(),
    prismjs({
      languages: ["json", "js", "go", "bash", "yaml", "markup"],
      plugins: ["line-numbers"],
      theme: "solarizedlight",
      css: true,
    }),
    AutoImport({
      resolvers: [ElementPlusResolver({ importStyle: "sass" })],
    }),
    Components({
      resolvers: [
        ElementPlusResolver({ importStyle: "sass" }),
      ],
    })
  ],
  server: {
    // if your frontend not in the localhost, please uncomment the https config meanwhile
    host: "127.0.0.1",
    port: 8081,
    // https: {
    //   ca: fs.readFileSync('../certs/root.crt'),
    //   key: fs.readFileSync('../certs/frontend.key'),
    //   cert: fs.readFileSync('../certs/frontend.crt')
    // },
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
