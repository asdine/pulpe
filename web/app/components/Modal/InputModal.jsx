import React from 'react';
import { Modal } from 'reactstrap';

const save = (onSave, input) => {
  if (!input.value.trim()) {
    return;
  }

  onSave(input.value);
};

const InputModal = ({ isOpen, onSave, hideModal, placeholder }) => {
  let input;

  return (
    <Modal isOpen={isOpen} toggle={hideModal}>
      <div className="modal-body">
        <div className="row">
          <div className="col-8">
            <input
              type="text"
              className="form-control"
              placeholder={placeholder}
              ref={node => {
                input = node;
                setTimeout(() => input && input.focus(), 0);
              }}
              onKeyPress={(e) => {
                if (e.key === 'Enter') {
                  save(onSave, input);
                }
              }}
            />
          </div>
          <div className="col-3 offset-1 input-modal-options">
            <button
              type="button"
              className="btn btn-secondary btn-sm btn-save"
              onClick={(e) => {
                e.preventDefault();
                save(onSave, input);
              }}
            >Save</button>
            <button type="button" className="close" aria-label="Close" onClick={() => hideModal()}>
              <span aria-hidden="true">×</span>
            </button>
          </div>
        </div>
      </div>
    </Modal>
  );
};

export default InputModal;
