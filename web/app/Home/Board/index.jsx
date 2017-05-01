import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Button } from 'reactstrap';
import { DragDropContext } from 'react-dnd';
import HTML5Backend from 'react-dnd-html5-backend';
import { getActiveBoard } from '@/Home/duck';
import { showModal } from '@/components/Modal/duck';
import BoardHeader from './Header';
import DragDropContainer from './DragDropContainer';
import List from './List';
import { getBoardSelector, fetchBoard } from './duck';
import { getListIDsSelector, MODAL_CREATE_LIST } from './List/duck';

@DragDropContext(HTML5Backend)
class Board extends Component {
  componentDidMount() {
    const { fetch, activeBoard } = this.props;
    if (!activeBoard) {
      return;
    }

    fetch(activeBoard.owner, activeBoard.slug);
  }

  componentDidUpdate(prevProps) {
    const { fetch, activeBoard = {} } = this.props;

    if (activeBoard.slug &&
       (!prevProps.activeBoard || activeBoard.slug !== prevProps.activeBoard.slug || activeBoard.owner !== prevProps.activeBoard.owner)) {
      fetch(activeBoard.owner, activeBoard.slug);
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
    activeBoard: getActiveBoard(state),
    board: getBoardSelector(state),
    lists: getListIDsSelector(state)
  }),
  {
    fetch: fetchBoard,
    onCreate: (board) => showModal(MODAL_CREATE_LIST, board)
  }
)(Board);
