import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
import client from '@/services/api/client';
import ajaxEpic, { successOf, requestOf } from '@/services/api/ajaxEpic';

export const DOMAIN = 'pulpe/login';

// types
export const LOGIN = 'pulpe/login';

// action creators
export const login = (payload) => ({
  type: requestOf(LOGIN),
  payload
});


// epics
const loginEpic = ajaxEpic(
  LOGIN,
  action => client.login(action.payload)
);

const redirectOnLoginSuccessEpic = action$ => action$.ofType(successOf(LOGIN))
  .do(() => window.location.replace('/'))
  .mergeMap(() => Observable.empty());


const reducer = (state = {}) => state;

export default {
  [DOMAIN]: reducer
};

export const epics = combineEpics(
  loginEpic,
  redirectOnLoginSuccessEpic
);

