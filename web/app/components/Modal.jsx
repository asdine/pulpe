import React from 'react';
import { Modal } from 'reactstrap';
import { browserHistory } from 'react-router';
import { LargeCreate } from './Card';
import { MainModal as ContainerModal } from '../containers/modal';

export class MainModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      modal: props.isOpen !== undefined ? props.isOpen : true
    };

    this.toggle = this.toggle.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      modal: nextProps.isOpen !== undefined ? nextProps.isOpen : true
    });
  }

  toggle() {
    const { onExit = () => {}, editLevel,
            decrementEditLevel, disableAllEditModes, hideModal } = this.props;

    disableAllEditModes();

    if (editLevel > 0) {
      decrementEditLevel();
      return;
    }

    this.setState({
      modal: !this.state.modal
    });

    hideModal();

    setTimeout(() => {
      onExit();
    }, 200);
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
        <div className="col-5 offset-1">
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
        <div className="col-5">
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
          <div className="col-8">
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
          <div className="col-3 offset-1 input-modal-options">
            <button
              type="button"
              className="btn btn-secondary btn-sm btn-save"
              onClick={(e) => {
                e.preventDefault();
                if (!input.value.trim()) {
                  return;
                }
                onSave(input.value);
              }}
            >Save</button>
            <button type="button" className="close" aria-label="Close" onClick={hideModal}>
              <span aria-hidden="true">×</span>
            </button>
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

export const DeleteBoardModal = ({
    id,
    isOpen,
    deleteBoard,
    hideModal,
    redirectTo }) => (
      <ConfirmModal
        isOpen={isOpen}
        onConfirm={() => {
          deleteBoard(id);
          return redirectTo !== undefined ?
            browserHistory.push(`${redirectTo}`) :
            browserHistory.push('/');
        }}
        text="Delete the board"
        toggle={hideModal}
      />);

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

export const CreateCardModal = ({ isOpen, hideModal, ...rest }) =>
  <ContainerModal
    isOpen={isOpen}
    text="Create a card"
    toggle={hideModal}
  >
    <LargeCreate hideModal={hideModal} {...rest} />
  </ContainerModal>;

export const DeleteCardModal = ({ card, redirectTo, isOpen, deleteCard, hideModal }) =>
  <ConfirmModal
    isOpen={isOpen}
    onConfirm={() => {
      deleteCard(card);
      browserHistory.push(`/${redirectTo}`);
    }}
    text="Delete the card"
    toggle={hideModal}
  />;
