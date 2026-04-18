import { defineConfig } from "vite";

export default defineConfig({
    envPrefix: "PB",
    base: "./",
    build: {
        chunkSizeWarningLimit: 1000,
        reportCompressedSize: false,
    },
    resolve: {
        alias: {
            "@": __dirname + "/src",
        },
    },
});
