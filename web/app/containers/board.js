import React from 'react';
import { connect } from 'react-redux';
import { getActiveBoard, getBoardByID, getListsByBoardID, isEditing } from '../reducers';
import * as actions from '../actions';
import Board from '../components/Board';

const mapStateToProps = (state) => {
  const board = getActiveBoard(state);
  const id = board ? board.id : null;
  return ({
    id,
    board,
    lists: getListsByBoardID(state, id),
    isEditing: isEditing(state, 'board-name')
  });
};

class BoardDetail extends React.Component {
  componentDidMount() {
    const { fetchBoard, id } = this.props;
    if (!id) {
      return;
    }
    fetchBoard(id);
  }

  componentDidUpdate(prevProps) {
    const { fetchBoard, id } = this.props;
    if (id && id !== prevProps.id) {
      fetchBoard(id);
    }
  }

  render() {
    if (!this.props.board) {
      return <div />;
    }

    const { fetchBoard, ...rest } = this.props; /* eslint no-unused-vars: 0 */
    return (
      <Board {...rest} />
    );
  }
}

export default connect(
  mapStateToProps,
  actions
)(BoardDetail);
