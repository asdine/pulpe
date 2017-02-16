if (process.env.NODE_ENV === 'development') {
  const DEV_SAVE = 'DEV_SAVE';
  module.exports.DEV_SAVE = DEV_SAVE;
}

export const FETCH_BOARDS = 'FETCH_BOARDS';
export const FETCH_BOARD = 'FETCH_BOARD';
export const CREATE_BOARD = 'CREATE_BOARD';
export const UPDATE_BOARD = 'UPDATE_BOARD';
export const DELETE_BOARD = 'DELETE_BOARD';

export const CREATE_LIST = 'CREATE_LIST';
export const UPDATE_LIST = 'UPDATE_LIST';
export const DELETE_LIST = 'DELETE_LIST';

export const FETCH_CARD = 'FETCH_CARD';
export const CREATE_CARD = 'CREATE_CARD';
export const UPDATE_CARD = 'UPDATE_CARD';
export const UPDATE_CARD_POSITION = 'UPDATE_CARD_POSITION';
export const DELETE_CARD = 'DELETE_CARD';

export const SET_ACTIVE_BOARD = 'SET_ACTIVE_BOARD';

export const SHOW_MODAL = 'SHOW_MODAL';
export const HIDE_MODAL = 'HIDE_MODAL';
export const MODAL_CREATE_BOARD = 'MODAL_CREATE_BOARD';
export const MODAL_DELETE_BOARD = 'MODAL_DELETE_BOARD';
export const MODAL_EDIT_BOARD = 'MODAL_EDIT_BOARD';
export const MODAL_CREATE_LIST = 'MODAL_CREATE_LIST';
export const MODAL_DELETE_LIST = 'MODAL_DELETE_LIST';
export const MODAL_CREATE_CARD = 'MODAL_CREATE_CARD';
export const MODAL_DELETE_CARD = 'MODAL_DELETE_CARD';

export const SET_EDIT_MODE = 'SET_EDIT_MODE';
export const TOGGLE_EDIT_MODE = 'TOGGLE_EDIT_MODE';
export const DISABLE_ALL_EDIT_MODES = 'DISABLE_ALL_EDIT_MODES';

export const INCREMENT_EDIT_LEVEL = 'INCREMENT_EDIT_LEVEL';
export const DECREMENT_EDIT_LEVEL = 'DECREMENT_EDIT_LEVEL';
export const SET_EDIT_LEVEL = 'SET_EDIT_LEVEL';

export const requestType = (type) => `${type}_REQUEST`;
export const successType = (type) => `${type}_SUCCESS`;
export const failureType = (type) => `${type}_FAILURE`;
