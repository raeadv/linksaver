import { fileURLToPath, URL } from 'node:url'
import { createRequire } from 'module'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import path from 'path'
import AutoImport from 'unplugin-auto-import/vite'

const _require = createRequire(import.meta.url)
const ionicons5Exports = Object.keys(_require('@vicons/ionicons5'))

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(__dirname, '../'), '')
  const baseUrl = env.VITE_BASE_URL

  const [host, port] = parseBaseUrl(baseUrl)
  const isDev = mode === 'development'

  return {
    envDir: path.resolve(__dirname, "../"),
    plugins: [
      vue(),
      isDev && vueDevTools(),
      AutoImport({
        imports: [
          'vue',
          'vue-router',
          { '@vicons/ionicons5': ionicons5Exports },
          {
            'axios': [
              // default imports
              ['default', 'axios'], // import { default as axios } from 'axios',
            ],
          }
        ],
        include: [
          /\.[tj]sx?$/, // .ts, .tsx, .js, .jsx
          /\.vue$/,
          /\.vue\?vue/, // .vue
          /\.vue\.[tj]sx?\?vue/, // .vue (vue-loader with experimentalInlineMatchResource enabled)
        ],
        vueTemplate: true,
        dts: './auto-imports.d.ts',
        viteOptimizeDeps: true,
      })
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
      cssMinify: true,
      reportCompressedSize: true,
      rollupOptions: {
        output: {
          manualChunks: {
            'lodash': ['lodash'],
            'vjs': ['vue', 'vue-router', 'pinia'],
            'n-ui': ['naive-ui'],
            'icons': ['@vicons/ionicons5'],
          }
        }
      }
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