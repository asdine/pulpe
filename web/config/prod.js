import validate from 'webpack-validator';
import merge from 'webpack-merge';
import baseConfig from './base';

export default validate(merge(baseConfig, {
  devtool: 'source-map',

  entry: [
    'babel-polyfill',
    './app/index.jsx'
  ],
}));
