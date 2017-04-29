import path from 'path';

const BUILD_DIR = path.resolve(__dirname, '../build');
const APP_DIR = path.resolve(__dirname, '../app');

export default {
  output: {
    path: BUILD_DIR,
    filename: '[chunkhash].[name].js',
    sourceMapFilename: '[name].map'
  },

  resolve: {
    alias: {
      '@': APP_DIR,
    },
    modules: ['node_modules'],
    extensions: ['.js', '.jsx', '.json']
  },
};
