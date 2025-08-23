import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	// ssr: {
	// 	noExternal: [],
	// 	external: ['layerchart']
	// },
	// optimizeDeps: {
	// 	include: ['layerchart']
	// },
	server: {
		host: '0.0.0.0',
		port: 5173,
		watch: {
			usePolling: true,
			interval: 1000,
		},
		hmr: {
			port: 5173,
			host: 'localhost'
		},
		proxy: {
			'/api': {
				// target: process.env.PUBLIC_API_URL as string || 'http://localhost:4000',
				target: 'http://server:4000',
				changeOrigin: true,
				secure: false,
			}
		}
	}
});