import { schema } from 'normalizr';
import { combineReducers } from 'redux';
import { combineEpics } from 'redux-observable';
import client from '@/services/api/client';
import { hideModal } from '@/components/Modal/duck';
import ajaxEpic, { successOf, requestOf } from '@/services/api/ajaxEpic';
import { FETCH as FETCH_BOARD, UPDATE as UPDATE_BOARD } from '@/Home/Board/duck';

export const DOMAIN = 'pulpe/home/board/list';

// types
export const CREATE = `${DOMAIN}/create`;
export const UPDATE = `${DOMAIN}/update`;
export const DELETE = `${DOMAIN}/delete`;
export const PATCH = `${DOMAIN}/patch`;

export const MODAL_CREATE_LIST = `${DOMAIN}/modalCreateList`;
export const MODAL_DELETE_LIST = `${DOMAIN}/modalDeleteList`;

// schemas
const listSchema = new schema.Entity('lists');

// action creators
export const createList = (boardID, name) => ({
  type: requestOf(CREATE),
  boardID,
  name
});

export const updateList = ({ id, ...rest }) => ({
  type: requestOf(UPDATE),
  id,
  toUpdate: rest
});

export const deleteList = (id) => ({
  type: requestOf(DELETE),
  id
});

export const patchList = ({ id, ...patch }) => ({
  type: PATCH,
  id,
  patch
});

// epics
const createListEpic = ajaxEpic(
  CREATE,
  action => client.createList(action),
  listSchema
);

const closeCreateListModalOnSuccessEpic = action$ => action$.ofType(successOf(CREATE))
  .map(hideModal);

const updateListEpic = ajaxEpic(
  UPDATE,
  action => client.updateList(action),
  listSchema
);

const deleteListEpic = ajaxEpic(
  DELETE,
  action => client.deleteList(action.id).map(() => ({
    ...action,
    type: successOf(DELETE),
  }))
);

const onListDeleteEpic = action$ => action$.ofType(successOf(DELETE))
  .map(() => hideModal());

export const epics = combineEpics(
  createListEpic,
  closeCreateListModalOnSuccessEpic,
  updateListEpic,
  deleteListEpic,
  onListDeleteEpic
);

// reducer
const byID = (state = {}, action) => {
  switch (action.type) {
    case successOf(FETCH_BOARD): {
      return action.response.entities.lists || {};
    }
    case successOf(CREATE): {
      return {
        ...state,
        [action.response.result]: action.response.entities.lists[action.response.result]
      };
    }
    case successOf(UPDATE): {
      return {
        ...state,
        ...action.response.entities.lists[action.response.result]
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
      return action.response.entities.boards[action.response.result].lists || [];
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

export const getListSelector = (state, id) => state[DOMAIN].byID[id];
export const getListIDsSelector = (state) => state[DOMAIN].ids;
