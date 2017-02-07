import webpack from 'webpack';
import validate from 'webpack-validator';
import merge from 'webpack-merge';
import baseConfig from './base';

export default validate(merge(baseConfig, {
  devtool: 'cheap-module-source-map',

  entry: [
    'react-hot-loader',
    'babel-polyfill',
    './app/index.jsx'
  ],

  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loaders: [
          'babel-loader',
          'eslint-loader'
        ],
      }
    ]
  },

  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin(),
  ],

  devServer: {
    hot: true,
    contentBase: './app/',
    historyApiFallback: true,
  }
}));
