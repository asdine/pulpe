import React from 'react';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { Router, Route, IndexRoute, browserHistory } from 'react-router';
import configureStore from '@/store';
import Home from '@/Home';
import Index from '@/Home/Routes/Index';
import Board from '@/Home/Routes/Board';
import Card from '@/Home/Routes/Card';

import '@/app.global.scss';

const store = configureStore();

render(
  <Provider store={store}>
    <Router history={browserHistory}>
      <Route path="/" component={Home}>
        <IndexRoute component={Index} />
        <Route path="/:board" component={Board}>
          <Route path="/:board/:list/:card" component={Card} />
        </Route>
      </Route>
    </Router>
  </Provider>,
  document.getElementById('root')
);

if (module.hot) {
  module.hot.accept('./store', () => {
    const nextRootReducer = require('./store'); // eslint-disable-line global-require
    store.replaceReducer(nextRootReducer);
  });
}
