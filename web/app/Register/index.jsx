import React from 'react';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import '@/app.global.scss';
import store from './store';
import Register from './Register';

render(
  <Provider store={store}>
    <Register />
  </Provider>,
  document.getElementById('root')
);

if (module.hot) {
  module.hot.accept('./store', () => {
    const nextRootReducer = require('./store'); // eslint-disable-line global-require
    store.replaceReducer(nextRootReducer);
  });
}
