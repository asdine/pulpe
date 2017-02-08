import { createStore, applyMiddleware, compose } from 'redux';
import { createEpicMiddleware } from 'redux-observable';
import rootReducer from './reducers';
import rootEpic from './epics';

const middlewares = [createEpicMiddleware(rootEpic)];

if (process.env.NODE_ENV === 'development') {
  middlewares.push(require('redux-logger')()); // eslint-disable-line
}

export default function configureStore() {
  const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose; /* eslint no-underscore-dangle: 0 */ /* eslint max-len: 0 */
  const store = createStore(
    rootReducer,
    composeEnhancers(
      applyMiddleware(...middlewares)
    )
  );

  return store;
}
