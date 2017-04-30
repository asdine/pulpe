import { combineEpics } from 'redux-observable';
import client from '@/services/api/client';
import ajaxEpic, { requestOf } from '@/services/api/ajaxEpic';

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

const reducer = (state = {}) => state;

export default {
  [DOMAIN]: reducer
};

export const epics = combineEpics(
  loginEpic,
);

