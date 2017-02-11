import webpack from 'webpack';
import merge from 'webpack-merge';
import ExtractTextPlugin from 'extract-text-webpack-plugin';
import ChunkManifestPlugin from 'chunk-manifest-webpack-plugin';
import WebpackChunkHash from 'webpack-chunk-hash';
import baseConfig from './base';
import packageConfig from '../package.json';

export default merge(baseConfig, {
  devtool: 'source-map',

  output: {
    publicPath: '/assets/'
  },

  entry: {
    main: ['babel-polyfill', './app/index.jsx'],
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
        test: /\.(jpg|png|gif)$/,
        use: 'file-loader'
      },
      {
        test: /\.(woff|woff2|eot|ttf|svg)$/,
        use: 'url-loader?limit=100000'
      }
    ]
  },

  plugins: [
    new ExtractTextPlugin('style.css'),
    // new webpack.optimize.CommonsChunkPlugin({
    //   name: ['vendor', 'manifest'],
    //   minChunks: Infinity,
    // }),
    // new ChunkManifestPlugin({
    //   filename: 'chunk-manifest.json',
    //   manifestVariable: 'webpackManifest'
    // })
    new webpack.optimize.CommonsChunkPlugin({
      name: ['vendor', 'manifest'], // vendor libs + extracted manifest
      minChunks: Infinity,
    }),
    new webpack.HashedModuleIdsPlugin(),
    new WebpackChunkHash(),
    new ChunkManifestPlugin({
      filename: 'chunk-manifest.json',
      manifestVariable: 'webpackManifest'
    })
  ]
});
