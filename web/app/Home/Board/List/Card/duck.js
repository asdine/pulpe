import { schema } from 'normalizr';
import { combineReducers } from 'redux';
import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
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
export const DROP = `${DOMAIN}/drop`;

export const MODAL_CREATE_CARD = `${DOMAIN}/modalCreateCard`;
export const MODAL_CARD_DETAIL = `${DOMAIN}/modalCardDetail`;

// schemas
const cardSchema = new schema.Entity('cards');

// action creators
export const createCard = ({ boardID, listID, name, description }) => ({
  type: CREATE,
  boardID,
  listID,
  name,
  description,
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

export const dropCard = (card, index, canceled) => ({
  type: DROP,
  card,
  index,
  canceled
});

// epics
const beforeCreateCardEpic = (action$, store) => action$.ofType(CREATE)
  .map(action => {
    const { listID } = action;
    const cards = getCardsByListIDSelector(store.getState(), listID);
    let position = 1 << 16;

    if (cards.length !== 0) {
      position += cards[cards.length - 1].position;
    }

    return {
      ...action,
      position,
      type: requestOf(CREATE),
    };
  });

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

const onDropCardEpic = (action$, store) => action$.ofType(DROP)
  .mergeMap(action => {
    const { card: outDatedCard, index, canceled } = action;

    if (canceled) {
      return Observable.of(patchCard({
        id: outDatedCard.id,
        listID: outDatedCard.listID
      }));
    }

    const card = getCardSelector(store.getState(), outDatedCard.id);
    const cards = getCardsByListIDSelector(store.getState(), card.listID);

    const patch = {
      id: card.id
    };

    if (outDatedCard.listID !== card.listID) {
      patch.listID = card.listID;
    }

    if (cards[index].id === card.id && !patch.listID) {
      return Observable.empty();
    }

    if (cards.length === 1) {
      patch.position = 1 << 16;
    } else if (index === 0) {
      const { position: nextPosition } = cards[0];
      patch.position = nextPosition / 2;
    } else if (index < cards.length - 1) {
      if (card.position > cards[index].position) {
        const { position: prevPosition } = cards[index - 1];
        const { position: nextPosition } = cards[index];
        patch.position = prevPosition + ((nextPosition - prevPosition) / 2);
      } else {
        const { position: prevPosition } = cards[index];
        const { position: nextPosition } = cards[index + 1];
        patch.position = prevPosition + ((nextPosition - prevPosition) / 2);
      }
    } else {
      const { position: prevPosition } = cards[index];
      patch.position = prevPosition + (1 << 16);
    }

    if (patch.position === card.position && !patch.listID) {
      return Observable.empty();
    }

    return Observable.of(patchCard(patch), updateCard(patch));
  });

export const epics = combineEpics(
  beforeCreateCardEpic,
  createCardEpic,
  fetchCardEpic,
  closeCreateCardModalOnSuccessEpic,
  updateCardEpic,
  deleteCardEpic,
  onDropCardEpic,
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
      const { [action.response.id]: _, ...newState } = state; /* eslint no-unused-vars: 0 */
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
      return state.filter(id => id !== action.response.id);
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
    .sort((a, b) => a.position > b.position)
    .filter(c => c.listID === listID);
