/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{svelte,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                'gray': {
                    '850': '#1f2937',
                    '900': '#111827',
                    '950': '#0d131f'
                },
            }
        },
    },
    plugins: [],
}