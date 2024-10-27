import type { Config } from 'tailwindcss';

export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'class',
	theme: {
		extend: {}
	},
	fontFamily: {
		sans: [
			'ui-sans-serif',
			'system-ui',
			'-apple-system',
			'BlinkMacSystemFont',
			'Roboto',
			'Inter',
			'"Segoe UI"',
			'"Helvetica Neue"',
			'Arial',
			'"Noto Sans"',
			'sans-serif',
			'"Apple Color Emoji"',
			'"Segoe UI Emoji"',
			'"Segoe UI Symbol"',
			'"Noto Color Emoji"'
		],
		serif: ['ui-serif', 'Georgia', 'Cambria', '"Times New Roman"', 'Times', 'serif'],
		mono: [
			'ui-monospace',
			'SFMono-Regular',
			'Menlo',
			'Monaco',
			'Consolas',
			'"Liberation Mono"',
			'"Courier New"',
			'monospace'
		]
	},
	plugins: []
} as Config;
