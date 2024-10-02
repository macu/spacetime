const webpack = require('webpack');
const VueLoader = require('vue-loader');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = {
	// mode: 'production',
	mode: 'development',
	stats: 'errors-warnings',
	entry: [
		'./src/app.js',
	],
	output: {
		filename: 'compiled.js',
		path: __dirname + '/js',
	},
	optimization: {
		minimize: true,
	},
	performance: {
		hints: 'warning',
		maxEntrypointSize: 250000, // JS output 250 kB
		maxAssetSize: 250000, // CSS output 250 kB
	},
	externals: {
		'vue': 'Vue',
		'vuex': 'Vuex',
		'vue-router': 'VueRouter',
		'element-plus': 'ElementPlus',
		'element-en': 'ElementPlusLocaleEn',
		'moment': 'moment',
	},
	resolve: {
		alias: {
			'@': __dirname + '/src',
		},
	},
	module: {
		rules: [
			{
				test: /\.vue$/,
				loader: 'vue-loader',
			},
			{
				test: /\.m?js$/,
				resolve: {
					fullySpecified: false,
				},
				use: {
					loader: 'babel-loader',
					options: {
						presets: [
							['@babel/preset-env', { targets: '>1%' }],
						],
					},
				},
			},
			{
				test: /\.s?css$/,
				use: [
					MiniCssExtractPlugin.loader, // add support for `import 'file.scss';` in JS
					{
						loader: 'css-loader',
						options: {
							url: false, // whether to resolve urls; leave urls in the code as written
						},
					},
					{
						loader: 'sass-loader',
						options: {
							implementation: require('sass'),
							sassOptions: {
								includePaths: [
									//__dirname + '/bower_components/bootstrap-sass/assets/stylesheets',
								],
							},
						},
					},
				],
			},
		],
	},
	plugins: [
		new VueLoader.VueLoaderPlugin(),
		new MiniCssExtractPlugin({
			// Output destination for compiled CSS
			filename: '../css/compiled.css',
		}),
		new webpack.DefinePlugin({
			'process.env.NODE_ENV': JSON.stringify('development'),
			// https://github.com/vuejs/core/tree/main/packages/vue#bundler-build-feature-flags
			__VUE_OPTIONS_API__: true,
			__VUE_PROD_DEVTOOLS__: process.env.NODE_ENV === 'development',
		}),
	],
};
