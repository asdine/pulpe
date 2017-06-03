import React, { Component } from 'react';
import { connect } from 'react-redux';
import { DragDropContext } from 'react-dnd';
import HTML5Backend from 'react-dnd-html5-backend';
import { Scrollbars } from 'react-custom-scrollbars';
import { getActiveBoard } from '@/Home/duck';
import Editable from '@/components/Editable';
import BoardHeader from './Header';
import DragDropContainer from './DragDropContainer';
import List from './List';
import { DraggablePreview } from './List/Draggable';
import { getBoardSelector, fetchBoard } from './duck';
import { getListIDsSelector, createList, dropList } from './List/duck';

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
      <div className="main-board">
        <BoardHeader {...this.props} />
        <BoardBody {...this.props} />
      </div>
    );
  }
}

const BoardBody = ({ board, lists = [], onCreate, onDrop }) => (
  <div className="board-area">
    <Scrollbars
      renderThumbHorizontal={props => <div {...props} className="thumb-horizontal" />}
      renderTrackHorizontal={props => <div {...props} className="track-horizontal" />}
    >
      <div className="scrollable board-content">
        <DragDropContainer
          className="plp-cards-list-wrapper"
          itemClassName="plp-cards-list-wrapper"
          onDrop={onDrop}
        >
          {lists.map((id) => (
            <List key={id} id={id} />
      ))}
        </DragDropContainer>
        <DraggablePreview board={board} lists={lists} />
        <div className="plp-cards-list-wrapper">
          <div className="plp-cards-list">
            <div className="plp-list-top">
              <Editable
                onSave={(name) => onCreate(board.id, name)}
                className="plp-list-top-edit"
              >
                <button className="btn btn-success btn-sm btn-block">+ Create a new list</button>
              </Editable>
            </div>
          </div>
        </div>
      </div>
    </Scrollbars>
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
    onCreate: createList,
    onDrop: dropList,
  }
)(Board);
