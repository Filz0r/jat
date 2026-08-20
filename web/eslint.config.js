//  @ts-check

import { tanstackConfig } from '@tanstack/eslint-config';
import queryPlugin from '@tanstack/eslint-plugin-query';

export default [
	...tanstackConfig,
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
		ignores: ['eslint.config.js', 'prettier.config.js'],
	},
];
