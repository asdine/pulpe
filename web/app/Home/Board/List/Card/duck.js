import { schema } from 'normalizr';
import { combineReducers } from 'redux';
import { combineEpics } from 'redux-observable';
import client from '@/services/api/client';
import ajaxEpic, { successOf, requestOf } from '@/services/api/ajaxEpic';
import { hideModal } from '@/components/Modal/duck';
import { FETCH as FETCH_BOARD, UPDATE as UPDATE_BOARD } from '@/Home/Board/duck';

export const DOMAIN = 'pulpe/home/board/list/card';

// types
export const CREATE = `${DOMAIN}/create`;
export const FETCH = `${DOMAIN}/fetch`;
export const UPDATE = `${DOMAIN}/update`;
export const DELETE = `${DOMAIN}/delete`;
export const PATCH = `${DOMAIN}/patch`;

export const MODAL_CREATE_CARD = `${DOMAIN}/modalCreateCard`;
export const MODAL_CARD_DETAIL = `${DOMAIN}/modalCardDetail`;

// schemas
const cardSchema = new schema.Entity('cards');

// action creators
export const createCard = ({ boardID, listID, name, description, position }) => ({
  type: requestOf(CREATE),
  boardID,
  listID,
  name,
  description,
  position
});

export const fetchCard = (id) => ({
  type: requestOf(FETCH),
  id
});

export const updateCard = ({ id, ...patch }) => ({
  type: requestOf(UPDATE),
  id,
  patch
});

export const deleteCard = (id) => ({
  type: requestOf(DELETE),
  id
});


export const patchCard = ({ id, ...patch }) => ({
  type: PATCH,
  id,
  patch
});

// epics
const createCardEpic = ajaxEpic(
  CREATE,
  action => client.createCard(action),
  cardSchema
);

const closeCreateCardModalOnSuccessEpic = action$ => action$.ofType(successOf(CREATE))
  .map(hideModal);

const fetchCardEpic = ajaxEpic(
  FETCH,
  action => client.getCard(action.id),
  cardSchema
);

const updateCardEpic = ajaxEpic(
  UPDATE,
  action => client.updateCard(action),
  cardSchema
);

const deleteCardEpic = ajaxEpic(
  DELETE,
  action => client.deleteCard(action.id).map(() => ({
    ...action,
    type: successOf(DELETE),
  }))
);

export const epics = combineEpics(
  createCardEpic,
  fetchCardEpic,
  closeCreateCardModalOnSuccessEpic,
  updateCardEpic,
  deleteCardEpic,
);

// reducer
const byID = (state = {}, action) => {
  switch (action.type) {
    case successOf(FETCH_BOARD): {
      return action.response.entities.cards || {};
    }
    case successOf(CREATE):
    case successOf(FETCH):
    case successOf(UPDATE): {
      return {
        ...state,
        [action.response.result]: action.response.entities.cards[action.response.result]
      };
    }
    case PATCH: {
      return {
        ...state,
        [action.id]: {
          ...state[action.id],
          ...action.patch
        },
      };
    }
    case successOf(DELETE): {
      const { [action.id]: _, ...newState } = state; /* eslint no-unused-vars: 0 */
      return newState;
    }
    default:
      return state;
  }
};

const ids = (state = [], action) => {
  switch (action.type) {
    case successOf(FETCH_BOARD): {
      return action.response.entities.boards[action.response.result].cards || [];
    }
    case successOf(CREATE): {
      return [
        ...state,
        action.response.result
      ];
    }
    case successOf(DELETE): {
      return state.filter(id => id !== action.id);
    }
    default:
      return state;
  }
};


export default {
  [DOMAIN]: combineReducers({
    byID,
    ids,
  })
};

export const getCardSelector = (state, id) => state[DOMAIN].byID[id];
export const getCardBySlugSelector = (state, slug) =>
  state[DOMAIN].ids
    .map(id => state[DOMAIN].byID[id])
    .find(c => c.slug === slug);

export const getCardsByListIDSelector = (state, listID) =>
  state[DOMAIN].ids
    .map(id => state[DOMAIN].byID[id])
    .filter(c => c.listID === listID);
