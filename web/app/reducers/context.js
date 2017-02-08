import { successType, SET_ACTIVE_BOARD, FETCH_CARD } from '../actions/types';

const context = (state = {}, action) => {
  switch (action.type) {
    case SET_ACTIVE_BOARD: {
      return {
        ...state,
        activeBoard: action.activeBoard
      };
    }
    case successType(FETCH_CARD): {
      return {
        ...state,
        activeBoard: action.response.entities.cards[action.response.result].boardID
      };
    }
    default:
      return state;
  }
};

export default context;

export const getActiveBoardID = (state) => state.activeBoard;
