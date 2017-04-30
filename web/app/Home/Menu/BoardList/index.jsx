import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { Button } from 'reactstrap';
import { getActiveBoard } from '@/Home/duck';
import { showModal } from '@/components/Modal/duck';
import * as duck from './duck';

class BoardList extends Component {
  componentDidMount() {
    const { fetchBoards } = this.props;
    fetchBoards();
  }

  render() {
    const { boards = [], onCreate, activeBoard = {} } = this.props;

    return (
      <div>
        <ul className="list-unstyled plp-boards-list">
          {boards.map((board) => (
            activeBoard.slug === board.slug ?
              <li key={board.id} className="is-active">
                {board.name}
              </li> :
              <li key={board.id}>
                <Link to={`/${board.owner.login}/${board.slug}`}>{board.name}</Link>
              </li>
          ))}
        </ul>
        <Button color="secondary" size="sm" onClick={onCreate}>+ Create a board</Button>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  boards: duck.getBoards(state),
  activeBoard: getActiveBoard(state),
});

export default connect(
  mapStateToProps,
  {
    onCreate: () => showModal(duck.MODAL_CREATE_BOARD),
    fetchBoards: duck.fetchBoards,
  },
)(BoardList);
