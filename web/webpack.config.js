import path from 'path';
import webpack from 'webpack';
import validate from 'webpack-validator';

const BUILD_DIR = path.resolve(__dirname, 'build');

export default validate({
  entry: [
    'webpack/hot/only-dev-server', // "only" prevents reload on syntax errors
    'babel-polyfill',
    './app/index.jsx',
    './app/index.html'
  ],

  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loaders: [
          'react-hot-loader',
          'babel-loader',
        ],
      },
      {
        test: /\.html$/,
        loader: "file-loader?name=[name].[ext]",
      },
      {
        test: /\.json$/,
        loader: 'json-loader'
      }
    ]
  },

  output: {
    path: BUILD_DIR,
    filename: 'app.bundle.js'
  },

  resolve: {
    extensions: ['.js', '.jsx', '.json']
  },

  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin(),
  ],

  devServer: {
    hot: true,
    contentBase: './app/'
  }
});
