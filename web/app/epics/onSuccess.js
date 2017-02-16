import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
import { browserHistory } from 'react-router';
import { successType, CREATE_BOARD, UPDATE_BOARD, CREATE_LIST } from '../actions/types';
import { hideModal } from '../actions';

const onCreateBoard = action$ => action$.ofType(successType(CREATE_BOARD))
  .do((action) => browserHistory.push(`/${action.response.entities.boards[action.response.result].slug}`))
  .map(() => hideModal());

const onUpdateBoard = action$ => action$.ofType(successType(UPDATE_BOARD))
  .do((action) => browserHistory.push(`/${action.response.entities.boards[action.response.result].slug}`))
  .mergeMap(() => Observable.empty());

const onCreateList = action$ => action$.ofType(successType(CREATE_LIST))
  .map(() => hideModal());

export default combineEpics(
  onCreateBoard,
  onUpdateBoard,
  onCreateList,
);
