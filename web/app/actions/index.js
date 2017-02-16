import * as types from './types';

if (process.env.NODE_ENV === 'development') {
  module.exports.devSave = (k, v) => ({
    type: types.DEV_SAVE,
    update: {
      [k]: v,
    }
  });
}

export const fetchBoards = (filters) => ({
  type: types.requestType(types.FETCH_BOARDS),
  filters
});

export const createBoard = ({ name }) => ({
  type: types.requestType(types.CREATE_BOARD),
  name
});

export const fetchBoard = (id) => ({
  type: types.requestType(types.FETCH_BOARD),
  id
});

export const updateBoard = ({ id, ...rest }) => ({
  type: types.requestType(types.UPDATE_BOARD),
  id,
  toUpdate: rest
});

export const deleteBoard = (id) => ({
  type: types.requestType(types.DELETE_BOARD),
  id
});

export const createList = ({ boardID, name }) => ({
  type: types.requestType(types.CREATE_LIST),
  boardID,
  name
});

export const updateList = ({ id, ...rest }) => ({
  type: types.requestType(types.UPDATE_LIST),
  id,
  toUpdate: rest
});

export const deleteList = ({ id, boardID }) => ({
  type: types.requestType(types.DELETE_LIST),
  id,
  boardID,
});

export const fetchCard = (id) => ({
  type: types.requestType(types.FETCH_CARD),
  id
});

export const createCard = ({ boardID, listID, name, description, position }) => ({
  type: types.requestType(types.CREATE_CARD),
  boardID,
  listID,
  name,
  description,
  position
});

export const updateCard = ({ id, name, description, position }) => ({
  type: types.requestType(types.UPDATE_CARD),
  id,
  name,
  description,
  position
});

export const updateCardPosition = ({ id, position }) => ({
  type: types.UPDATE_CARD_POSITION,
  id,
  position
});

export const deleteCard = ({ id, boardID }) => ({
  type: types.requestType(types.DELETE_CARD),
  id,
  boardID
});

export const setActiveBoard = (slug) => ({
  type: types.SET_ACTIVE_BOARD,
  activeBoard: slug
});

export const showModal = (modalType, modalProps = {}) => ({
  type: types.SHOW_MODAL,
  modalType,
  modalProps
});

export const hideModal = (modalType, modalProps = {}) => ({
  type: types.HIDE_MODAL,
  modalType,
  modalProps
});

export const setEditMode = (item, status) => ({
  type: types.SET_EDIT_MODE,
  item,
  status
});

export const toggleEditMode = (item) => ({
  type: types.TOGGLE_EDIT_MODE,
  item,
});

export const disableAllEditModes = () => ({
  type: types.DISABLE_ALL_EDIT_MODES
});


export const incrementEditLevel = () => ({
  type: types.INCREMENT_EDIT_LEVEL
});

export const decrementEditLevel = () => ({
  type: types.DECREMENT_EDIT_LEVEL
});

export const setEditLevel = (value) => ({
  type: types.SET_EDIT_LEVEL,
  value
});
