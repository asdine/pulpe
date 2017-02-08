import { combineReducers } from 'redux';
import context, * as fromContext from './context';
import boards, * as fromBoards from './boards';
import lists, * as fromLists from './lists';
import cards, * as fromCards from './cards';
import modals, * as fromModals from './modals';
import edits, * as fromEdits from './edits';

const reducers = {
  context,
  boards,
  lists,
  cards,
  modals,
  edits,
};

if (process.env.NODE_ENV === 'development') {
  const dev = require('./dev'); // eslint-disable-line global-require
  reducers.dev = dev.default;
  module.exports.getDevStore = (state) => dev.getDevStore(state.dev);
}

export default combineReducers(reducers);

export const getBoardByID = (state, id) =>
  fromBoards.getBoardByID(state.boards, id);

export const getBoards = (state) =>
  fromBoards.getBoards(state.boards);

export const getFirstBoardID = (state) =>
  fromBoards.getFirstBoardID(state.boards);

export const getListsByBoardID = (state, boardID) =>
  fromLists.getListsByBoardID(state.lists, boardID);

export const getListByID = (state, id) =>
  fromLists.getListByID(state.lists, id);

export const getCardByID = (state, id) =>
  fromCards.getCardByID(state.cards, id);

export const getCardsByListID = (state, boardID, listID) =>
  fromCards.getCardsByListID(state.cards, boardID, listID);

export const getActiveBoardID = (state) =>
  fromContext.getActiveBoardID(state.context);

export const getActiveBoard = (state) =>
  getBoardByID(state, getActiveBoardID(state));

export const getModalProps = (state) => fromModals.getModalProps(state.modals);
export const getModalType = (state) => fromModals.getModalType(state.modals);

export const isEditing = (state, item) => fromEdits.isEditing(state.edits, item);
export const getEditLevel = (state) => fromEdits.getEditLevel(state.edits);
