import webpack from 'webpack';
import merge from 'webpack-merge';
import path from 'path';
import ExtractTextPlugin from 'extract-text-webpack-plugin';
import ChunkManifestPlugin from 'chunk-manifest-webpack-plugin';
import WebpackChunkHash from 'webpack-chunk-hash';
import HtmlWebpackPlugin from 'html-webpack-plugin';
import baseConfig from './base';
import packageConfig from '../package.json';

export default merge(baseConfig, {
  devtool: 'source-map',

  output: {
    publicPath: '/assets/',
    path: path.resolve(__dirname, '../../dist')
  },

  entry: {
    home: './app/Home/index.jsx',
    register: './app/Register/index.jsx',
    login: './app/Login/index.jsx',
    vendor: Object.keys(packageConfig.dependencies)
        .filter(dep => packageConfig.excludedFromBuild.findIndex(exl => exl === dep) === -1)
  },

  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loaders: ['babel-loader'],
      },
      {
        test: /\.s?css$/,
        use: ExtractTextPlugin.extract({
          fallback: 'style-loader',
          use: 'css-loader?!postcss-loader!sass-loader'
        })
      },
      {
        test: /\.(jpg|png|gif|svg)$/,
        use: 'file-loader'
      },
      {
        test: /\.(woff|woff2|eot|ttf)$/,
        use: 'url-loader?limit=100000'
      }
    ]
  },

  plugins: [
    new ExtractTextPlugin('style.css'),
    new webpack.optimize.CommonsChunkPlugin({
      name: ['vendor', 'manifest'], // vendor libs + extracted manifest
      minChunks: Infinity,
    }),
    new webpack.HashedModuleIdsPlugin(),
    new WebpackChunkHash(),
    new ChunkManifestPlugin({
      filename: 'chunk-manifest.json',
      manifestVariable: 'webpackManifest'
    }),
    new HtmlWebpackPlugin({
      filename: 'board.tmpl.html',
      template: 'app/templates/board.tmpl.html',
      chunks: ['manifest', 'home', 'vendor'],
    }),
    new HtmlWebpackPlugin({
      filename: 'register.tmpl.html',
      template: 'app/templates/register.tmpl.html',
      chunks: ['manifest', 'register', 'vendor'],
    }),
    new HtmlWebpackPlugin({
      filename: 'login.tmpl.html',
      template: 'app/templates/login.tmpl.html',
      chunks: ['manifest', 'login', 'vendor'],
    }),
  ]
});
