const DOMAIN = 'pulpe/modal';

// types
export const SHOW_MODAL = `${DOMAIN}/showModal`;
export const HIDE_MODAL = `${DOMAIN}/hideModal`;

// action creators
export const showModal = (modalType, modalProps = {}) => ({
  type: SHOW_MODAL,
  modalType,
  modalProps
});

export const hideModal = (modalType, modalProps = {}) => ({
  type: HIDE_MODAL,
  modalType,
  modalProps
});

// reducer
const initialState = {
  modalType: null,
  modalProps: {}
};

const modals = (state = initialState, action) => {
  switch (action.type) {
    case SHOW_MODAL:
      return {
        modalType: action.modalType,
        modalProps: action.modalProps
      };
    case HIDE_MODAL:
      return initialState;
    default:
      return state;
  }
};


export default {
  [DOMAIN]: modals
};

export const getModalProps = (state) => state[DOMAIN].modalProps;
export const getModalType = (state) => state[DOMAIN].modalType;
