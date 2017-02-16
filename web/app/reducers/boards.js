import { combineReducers } from 'redux';
import { successType, CREATE_BOARD, DELETE_BOARD, FETCH_BOARDS } from '../actions/types';

const byID = (state = {}, action) => {
  if (action.response) {
    return {
      ...state,
      ...action.response.entities.boards
    };
  }

  switch (action.type) {
    case successType(DELETE_BOARD): {
      const { [action.id]: _, ...newState } = state; /* eslint no-unused-vars: 0 */
      return newState;
    }
    default:
      return state;
  }
};

const allIDs = (state = [], action) => {
  switch (action.type) {
    case successType(CREATE_BOARD): {
      return [
        ...state,
        action.response.result
      ];
    }
    case successType(FETCH_BOARDS): {
      return [...new Set([...state, ...action.response.result])];
    }
    case successType(DELETE_BOARD):
      return state.filter((id) => id !== action.id);
    default:
      return state;
  }
};

const boards = combineReducers({
  byID,
  allIDs,
});

export default boards;

export const getBoardByID = (state, id) => state.byID[id];
export const getBoardBySlug = (state, slug) =>
  Object.keys(state.byID).map(k => state.byID[k]).find(board => board.slug === slug);
export const getBoards = (state) => state.allIDs.map(id => state.byID[id]);
export const getFirstBoardSlug = (state) =>
  state.allIDs.length > 0 ? state.byID[state.allIDs[0]].slug : undefined;
