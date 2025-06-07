import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: './dist',
    emptyOutDir: true,
  },
})
