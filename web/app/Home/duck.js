import { Map } from 'immutable';

export const DOMAIN = 'pulpe/home';

// types
const SET_ACTIVE_BOARD = `${DOMAIN}/setActiveBoard`;

// action creators
export const setActiveBoard = (owner, slug) => ({
  type: SET_ACTIVE_BOARD,
  activeBoard: {
    owner, slug
  },
});

// reducer
const initialState = new Map({});

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case SET_ACTIVE_BOARD: {
      return state.set('activeBoard', action.activeBoard);
    }
    default:
      return state;
  }
};

export default {
  [DOMAIN]: reducer
};

export const getActiveBoard = (state) => state[DOMAIN].get('activeBoard');
