import React from 'react';
import { connect } from 'react-redux';
import Modal from '@/components/Modal';
import { getModalProps, getModalType, hideModal } from '@/components/Modal/duck';
import { MODAL_CREATE_CARD, createCard, getCardsByListIDSelector } from './duck';

const CreateCardModal = (props) => {
  const { isOpen, toggle } = props;
  return (
    <Modal isOpen={isOpen} onRequestClose={toggle} contentLabel="Create a card">
      <CreateCardForm {...props} />
    </Modal>
  );
};

const CreateCardForm = (props) => {
  let input;
  let textarea;

  const { list = {}, onSave, toggle } = props;

  const save = () => {
    const name = input.value.trim();
    const description = textarea.value.trim();

    if (!name) {
      return;
    }

    const newCard = {
      listID: list.id,
      name,
      description
    };

    onSave(newCard);
  };

  return (
    <div className="modal-content">
      <div className="modal-header">
        <div className="container">
          <div className="row">
            <div className="col-8">
              <div className="form-group">
                <input
                  type="text"
                  className="form-control"
                  placeholder="Card name"
                  ref={node => {
                    input = node;
                    setTimeout(() => input && input.focus(), 0);
                  }}
                />
              </div>
            </div>
            <div className="col-3 offset-1 close-save-options">
              <button type="button" className="close" data-dismiss="modal" aria-label="Close" onClick={toggle}>
                <span aria-hidden="true">&times;</span>
              </button>
              <button
                type="button"
                className="btn btn-secondary btn-sm btn-save"
                onClick={(e) => {
                  e.preventDefault();
                  save(input, textarea);
                }}
              >
                Save
              </button>
            </div>
          </div>
        </div>
      </div>
      <div className="modal-body">
        <div className="form-group">
          <label htmlFor="card-content">Content</label>
          <textarea
            className="form-control"
            id="card-content"
            rows="3"
            ref={node => {
              textarea = node;
            }}
          />
        </div>
      </div>
    </div>
  );
};

export default connect(
  (state) => {
    const list = getModalProps(state);

    return ({
      list,
      cards: getCardsByListIDSelector(state, list.boardID, list.id),
      isOpen: getModalType(state) === MODAL_CREATE_CARD
    });
  },
  {
    toggle: hideModal,
    onSave: createCard,
  },
)(CreateCardModal);
