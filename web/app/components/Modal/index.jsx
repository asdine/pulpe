import React from 'react';
import ReactModal from 'react-modal';

const Modal = ({ children, ...rest }) => (
  <ReactModal
    {...rest}
    overlayClassName="modal fade show"
    className="modal-dialog"
  >
    {children}
  </ReactModal>
  );

export default Modal;
