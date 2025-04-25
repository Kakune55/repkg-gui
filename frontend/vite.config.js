import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue({
    template: {
      compilerOptions: {
        //允许s-开头的自定义组件
        isCustomElement: (tag) => tag.startsWith('s-')
      }
    }
  })]
})
