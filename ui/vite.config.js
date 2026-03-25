import { defineConfig, loadEnv } from 'vite';
import { svelte, vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { fileURLToPath } from 'node:url';
import { dirname } from 'node:path';

const __dirname = dirname(fileURLToPath(import.meta.url));

// see https://vitejs.dev/config
export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, __dirname, '');
    const productName = env.PB_PRODUCT_NAME || 'PocketBase';

    return {
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
                preprocess: [vitePreprocess()],
                onwarn: (warning, handler) => {
                    if (warning.code.startsWith('a11y-')) {
                        return; // silence a11y warnings
                    }
                    handler(warning);
                },
            }),
            {
                name: 'inject-product-name-html',
                transformIndexHtml(html) {
                    return html.replace(/__PB_PRODUCT_NAME__/g, productName);
                },
            },
        ],
        resolve: {
            alias: {
                '@': __dirname + '/src',
            }
        },
    };
});
