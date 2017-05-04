import { schema } from 'normalizr';
import { combineReducers } from 'redux';
import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
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
export const DROP = `${DOMAIN}/drop`;

export const MODAL_CREATE_LIST = `${DOMAIN}/modalCreateList`;
export const MODAL_DELETE_LIST = `${DOMAIN}/modalDeleteList`;

// schemas
const listSchema = new schema.Entity('lists');

// action creators
export const createList = (boardID, name) => ({
  type: CREATE,
  boardID,
  name
});

export const updateList = ({ id, ...patch }) => ({
  type: requestOf(UPDATE),
  id,
  patch
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

export const dropList = (id, droppedOnId, index) => ({
  type: DROP,
  id,
  droppedOnId,
  index
});

// epics
const beforeCreateListEpic = (action$, store) => action$.ofType(CREATE)
  .map(action => {
    const { boardID } = action;
    const lists = getListsSelector(store.getState());
    let position = 1 << 16;

    if (lists.length !== 0) {
      position += lists[lists.length - 1].position;
    }

    return {
      ...action,
      position,
      type: requestOf(CREATE),
    };
  });

const createListEpic = ajaxEpic(
  CREATE,
  action => client.createList(action),
  listSchema
);

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

const onDropListEpic = (action$, store) => action$.ofType(DROP)
  .mergeMap(action => {
    const { id, index } = action;

    const patch = {
      id,
    };

    const list = getListSelector(store.getState(), id);
    const lists = getListsSelector(store.getState());

    if (lists[index].id === list.id) {
      return Observable.empty();
    }

    if (lists.length === 1) {
      patch.position = 1 << 16;
    } else if (index === 0) {
      const { position: nextPosition } = lists[0];
      patch.position = nextPosition / 2;
    } else if (index < lists.length - 1) {
      if (list.position > lists[index].position) {
        const { position: prevPosition } = lists[index - 1];
        const { position: nextPosition } = lists[index];
        patch.position = prevPosition + ((nextPosition - prevPosition) / 2);
      } else {
        const { position: prevPosition } = lists[index];
        const { position: nextPosition } = lists[index + 1];
        patch.position = prevPosition + ((nextPosition - prevPosition) / 2);
      }
    } else {
      const { position: prevPosition } = lists[index];
      patch.position = prevPosition + (1 << 16);
    }

    if (patch.position === list.position) {
      return Observable.empty();
    }

    return Observable.of(patchList(patch), updateList(patch));
  });

export const epics = combineEpics(
  beforeCreateListEpic,
  createListEpic,
  updateListEpic,
  deleteListEpic,
  onListDeleteEpic,
  onDropListEpic,
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
export const getListIDsSelector = (state) => getListsSelector(state).map(list => list.id);
export const getListsSelector = (state) =>
  state[DOMAIN].ids
    .map(id => getListSelector(state, id))
    .sort((a, b) => a.position > b.position);
