/** @type {import('tailwindcss').Config}*/
const config = {
  content: [
		'./src/**/*.{html,js,svelte,ts}',
		'./node_modules/stwui/**/*.{svelte,js,ts,html}'
	],
	plugins: [
		require('@tailwindcss/forms'),
		require('stwui/plugin')
	]

};
module.exports = config;
