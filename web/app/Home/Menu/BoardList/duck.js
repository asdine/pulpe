import { schema } from 'normalizr';
import { browserHistory } from 'react-router';
import { combineEpics } from 'redux-observable';
import client from '../../../services/api/client';
import ajaxEpic, { successOf, requestOf } from '../../../services/api/ajaxEpic';
import { hideModal } from '../../../components/Modal/duck';
import { DELETE as BOARD_DELETE } from '../../Board/duck';

const DOMAIN = 'pulpe/home/menu/boardList';

// types
const FETCH = `${DOMAIN}/fetch`;
const CREATE = `${DOMAIN}/create`;
export const MODAL_CREATE_BOARD = `${DOMAIN}/modalCreateBoard`;

// schemas
const boardSchema = new schema.Entity('boards');

// action creators
export const fetchBoards = (filters) => ({
  type: requestOf(FETCH),
  filters
});

export const createBoard = (name) => ({
  type: requestOf(CREATE),
  name
});

// epics
const fetchBoardsEpic = ajaxEpic(
  FETCH,
  (action) => client.getBoards(action.filters),
  [boardSchema]
);

const createBoardEpic = ajaxEpic(
  CREATE,
  action => client.createBoard(action),
  boardSchema
);

const redirectOnBoardCreationEpic = action$ => action$.ofType(successOf(CREATE))
  .do((action) =>
    setTimeout(() =>
      browserHistory.push(`/${action.response.entities.boards[action.response.result].slug}`),
      450
    )
  )
  .mapTo({ type: '' });

const closeModalOnCreationEpic = action$ => action$.ofType(successOf(CREATE))
  .map(hideModal);

export const epics = combineEpics(
  fetchBoardsEpic,
  createBoardEpic,
  closeModalOnCreationEpic,
  redirectOnBoardCreationEpic,
);

// reducer
const reducer = (state = [], action = {}) => {
  switch (action.type) {
    case successOf(FETCH): {
      return action.response.result.map(id => action.response.entities.boards[id]);
    }
    case successOf(CREATE): {
      return [
        ...state,
        action.response.entities.boards[action.response.result]
      ];
    }
    case successOf(BOARD_DELETE): {
      const idx = state.findIndex(b => b.id === action.id);
      if (idx === -1) {
        return state;
      }

      return [
        ...state.slice(0, idx),
        ...state.slice(idx + 1)
      ];
    }
    default:
      return state;
  }
};

export default {
  [DOMAIN]: reducer
};

export const getBoards = (state) => state[DOMAIN];
