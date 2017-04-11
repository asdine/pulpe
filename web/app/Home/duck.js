import { Map } from 'immutable';

export const DOMAIN = 'pulpe/home';

// types
const SET_ACTIVE_BOARD_ID = `${DOMAIN}/setActiveBoardID`;

// action creators
export const setActiveBoardID = (id) => ({
  type: SET_ACTIVE_BOARD_ID,
  activeBoardID: id
});

// reducer
const initialState = new Map({});

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case SET_ACTIVE_BOARD_ID: {
      return state.set('activeBoardID', action.activeBoardID);
    }
    // TODO: add this
    // case successType(FETCH_CARD): {
    //   return {
    //     ...state,
    //     activeBoard: action.response.entities.cards[action.response.result].boardID
    //   };
    // }
    default:
      return state;
  }
};

export default {
  [DOMAIN]: reducer
};

export const getActiveBoardID = (state) => state[DOMAIN].get('activeBoardID');
