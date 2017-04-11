import 'rxjs/add/operator/mergeMap';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/mapTo';
import 'rxjs/add/operator/ignoreElements';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/do';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/delay';
import 'rxjs/add/observable/of';
import 'rxjs/add/observable/empty';
import 'rxjs/add/observable/dom/ajax';
import { createStore, applyMiddleware, compose, combineReducers } from 'redux';
import { createEpicMiddleware, combineEpics } from 'redux-observable';
import { reducers as HomeReducers, epics as HomeEpics } from './Home/store';
import ModalReducer from './components/Modal/duck';
import SubReducer from './components/Sub/duck';

const reducers = {
  ...HomeReducers,
  ...ModalReducer,
  ...SubReducer,
};

const epics = combineEpics(
  HomeEpics,
);

const middlewares = [createEpicMiddleware(epics)];

if (process.env.NODE_ENV === 'development') {
  // const dev = require('./dev'); // eslint-disable-line global-require
  // reducers.dev = dev.default;
  // module.exports.getDevStore = (state) => dev.getDevStore(state.dev);
  middlewares.push(require('redux-logger')()); // eslint-disable-line
}

export default function configureStore() {
  const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose; /* eslint no-underscore-dangle: 0 */ /* eslint max-len: 0 */
  const store = createStore(
    combineReducers(reducers),
    composeEnhancers(
      applyMiddleware(...middlewares)
    )
  );

  return store;
}
