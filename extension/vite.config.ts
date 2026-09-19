import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
    plugins: [react()],
    build: {
        outDir: 'dist',
        rollupOptions: {
            input: {
                popup: resolve(__dirname, 'popup.html'),
                'content-leetcode-cn': resolve(__dirname, 'src/content/leetcode-cn.ts'),
            },
            output: {
                entryFileNames: (chunkInfo) => {
                    if (chunkInfo.name === 'content-leetcode-cn') {
                        return 'content/[name].js';
                    }
                    return 'assets/[name]-[hash].js';
                },
            },
        },
    },
});