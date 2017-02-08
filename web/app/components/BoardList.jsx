import React from 'react';
import { Link } from 'react-router';
import { Button } from 'reactstrap';
import * as ActionTypes from '../actions/types';

const BoardList = ({ boards = [], showModal, activeBoard = {} }) =>
  <div>
    <ul className="list-unstyled plp-boards-list">
      {boards.map((board) => (
        activeBoard.id === board.id ?
          <li key={board.id} className="is-active">
            {board.name}
          </li> :
          <li key={board.id}>
            <Link to={`/b/${board.id}`}>{board.name}</Link>
          </li>
      ))}
    </ul>
    <Button color="secondary" size="sm" onClick={() => showModal(ActionTypes.MODAL_CREATE_BOARD)}>+ Create a board</Button>
  </div>;

export default BoardList;
