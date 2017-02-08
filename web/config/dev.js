import webpack from 'webpack';
import merge from 'webpack-merge';
import baseConfig from './base';

export default merge(baseConfig, {
  devtool: 'cheap-module-source-map',

  output: {
    filename: '[name].js',
  },

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
      },
      {
        test: /\.s?css$/,
        loaders: [
          'style-loader',
          'css-loader?!postcss-loader!sass-loader'
        ]
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
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin()
  ],

  devServer: {
    hot: true,
    contentBase: './app/',
    historyApiFallback: true,
    proxy: {
      '/v1': 'http://localhost:4000'
    }
  }
});
