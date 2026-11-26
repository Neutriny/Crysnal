import { fileURLToPath, URL } from 'node:url'

import { defineConfig, lazyPlugins } from 'vite-plus'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  lint: {
    jsPlugins: [{ name: 'vite-plus', specifier: 'vite-plus/oxlint-plugin' }],
    rules: { 'vite-plus/prefer-vite-plus-imports': 'error' },
    options: { typeAware: true, typeCheck: true },
  },
  fmt: {
    semi: false,
    singleQuote: true,
  },
  plugins: lazyPlugins(() => [UnoCSS(), vue(), vueDevTools()]),
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
