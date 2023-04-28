import {defineConfig} from 'vite';
import {resolve} from "path";
import dts from 'vite-plugin-dts';
import tsc from '@mesh/tsc';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        dts({
            insertTypesEntry: true,
        }),
        tsc(),
    ],
    resolve: {
        alias: {
            '@': resolve(__dirname, './src'),
        },
    },
    build: {
        lib: {
            // Could also be a dictionary or array of multiple entry points
            entry: resolve(__dirname, 'src/index.ts'),
            name: 'mesh',
            fileName: 'mesh'
        },
        rollupOptions: {}
    },
})
