import path from 'path';
import HtmlWebpackPlugin from 'html-webpack-plugin';

const BUILD_DIR = path.resolve(__dirname, '../build');

export default {
  output: {
    path: BUILD_DIR,
    filename: '[chunkhash].[name].js',
    sourceMapFilename: '[name].map'
  },

  resolve: {
    extensions: ['.js', '.jsx', '.json']
  },

  plugins: [
    new HtmlWebpackPlugin({
      title: 'Pulpe',
      template: 'app/index.html',
      chunksSortMode: 'dependency'
    })
  ]
};
