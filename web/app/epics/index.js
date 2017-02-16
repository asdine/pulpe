import 'rxjs/add/operator/mergeMap';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/do';
import 'rxjs/add/operator/filter';
import 'rxjs/add/observable/of';
import 'rxjs/add/observable/empty';
import 'rxjs/add/observable/dom/ajax';
import { combineEpics } from 'redux-observable';
import client from './client';
import * as schema from './schema';
import * as ActionTypes from '../actions/types';
import ajaxEpic from './ajaxEpic';
import fixCardPositionEpic from './fixCardPosition';
import onSuccess from './onSuccess';

export default combineEpics(
  ajaxEpic(
    ActionTypes.FETCH_BOARDS,
    (action) => client.getBoards(action.filters),
    [schema.board]
  ),
  ajaxEpic(
    ActionTypes.FETCH_BOARD,
    action => client.getBoard(action.id),
    schema.board
  ),
  ajaxEpic(
    ActionTypes.CREATE_BOARD,
    action => client.createBoard(action),
    schema.board
  ),
  ajaxEpic(
    ActionTypes.UPDATE_BOARD,
    action => client.updateBoard(action),
    schema.board
  ),
  ajaxEpic(
    ActionTypes.DELETE_BOARD,
    action => client.deleteBoard(action.id).map(() => ({
      ...action,
      type: ActionTypes.successType(ActionTypes.DELETE_BOARD),
    }))
  ),
  ajaxEpic(
    ActionTypes.CREATE_LIST,
    action => client.createList(action),
    schema.list
  ),
  ajaxEpic(
    ActionTypes.UPDATE_LIST,
    action => client.updateList(action),
    schema.list
  ),
  ajaxEpic(
    ActionTypes.DELETE_LIST,
    action => client.deleteList(action.id).map(() => ({
      ...action,
      type: ActionTypes.successType(ActionTypes.DELETE_LIST),
    }))
  ),
  ajaxEpic(
    ActionTypes.FETCH_CARD,
    action => client.getCard(action.id),
    schema.card
  ),
  ajaxEpic(
    ActionTypes.CREATE_CARD,
    action => client.createCard(action),
    schema.card
  ),
  ajaxEpic(
    ActionTypes.UPDATE_CARD,
    action => client.updateCard(action),
    schema.card
  ),
  ajaxEpic(
    ActionTypes.DELETE_CARD,
    action => client.deleteCard(action.id).map(() => ({
      ...action,
      type: ActionTypes.successType(ActionTypes.DELETE_CARD),
    }))
  ),
  fixCardPositionEpic,
  onSuccess
);
