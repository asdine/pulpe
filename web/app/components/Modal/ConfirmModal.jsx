import React from 'react';
import { Modal } from 'reactstrap';

const ConfirmModal = ({ isOpen = false, onConfirm, text, toggle }) =>
  <Modal isOpen={isOpen} toggle={toggle}>
    <div className="modal-body">
      <div className="row">
        <div className="col-5 offset-1">
          <button
            type="button"
            className="btn btn-danger btn-block"
            onClick={onConfirm}
          >
            {text}
          </button>
        </div>
        <div className="col-5">
          <button type="button" className="btn btn-secondary btn-block" onClick={() => toggle()}>Cancel</button>
        </div>
      </div>
    </div>
  </Modal>;

export default ConfirmModal;
