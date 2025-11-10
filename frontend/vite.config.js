import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path'

export default defineConfig({
    plugins: [svelte()],
    resolve: {
        alias: {
            '@wailsjs/runtime': path.resolve(__dirname, './wailsjs/runtime'),
            '@wailsjs/go': path.resolve(__dirname, './wailsjs/go')
        }
    },
    build: {
        outDir: '../build',
        emptyOutDir: true
    },
    server: {
        port: 5173,
        strictPort: true,
        host: 'localhost'
    }
})