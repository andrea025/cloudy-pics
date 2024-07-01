import {fileURLToPath, URL} from 'node:url'

import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig(({command, mode, ssrBuild}) => {
	const ret = {
		plugins: [vue()],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url))
			}
		},
	};
	ret.define = {
		"__API_URL__": JSON.stringify("http://18.232.92.85:3000"),
		"__BASE_URL__": JSON.stringify("https://cloudy-pics.s3.amazonaws.com/"),
	};
	return ret;
})
