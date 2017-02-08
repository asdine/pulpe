import merge from 'webpack-merge';
import ExtractTextPlugin from 'extract-text-webpack-plugin';
import baseConfig from './base';

export default merge(baseConfig, {
  devtool: 'source-map',

  entry: [
    'babel-polyfill',
    './app/index.jsx'
  ],

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
          use: [
            { loader: 'css-loader', query: { modules: true, sourceMaps: true } },
            { loader: 'sass-loader' },
            { loader: 'postcss-loader' },
          ]
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
  ]
});
