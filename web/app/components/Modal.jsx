import React from 'react';
import { Modal } from 'reactstrap';
import { hashHistory } from 'react-router';

export class MainModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      modal: true
    };

    this.toggle = this.toggle.bind(this);
  }

  toggle() {
    const { board, editLevel, decrementEditLevel, disableAllEditModes } = this.props;

    disableAllEditModes();

    if (editLevel > 0) {
      decrementEditLevel();
      return;
    }

    this.setState({
      modal: !this.state.modal
    });

    return board ?
      setTimeout(() => hashHistory.push(`/b/${board.id}`), 500) :
      setTimeout(() => hashHistory.push('/'), 500);
  }

  render() {
    const childrenWithProps = React.Children.map(this.props.children,
     (child) => React.cloneElement(child, {
       ...child.props,
       toggle: this.toggle
     })
    );

    return (
      <div>
        <Modal
          isOpen={this.state.modal}
          toggle={this.toggle}
          className={this.props.className}
        >
          { childrenWithProps }
        </Modal>
      </div>
    );
  }
}


export const ConfirmModal = ({ isOpen = false, onConfirm, text, toggle, delay }) =>
  <Modal isOpen={isOpen} toggle={toggle}>
    <div className="modal-body">
      <div className="row">
        <div className="col-xs-5 offset-xs-1">
          <button
            type="button"
            className="btn btn-danger btn-block"
            onClick={() => {
              toggle();
              return delay ? setTimeout(onConfirm, delay) : onConfirm();
            }}
          >
            {text}
          </button>
        </div>
        <div className="col-xs-5">
          <button type="button" className="btn btn-secondary btn-block" onClick={toggle}>Cancel</button>
        </div>
      </div>
    </div>
  </Modal>;

export const InputModal = ({ isOpen, onSave, hideModal, placeholder }) => {
  let input;

  return (
    <Modal isOpen={isOpen} toggle={hideModal}>
      <div className="modal-body">
        <div className="row">
          <div className="col-xs-8">
            <input
              type="text"
              className="form-control"
              placeholder={placeholder}
              ref={node => {
                input = node;
                setTimeout(() => input && input.focus(), 0);
              }}
            />
          </div>
          <div className="col-xs-3 offset-xs-1 clearfix">
            <button type="button" className="close" aria-label="Close" onClick={hideModal}>
              <span aria-hidden="true">×</span>
            </button>
            <button
              type="button"
              className="btn btn-secondary btn-sm float-xs-right"
              onClick={(e) => {
                e.preventDefault();
                if (!input.value.trim()) {
                  return;
                }
                onSave(input.value);
              }}
            >Save</button>
          </div>
        </div>
      </div>
    </Modal>
  );
};

export const CreateBoardModal = ({ createBoard, ...rest }) =>
  <InputModal
    onSave={name => createBoard({ name })}
    placeholder="Board name"
    {...rest}
  />;

export const DeleteBoardModal = ({ id, isOpen, deleteBoard, hideModal }) =>
  <ConfirmModal
    isOpen={isOpen}
    onConfirm={() => {
      deleteBoard(id);
        // TODO change route on delete success
      hashHistory.push('/');
    }}
    text="Delete the board"
    toggle={hideModal}
  />;

export const CreateListModal = ({ createList, board, ...rest }) =>
  <InputModal
    onSave={name => createList({ boardID: board.id, name })}
    placeholder="List name"
    {...rest}
  />;

export const DeleteListModal = ({ list, isOpen, deleteList, hideModal }) =>
  <ConfirmModal
    isOpen={isOpen}
    onConfirm={() => deleteList(list)}
    text="Delete the list"
    toggle={hideModal}
  />;

export const DeleteCardModal = ({ card, isOpen, deleteCard, hideModal }) =>
  <ConfirmModal
    isOpen={isOpen}
    onConfirm={() => {
      deleteCard(card);
      hashHistory.push(`/b/${card.boardID}`);
    }}
    text="Delete the card"
    toggle={hideModal}
  />;