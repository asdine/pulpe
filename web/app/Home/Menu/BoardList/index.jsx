import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import Editable from '@/components/Editable';
import { getActiveBoard } from '@/Home/duck';
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
        <ul className="list-unstyled left-menu__board-list">
          {boards.map((board) => (
            activeBoard.slug === board.slug ?
              <li key={board.id} className="left-menu__item is-active">
                {board.name}
              </li> :
              <li key={board.id} className="left-menu__item">
                <Link to={`/${board.owner.login}/${board.slug}`}>{board.name}</Link>
              </li>
          ))}
        </ul>
        <Editable editorClassName="left-menu__create-input" onSave={onCreate} >
          <button className="left-menu__create-button">Create a new board</button>
        </Editable>
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
    onCreate: duck.createBoard,
    fetchBoards: duck.fetchBoards,
  },
)(BoardList);
