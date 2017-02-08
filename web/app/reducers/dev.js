import { DEV_SAVE } from '../actions/types';

const devStore = (state = {}, action) => {
  switch (action.type) {
    case DEV_SAVE: {
      return {
        ...state,
        ...action.update
      };
    }
    default:
      return state;
  }
};

export default devStore;

export const getDevStore = (state) => state;
