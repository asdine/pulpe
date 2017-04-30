import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
import client from '@/services/api/client';
import ajaxEpic, { successOf, requestOf, failureOf } from '@/services/api/ajaxEpic';

export const DOMAIN = 'pulpe/register';

// types
export const REGISTER = 'pulpe/register';

// action creators
export const register = (payload) => ({
  type: requestOf(REGISTER),
  payload
});


// epics
const registerEpic = ajaxEpic(
  REGISTER,
  action => client.register(action.payload)
);

const redirectOnRegisterSuccessEpic = action$ => action$.ofType(successOf(REGISTER))
  .do(() => window.location.replace('/'))
  .mergeMap(() => Observable.empty());

const reducer = (state = {}, action) => {
  switch (action.type) {
    case failureOf(REGISTER): {
      return {
        ...action.payload,
      };
    }
    default:
      return state;
  }
};

export default {
  [DOMAIN]: reducer
};

export const epics = combineEpics(
  registerEpic,
  redirectOnRegisterSuccessEpic
);

export const getErrors = (state) => state[DOMAIN].fields || {};
