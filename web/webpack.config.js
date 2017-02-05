import path from 'path';
import validate from 'webpack-validator';

const BUILD_DIR = path.resolve(__dirname, 'build');

export default validate({
  context: path.join(__dirname, 'app'),

  entry: [
    './index.js',
  ],

  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loaders: [
          'babel-loader',
        ],
      },
      {
        test: /\.json$/,
        loader: 'json-loader'
      }
    ]
  },

  output: {
    path: BUILD_DIR,
    filename: 'app.bundle.js',

    // https://github.com/webpack/webpack/issues/1114
    libraryTarget: 'commonjs2'
  },

  resolve: {
    extensions: ['.js', '.jsx', '.json']
  },
});
