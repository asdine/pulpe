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
import { combineReducers } from 'redux';
import { combineEpics } from 'redux-observable';
import configureStore from '@/store';
import LoginReducer, { epics as LoginEpics } from './duck';

const epics = combineEpics(
  LoginEpics,
);

const reducers = combineReducers({
  ...LoginReducer,
});

export default configureStore(reducers, epics);
