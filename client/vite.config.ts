import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import svgLoader from 'vite-svg-loader';

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), svgLoader()],
  build: {
    emptyOutDir: true,
    outDir: '../static',
  },
});
