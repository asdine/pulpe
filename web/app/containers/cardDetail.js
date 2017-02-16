import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { MainModal } from './modal';
import { Large } from '../components/Card';
import { getCardByID, isEditing, getBoardByID } from '../reducers';
import * as actions from '../actions';

const mapStateToProps = (state, { params }) => {
  const card = getCardByID(state, params.id) || {};
  return ({
    card,
    board: getBoardByID(state, card.boardID),
    isEditingName: isEditing(state, 'card-name'),
    isEditingDesc: isEditing(state, 'card-desc')
  });
};

class CardDetail extends React.Component {
  constructor(props) {
    super(props);

    this.saveName = this.saveName.bind(this);
    this.saveDesc = this.saveDesc.bind(this);
    this.onExit = this.onExit.bind(this);
  }

  componentDidMount() {
    const { fetchCard, params } = this.props;
    fetchCard(params.id);
  }

  componentDidUpdate(prevProps) {
    const { fetchCard, params } = this.props;
    if (prevProps.params.id !== params.id) {
      fetchCard(params.id);
    }
  }

  onExit() {
    browserHistory.push(`/${this.props.board.slug || ''}`);
  }

  saveName(input) {
    const { card, updateCard, disableAllEditModes } = this.props;
    const name = input.value.trim();

    if (!name) {
      disableAllEditModes();
      return;
    }

    const update = {
      id: card.id
    };
    let updated = false;

    if (card.name !== name) {
      updated = true;
      update.name = name;
    }

    if (updated) {
      updateCard(update);
    }

    disableAllEditModes();
  }

  saveDesc(textarea) {
    const { card, updateCard, disableAllEditModes } = this.props;
    const description = textarea.value.trim();

    const update = {
      id: card.id
    };
    let updated = false;

    if (card.description !== description) {
      updated = true;
      update.description = description;
    }

    if (updated) {
      updateCard(update);
    }

    disableAllEditModes();
  }

  render() {
    return (
      <MainModal onExit={this.onExit}>
        <Large
          saveName={this.saveName}
          saveDesc={this.saveDesc}
          {...this.props}
        />
      </MainModal>
    );
  }
}

export default connect(
  mapStateToProps,
  actions,
)(CardDetail);
