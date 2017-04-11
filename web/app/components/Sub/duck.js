import { combineReducers } from 'redux';

const DOMAIN = 'pulpe/sub';

// types
export const ADD_SUB = `${DOMAIN}/add`;
export const POP_SUB = `${DOMAIN}/pop`;
export const CLOSE_SUB = `${DOMAIN}/close`;
export const CLOSE_ALL_SUBS = `${DOMAIN}/closeAll`;
export const REGISTER_SUB = `${DOMAIN}/register`;

// action creators
export const addSub = (name) => ({
  type: ADD_SUB,
  name,
});

export const popSub = () => ({
  type: POP_SUB,
});

export const closeSub = (name) => ({
  type: CLOSE_SUB,
  name,
});

export const closeAllSubs = () => ({
  type: CLOSE_ALL_SUBS,
});

export const registerSub = (name) => ({
  type: REGISTER_SUB,
  name,
});

// reducer
const opened = (state = [], action) => {
  switch (action.type) {
    case ADD_SUB:
      return [
        ...state,
        action.name
      ];
    case POP_SUB:
      if (state.length > 0) {
        return state.slice(0, -1);
      }
      return [];
    case CLOSE_ALL_SUBS:
      return [];
    default:
      return state;
  }
};

const subs = (state = {}, action) => {
  switch (action.type) {
    case REGISTER_SUB:
      return {
        ...state,
        [action.name]: false
      };
    case ADD_SUB:
      return {
        ...state,
        [action.name]: true
      };
    case CLOSE_SUB:
      return {
        ...state,
        [action.name]: false
      };
    case CLOSE_ALL_SUBS:
      return Object.keys(state).reduce((acc, k) => ({ ...acc, [k]: false }), {});
    default:
      return state;
  }
};

export default {
  [DOMAIN]: combineReducers({
    opened,
    subs,
  })
};

export const subsStillOpened = (state) => state[DOMAIN].opened.length > 0;
export const subIsOpened = (state, name) => !!state[DOMAIN].subs[name];
export const getLastOpened = (state) => state[DOMAIN].opened.length > 0 ? state[DOMAIN].opened[0] : null;
