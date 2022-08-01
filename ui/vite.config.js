import { defineConfig } from 'vite';
import { svelte }       from '@sveltejs/vite-plugin-svelte';

// see https://vitejs.dev/config
export default defineConfig({
    server: {
        port: 3000,
    },
    envPrefix: 'PB',
    base: './',
    build: {
        chunkSizeWarningLimit: 1000,
        reportCompressedSize: false,
    },
    plugins: [
        svelte({
            experimental: {
                useVitePreprocess: true,
            },
        }),
    ],
    resolve: {
        alias: {
            '@': __dirname + '/src',
        }
    },
    css: {
        preprocessorOptions: {
            scss: {
                includePaths: [__dirname + '/src/scss'],
                prependData: `@import "${__dirname}/src/scss/mixins"; @import "${__dirname}/src/scss/vars";`,
            },
        },
    },
})
