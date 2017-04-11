import React from 'react';
import { connect } from 'react-redux';
import InputModal from '../../../components/Modal/InputModal';
import { getModalType, hideModal } from '../../../components/Modal/duck';
import { getBoardSelector } from '../duck';
import { createList, MODAL_CREATE_LIST } from './duck';


const CreateListModal = ({ onSave, board, ...rest }) =>
  <InputModal
    placeholder="List name"
    onSave={(name) => {
      onSave(board, name);
    }}
    {...rest}
  />;

export default connect(
  (state) => ({
    isOpen: getModalType(state) === MODAL_CREATE_LIST,
    board: getBoardSelector(state),
  }),
  {
    hideModal,
    onSave: (board = {}, name) => createList(board.id, name),
  },
)(CreateListModal);
