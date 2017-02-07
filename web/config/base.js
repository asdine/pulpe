import path from 'path';
import validate from 'webpack-validator';
import HtmlWebpackPlugin from 'html-webpack-plugin';

const BUILD_DIR = path.resolve(__dirname, '../build');

export default validate({
  output: {
    path: BUILD_DIR,
    filename: '[name].bundle.js',
    sourceMapFilename: '[name].map'
  },

  resolve: {
    extensions: ['.js', '.jsx', '.json']
  },

  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loaders: ['babel-loader'],
      }
    ]
  },

  plugins: [
    new HtmlWebpackPlugin({
      title: 'Pulpe',
      template: 'app/index.html',
      chunksSortMode: 'dependency'
    }),
  ]
});
