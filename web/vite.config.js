import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  build:{
    cssMinify: false,
    outDir: '../cmd/server/web',
  },
  plugins: [react()],
})
