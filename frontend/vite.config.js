import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
    plugins: [react()],
    define: {
        // 将环境变量嵌入到应用程序中
        'process.env.INFURA_API_KEY': JSON.stringify(process.env.INFURA_API_KEY),
    },
});