import React from 'react';
import { connect } from 'react-redux';
import ConfirmModal from '../../../components/Modal/ConfirmModal';
import { getModalType, getModalProps, hideModal } from '../../../components/Modal/duck';
import { deleteList, MODAL_DELETE_LIST } from './duck';

const DeleteListModal = ({ id, onDelete, ...rest }) =>
  <ConfirmModal
    onConfirm={() => onDelete(id)}
    text="Delete the list"
    {...rest}
  />;

export default connect(
  (state) => ({
    isOpen: getModalType(state) === MODAL_DELETE_LIST,
    id: getModalProps(state),
  }),
  {
    toggle: hideModal,
    onDelete: deleteList,
  },
)(DeleteListModal);
