import { combineEpics } from 'redux-observable';
import client from '@/services/api/client';
import ajaxEpic, { requestOf } from '@/services/api/ajaxEpic';

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

const reducer = (state = {}) => state;

export default {
  [DOMAIN]: reducer
};

export const epics = combineEpics(
  registerEpic,
);

