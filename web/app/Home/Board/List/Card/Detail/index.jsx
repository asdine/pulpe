import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Button, Modal, ModalBody } from 'reactstrap';
import { browserHistory } from 'react-router';
import { getModalProps, getModalType, hideModal } from '@/components/Modal/duck';
import Sub, { SubOpener, SubClosed, SubOpened } from '@/components/Sub';
import { subsStillOpened, popSub, closeSub, closeAllSubs, getLastOpened } from '@/components/Sub/duck';
import { getBoardSelector } from '@/Home/Board/duck';
import { MODAL_CARD_DETAIL, fetchCard, updateCard, patchCard, deleteCard, getCardBySlugSelector } from '@/Home/Board/List/Card/duck';

const DetailModal = (props) => {
  const { isOpen, toggle } = props;

  const close = () => toggle(props);

  return (
    <Modal
      isOpen={isOpen}
      toggle={close}
    >
      <Detail close={close} {...props} />
    </Modal>
  );
};

class Detail extends Component {
  componentDidMount() {
    const { fetch, card = {} } = this.props;
    if (!card.id) {
      return;
    }

    fetch(card.id);
  }

  componentWillReceiveProps(nextProps) {
    const { fetch, card = {} } = nextProps;

    if (card.id && (!this.props.card || card.id !== this.props.card.id)) {
      fetch(card.id);
    }
  }

  render() {
    const { card = {}, init } = this.props;
    if (!card.id) {
      return null;
    }

    return (
      <div className="plp-modal-large-card" onClick={init}>
        <Header {...this.props} />
        <Body {...this.props} />
      </div>
    );
  }
}

const Header = ({ card, close, onSave, onDelete }) =>
  <div className="modal-header">
    <div className="modal-title">
      <NameEditor card={card} onSave={onSave} close={close} />
    </div>

    <div className="modal-options clearfix">
      <Button className="close" data-dismiss="modal" aria-label="Close" onClick={close}>
        <span aria-hidden="true">&times;</span>
      </Button>
      <Button
        color="danger"
        size="sm"
        className="float-xs-right"
        onClick={() => {
          onDelete(card.id);
        }}
      >
        Delete
      </Button>
    </div>
  </div>;

const Body = (props) =>
  <ModalBody>
    <DescEditor {...props} />
  </ModalBody>;

const NameEditor = ({ card = {}, onSave, close }) => {
  let input;

  const save = () => {
    const value = input.value.trim();

    if (value && value !== card.name) {
      onSave(card.id, { name: value });
    } else {
      onSave(card.id, null);
    }
  };

  const onKeyPress = (e) => {
    if (e.key === 'Enter') {
      save();
      close();
    }
  };

  return (
    <Sub name="nameEditor" className="card-name">
      <SubClosed>
        <SubOpener>
          <h5>{ card.name }</h5>
        </SubOpener>
      </SubClosed>
      <SubOpened>
        <input
          type="text"
          autoFocus
          defaultValue={card.name}
          onBlur={save}
          onKeyPress={onKeyPress}
          ref={(node) => { input = node; }}
        />
      </SubOpened>
    </Sub>
  );
};

const DescEditor = ({ card = {}, close, onSave }) => {
  let input;

  const save = () => {
    const value = input.value.trim();

    if (value && value !== card.description) {
      onSave(card.id, { description: value });
    } else {
      onSave(card.id, null);
    }

    close();
  };

  const onKeyPress = (e) => {
    if (e.key === 'Enter') {
      save();
      close();
    }
  };

  return (
    <Sub name="descEditor">
      <SubClosed>
        <SubOpener>
          <div className="large-card-description">
            {card.description || <div className="large-card-description__no-description">Click here to add content</div>}
          </div>
        </SubOpener>
      </SubClosed>
      <SubOpened>
        <div className="large-card-description-edit">
          <textarea
            autoFocus
            defaultValue={card.description}
            onKeyPress={onKeyPress}
            ref={(node) => { input = node; }}
          />
          <div className="large-card-description-edit__footer">
            <button type="button" className="btn btn-secondary cancel-btn" onClick={close}>Cancel</button>
            <button type="button" className="btn btn-primary save-btn" onClick={save}>Save</button>
          </div>
        </div>
      </SubOpened>
    </Sub>
  );
};

export default connect(
  (state) => {
    const card = getCardBySlugSelector(state, getModalProps(state).card);

    return ({
      card,
      board: getBoardSelector(state),
      isOpen: getModalType(state) === MODAL_CARD_DETAIL,
      stillOpened: subsStillOpened(state),
      lastOpened: getLastOpened(state),
    });
  },
  (dispatch) => ({
    toggle: ({ stillOpened, lastOpened, board }) => {
      if (stillOpened) {
        dispatch(popSub());
        dispatch(closeSub(lastOpened));
      } else {
        dispatch(hideModal());
        browserHistory.push(`/${board.owner.login}/${board.slug}`);
      }
    },
    init: () => dispatch(closeAllSubs()),
    fetch: (id) => dispatch(fetchCard(id)),
    onSave: (id, patch) => {
      if (patch) {
        dispatch(updateCard({ id, ...patch }));
        dispatch(patchCard({ id, ...patch }));
      }
    },
    onDelete: (id) => dispatch(deleteCard(id))
  }),
)(DetailModal);
