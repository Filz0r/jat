//  @ts-check

/** @type {import('prettier').Config} */
const config = {
	semi: true,
	trailingComma: 'all',
	singleQuote: true,
	printWidth: 100,
	tabWidth: 4,
	useTabs: true,
	plugins: ['prettier-plugin-tailwindcss'],
};

export default config;
