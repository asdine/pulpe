import React from 'react';
import { connect } from 'react-redux';
import InputModal from '@/components/Modal/InputModal';
import { getModalType, hideModal } from '@/components/Modal/duck';
import * as duck from './duck';

const CreateBoardModal = (props) =>
  <InputModal
    placeholder="Board name"
    {...props}
  />;

export default connect(
  (state) => ({
    isOpen: getModalType(state) === duck.MODAL_CREATE_BOARD
  }),
  {
    hideModal,
    onSave: duck.createBoard,
  },
)(CreateBoardModal);
