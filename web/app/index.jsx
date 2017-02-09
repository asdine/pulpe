import React from 'react';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { Router, Route, browserHistory, IndexRoute } from 'react-router';
import configureStore from './configureStore';
import FirstBoard from './containers/firstBoard';
import BoardActivate from './containers/boardActivate';
import CardDetail from './containers/cardDetail';
import CardCreate from './containers/cardCreate';
import { MainModal } from './containers/modal';
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
        <Route path="/b/:id" component={BoardActivate} />
        <Route path="/b/:id/:listID" component={MainModal}>
          <Route path="/b/:id/:listID/newcard" component={CardCreate} />
        </Route>
        <Route path="/c" component={MainModal}>
          <Route path="/c/:id" component={CardDetail} />
        </Route>
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
