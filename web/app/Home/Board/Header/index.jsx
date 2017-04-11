import React from 'react';
import { connect } from 'react-redux';
import { Button } from 'reactstrap';
import { showModal } from '@/components/Modal/duck';
import Editable from '@/components/Editable';
import { fetchBoard, updateBoard, patchBoard } from '@/Home/Board/duck';
import { MODAL_DELETE_BOARD } from './duck';

const BoardHeader = ({ board = {}, onSave, onDelete }) => (
  <header>
    <Editable
      className="board-name"
      value={board.name}
      onSave={(value) => onSave({ id: board.id, name: value })}
    >
      <h2>{board.name}</h2>
    </Editable>
    <div className="board-options">
      <Button
        color="danger"
        size="sm"
        onClick={() => onDelete(board)}
      >Delete</Button>
    </div>
  </header>
);

export default connect(
  null,
  (dispatch) => ({
    fetch: (id) => dispatch(fetchBoard(id)),
    onSave: (patch) => {
      dispatch(patchBoard(patch));
      dispatch(updateBoard(patch));
    },
    onDelete: (board) => dispatch(showModal(MODAL_DELETE_BOARD, board))
  })
)(BoardHeader);
