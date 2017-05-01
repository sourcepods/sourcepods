const CopyWebpackPlugin = require('copy-webpack-plugin');

const externals = {
	'jquery': 'jQuery',
	'uikit': 'UIkit',
	'vue': 'Vue',
	'vue-router': 'VueRouter',
	'vuex': 'Vuex',
};

const loaders = [
	{loader: 'vue-loader', test: /\.vue$/},
	{loader: 'babel-loader', test: /\.js/, exclude: /node_modules/},
	{loader: 'file-loader', test: /\.(png|jpg|gif|svg)$/, options: {name: '[name].[ext]'}},
];

module.exports = [
	{
		entry: './ui/main.js',
		output: {
			filename: './public/js/main.js'
		},
		externals,
		module: {
			loaders
		},
		plugins: [
			new CopyWebpackPlugin([
				{from: './node_modules/jquery/dist/jquery.min.js', to: './public/js'},
				{from: './node_modules/uikit/dist/css/uikit.css', to: './public/css'},
				{from: './node_modules/uikit/dist/js/uikit-icons.js', to: './public/js'},
				{from: './node_modules/uikit/dist/js/uikit.js', to: './public/js'},
				{from: './node_modules/vue-router/dist/vue-router.js', to: './public/js'},
				{from: './node_modules/vue/dist/vue.js', to: './public/js'},
				{from: './node_modules/vuex/dist/vuex.js', to: './public/js'},
			])
		]
	}
];






