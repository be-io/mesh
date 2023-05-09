import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';
import {resolve} from "path";

// https://vitejs.dev/config/
export default defineConfig({
    base: 'mesh',
    plugins: [react()],
    resolve: {
        alias: {
            '@': resolve(__dirname, './src'),
        },
    },
    build: {
        outDir: resolve(__dirname, './static'),
    },
    server: {
        proxy: {
            '/mesh/invoke': {
                secure: false,
                target: 'https://127.0.0.1',
            },
        }
    }
})
