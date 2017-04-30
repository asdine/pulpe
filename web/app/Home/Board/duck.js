import { schema } from 'normalizr';
import { combineEpics } from 'redux-observable';
import client from '@/services/api/client';
import ajaxEpic, { successOf, requestOf } from '@/services/api/ajaxEpic';

export const DOMAIN = 'pulpe/home/board';

// types
export const FETCH = `${DOMAIN}/fetch`;
export const UPDATE = `${DOMAIN}/update`;
export const DELETE = `${DOMAIN}/delete`;
export const PATCH = `${DOMAIN}/patch`;

// schemas
const listSchema = new schema.Entity('lists');
const cardSchema = new schema.Entity('cards');

const boardSchema = new schema.Entity('boards', {
  lists: [listSchema],
  cards: [cardSchema]
});

// action creators
export const fetchBoard = (owner, slug) => ({
  type: requestOf(FETCH),
  owner,
  slug
});

export const updateBoard = ({ id, ...patch }) => ({
  type: requestOf(UPDATE),
  id,
  patch
});

export const deleteBoard = (id) => ({
  type: requestOf(DELETE),
  id
});

export const patchBoard = ({ id, ...patch }) => ({
  type: PATCH,
  id,
  patch
});

// epics
const fetchBoardEpic = ajaxEpic(
  FETCH,
  action => client.getBoard(action.owner, action.slug),
  boardSchema
);

const updateBoardEpic = ajaxEpic(
  UPDATE,
  action => client.updateBoard(action),
  boardSchema
);

const deleteBoardEpic = ajaxEpic(
  DELETE,
  action => client.deleteBoard(action.id).map(() => ({
    ...action,
    type: successOf(DELETE),
  }))
);

export const epics = combineEpics(
  fetchBoardEpic,
  updateBoardEpic,
  deleteBoardEpic,
);

// reducer
const reducer = (state = {}, action) => {
  switch (action.type) {
    case successOf(FETCH) || successOf(UPDATE): {
      return action.response.entities.boards[action.response.result];
    }
    case successOf(DELETE): {
      return {};
    }
    case PATCH: {
      return {
        ...state,
        ...action.patch,
      };
    }
    default:
      return state;
  }
};

export default {
  [DOMAIN]: reducer
};

export const getBoardSelector = (state) => state[DOMAIN];
