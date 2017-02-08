import { createStore, applyMiddleware, compose } from 'redux';
import { createEpicMiddleware } from 'redux-observable';
import createLogger from 'redux-logger';
import rootReducer from './reducers';
import rootEpic from './epics';

const epicMiddleware = createEpicMiddleware(rootEpic);
const logger = createLogger();
export default function configureStore() {
  const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose; /* eslint no-underscore-dangle: 0 */ /* eslint max-len: 0 */
  const store = createStore(
    rootReducer,
    composeEnhancers(
      applyMiddleware(epicMiddleware, logger)
    )
  );

  return store;
}
