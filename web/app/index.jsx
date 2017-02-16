import React from 'react';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { Router, Route, browserHistory, IndexRoute } from 'react-router';
import configureStore from './configureStore';
import FirstBoard from './containers/firstBoard';
import BoardActivate from './containers/boardActivate';
import CardDetail from './containers/cardDetail';
import { fetchBoards } from './actions';
import Home from './layouts/Home';

import './app.global.scss';

const store = configureStore();

store.dispatch(fetchBoards());

render(
  <Provider store={store}>
    <Router history={browserHistory}>
      <Route path="/" component={Home}>
        <IndexRoute component={FirstBoard} />
        <Route path="/:slug" component={BoardActivate} />
        <Route path="/:slug/:id" component={CardDetail} />
      </Route>
    </Router>
  </Provider>,
  document.getElementById('root')
);

if (module.hot) {
  module.hot.accept('./reducers', () => {
    const nextRootReducer = require('./reducers/index'); // eslint-disable-line global-require
    store.replaceReducer(nextRootReducer);
  });
}
