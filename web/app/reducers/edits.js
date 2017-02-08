import { combineReducers } from 'redux';
import { SET_EDIT_MODE, TOGGLE_EDIT_MODE, INCREMENT_EDIT_LEVEL, DECREMENT_EDIT_LEVEL, SET_EDIT_LEVEL, DISABLE_ALL_EDIT_MODES } from '../actions/types';

const editLevel = (state = 0, action) => {
  switch (action.type) {
    case INCREMENT_EDIT_LEVEL: {
      return state + 1;
    }
    case DECREMENT_EDIT_LEVEL: {
      return state - 1;
    }
    case SET_EDIT_LEVEL: {
      return action.value;
    }
    default:
      return state;
  }
};

const editModes = (state = {}, action) => {
  switch (action.type) {
    case SET_EDIT_MODE: {
      return {
        ...state,
        [action.item]: action.status
      };
    }
    case TOGGLE_EDIT_MODE: {
      return {
        ...state,
        [action.item]: !state[action.item]
      };
    }
    case DISABLE_ALL_EDIT_MODES: {
      return {};
    }
    default:
      return state;
  }
};

export default combineReducers({
  editLevel,
  editModes,
});

export const isEditing = (state, item) => state.editModes[item];
export const getEditLevel = (state) => state.editLevel;
