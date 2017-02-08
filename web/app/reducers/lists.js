import { combineReducers } from 'redux';
import { successType, CREATE_LIST, DELETE_LIST, FETCH_BOARD } from '../actions/types';

const byID = (state = {}, action) => {
  if (action.response) {
    return {
      ...state,
      ...action.response.entities.lists
    };
  }

  switch (action.type) {
    case successType(DELETE_LIST): {
      const { [action.id]: _, ...newState } = state; /* eslint no-unused-vars: 0 */
      return newState;
    }
    default:
      return state;
  }
};

const IDsByBoardID = (state = {}, action) => {
  switch (action.type) {
    case successType(FETCH_BOARD): {
      return {
        ...state,
        [action.response.result]: action.response.entities.boards[action.response.result].lists
      };
    }
    case successType(CREATE_LIST): {
      const list = action.response.entities.lists[action.response.result];

      return {
        ...state,
        [list.boardID]: [
          ...state[list.boardID],
          action.response.result
        ]
      };
    }
    case successType(DELETE_LIST):
      return {
        ...state,
        [action.boardID]: state[action.boardID].filter((id) => id !== action.id)
      };

    default:
      return state;
  }
};

const lists = combineReducers({
  byID,
  IDsByBoardID,
});

export default lists;

export const getListByID = (state, id) => state.byID[id];
export const getListsByBoardID = (state, boardID) => state.IDsByBoardID[boardID] !== undefined ?
    state.IDsByBoardID[boardID].map(id => state.byID[id]) : [];
