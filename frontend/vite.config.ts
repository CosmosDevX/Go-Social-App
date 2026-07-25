import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      // optional: if you want to proxy API during dev
      // '/api': 'http://localhost:8080'
    }
  }
})
