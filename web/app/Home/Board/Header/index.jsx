import React from 'react';
import { connect } from 'react-redux';
import { showModal } from '@/components/Modal/duck';
import Editable from '@/components/Editable';
import { fetchBoard, updateBoard, patchBoard } from '@/Home/Board/duck';
import { MODAL_DELETE_BOARD } from './duck';

const BoardHeader = ({ board = {}, onSave, onDelete }) => (
  <header>
    <Editable
      className="board-name"
      childrenClassName="board-name-title"
      editorClassName="board-name-edit"
      value={board.name}
      onSave={(value) => onSave({ id: board.id, name: value })}
      editor={Editor}
    >
      <h2>{board.name}</h2>
    </Editable>
    <div className="board-options">
      <button
        className=""
        onClick={() => onDelete(board)}
      >Delete</button>
    </div>
  </header>
);

const Editor = ({ onRef, ...rest }) => (
  <input
    type="text"
    ref={onRef}
    {...rest}
  />
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

