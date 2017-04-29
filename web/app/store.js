import { createStore, applyMiddleware, compose } from 'redux';
import { createEpicMiddleware } from 'redux-observable';

export default function configureStore(reducers, epics) {
  const middlewares = [createEpicMiddleware(epics)];

  if (process.env.NODE_ENV === 'development') {
    // const dev = require('./dev'); // eslint-disable-line global-require
    // reducers.dev = dev.default;
    // module.exports.getDevStore = (state) => dev.getDevStore(state.dev);
    middlewares.push(require('redux-logger')()); // eslint-disable-line
  }

  const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose; /* eslint no-underscore-dangle: 0 */ /* eslint max-len: 0 */
  const store = createStore(
    reducers,
    composeEnhancers(
      applyMiddleware(...middlewares)
    )
  );

  return store;
}
