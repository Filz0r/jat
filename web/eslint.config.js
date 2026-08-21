//  @ts-check

import { tanstackConfig } from '@tanstack/eslint-config';
import queryPlugin from '@tanstack/eslint-plugin-query';
import eslintPluginTailwindcss from 'eslint-plugin-tailwindcss';

export default [
	...tanstackConfig,
	eslintPluginTailwindcss.configs.recommended,
	{
		plugins: {
			'@tanstack/query': queryPlugin,
		},
		rules: {
			...queryPlugin.configs.recommended.rules,
			'import/no-cycle': 'off',
			'import/order': 'off',
			'sort-imports': 'off',
			'@typescript-eslint/array-type': 'off',
			'@typescript-eslint/require-await': 'off',
			'pnpm/json-enforce-catalog': 'off',
		},
	},
	{
		ignores: [
			'eslint.config.js',
			'prettier.config.js',
			'src/components/ui/**.ts',
			'src/components/ui/**.tsx',
			'src/components/ui/**.js',
		],
	},
	{
		settings: {
			// Define the tailwindcss settings with the MANDATORY `cssConfigPath`
			tailwindcss: {
				cssConfigPath: './src/styles.css',
			},
		},
		// Optional: Customize the rules to your needs
		rules: {
			'tailwindcss/classnames-order': 'warn',
			'tailwindcss/no-arbitrary-value': 'warn',
			'tailwindcss/no-custom-classname': ['warn', { whitelist: ['custom\\-*'] }],
			'tailwindcss/no-contradicting-classname': 'warn',
		},
	},
];
