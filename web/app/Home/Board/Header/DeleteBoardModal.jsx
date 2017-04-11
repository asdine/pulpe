import React from 'react';
import { connect } from 'react-redux';
import ConfirmModal from '../../../components/Modal/ConfirmModal';
import { getModalType, getModalProps, hideModal } from '../../../components/Modal/duck';
import { deleteBoard } from '../duck';
import { MODAL_DELETE_BOARD } from './duck';

const DeleteBoardModal = ({ id, onDelete, ...rest }) => (
  <ConfirmModal
    onConfirm={() => onDelete(id)}
    text="Delete the board"
    {...rest}
  />);

export default connect(
  (state) => ({
    id: getModalProps(state).id,
    isOpen: getModalType(state) === MODAL_DELETE_BOARD,
  }),
  {
    onDelete: deleteBoard,
    toggle: hideModal
  },
)(DeleteBoardModal);
