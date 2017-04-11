import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Button } from 'reactstrap';
import { getActiveBoardID } from '@/Home/duck';
import { showModal } from '@/components/Modal/duck';
import BoardHeader from './Header';
import DragDropContainer from './DragDropContainer';
import List from './List';
import { getBoardSelector, fetchBoard } from './duck';
import { getListIDsSelector, MODAL_CREATE_LIST } from './List/duck';

class Board extends Component {
  componentDidMount() {
    const { fetch, id } = this.props;
    if (!id) {
      return;
    }

    fetch(id);
  }

  componentDidUpdate(prevProps) {
    const { fetch, id } = this.props;
    if (id && id !== prevProps.id) {
      fetch(id);
    }
  }

  render() {
    if (!this.props.board.id) {
      return null;
    }

    return (
      <div className="plp-board">
        <BoardHeader {...this.props} />
        <BoardBody {...this.props} />
      </div>
    );
  }
}

const BoardBody = ({ board, lists = [], onCreate }) => (
  <div className="plp-board-content gridBoard-horizontal">
    <DragDropContainer
      className="plp-cards-list-wrapper"
      itemClassName="plp-cards-list-wrapper"
    >
      {lists.map((id) => (
        <List key={id} id={id} />
      ))}
    </DragDropContainer>

    <div className="plp-cards-list-wrapper">
      <div className="plp-cards-list">
        <Button color="success" size="sm" className="btn-block" onClick={() => onCreate(board)}>+ Create a new list</Button>
      </div>
    </div>
  </div>
);

export default connect(
  (state) => ({
    id: getActiveBoardID(state),
    board: getBoardSelector(state),
    lists: getListIDsSelector(state)
  }),
  {
    fetch: fetchBoard,
    onCreate: (board) => showModal(MODAL_CREATE_LIST, board)
  }
)(Board);
