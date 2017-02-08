import { combineReducers } from 'redux';
import { successType, UPDATE_CARD_POSITION, DELETE_CARD, CREATE_CARD, FETCH_BOARD } from '../actions/types';

const byID = (state = {}, action) => {
  if (action.response) {
    return {
      ...state,
      ...action.response.entities.cards
    };
  }

  switch (action.type) {
    case UPDATE_CARD_POSITION: {
      return {
        ...state,
        [action.id]: {
          ...state[action.id],
          position: action.position
        }
      };
    }
    case successType(DELETE_CARD): {
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
        [action.response.result]: action.response.entities.boards[action.response.result].cards
      };
    }
    case successType(CREATE_CARD): {
      const card = action.response.entities.cards[action.response.result];
      return {
        ...state,
        [card.boardID]: [
          ...state[card.boardID],
          card.id
        ]
      };
    }
    case successType(DELETE_CARD):
      return {
        ...state,
        [action.boardID]: state[action.boardID].filter((id) => id !== action.id)
      };
    default:
      return state;
  }
};

const cards = combineReducers({
  byID,
  IDsByBoardID,
});

export default cards;

export const getCardByID = (state, id) => state.byID[id];
export const getCardsByBoardID = (state, boardID) => state.IDsByBoardID[boardID];
export const getCardsByListID = (state, boardID, listID) => {
  const cardsOfBoard = state.IDsByBoardID[boardID];
  if (!cardsOfBoard) {
    return [];
  }
  return cardsOfBoard
    .map(id => state.byID[id])
    .filter(c => c.listID === listID)
    .sort((a, b) => a.position - b.position);
};

export const getCards = (state) => state.allIDs.map(id => state.byID[id]);
