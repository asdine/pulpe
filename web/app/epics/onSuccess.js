import { combineEpics } from 'redux-observable';
import { browserHistory } from 'react-router';
import { successType, CREATE_BOARD, CREATE_LIST } from '../actions/types';
import { hideModal } from '../actions';

const onCreateBoard = action$ => action$.ofType(successType(CREATE_BOARD))
  .do((action) => browserHistory.push(`/b/${action.response.result}`))
  .map(() => hideModal());

const onCreateList = action$ => action$.ofType(successType(CREATE_LIST))
  .map(() => hideModal());

export default combineEpics(
  onCreateBoard,
  onCreateList,
);
