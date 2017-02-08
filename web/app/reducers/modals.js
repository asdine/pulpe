import * as ActionTypes from '../actions/types';

const initialState = {
  modalType: null,
  modalProps: {}
};

const modals = (state = initialState, action) => {
  switch (action.type) {
    case ActionTypes.SHOW_MODAL:
      return {
        modalType: action.modalType,
        modalProps: action.modalProps
      };
    case ActionTypes.HIDE_MODAL:
      return initialState;
    default:
      return state;
  }
};

export default modals;

export const getModalProps = (state) => state.modalProps;
export const getModalType = (state) => state.modalType;
