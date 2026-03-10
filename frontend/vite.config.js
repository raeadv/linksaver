import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import path from 'path'


// https://vite.dev/config/
export default defineConfig(mode => {
  const env = loadEnv(mode, path.resolve(__dirname, '../'), '')
  const baseUrl = env.VITE_BASE_URL

  const [host, port] = parseBaseUrl(baseUrl)

  return {
    envDir: path.resolve(__dirname, "../"),
    plugins: [
      vue(),
      vueDevTools(),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      },
    },
    build: {
      emptyOutDir: true,
      outDir: path.resolve(__dirname, './../compiled/html'),
      minify: true,
    },
    server: {
      host: host,
      port: port
    }
  }
})


const parseBaseUrl = base_url => {

  const parts = String(base_url).split(':')
  const host = String(parts[1]).replace("//", "")

  return [host, parts[2]]
}