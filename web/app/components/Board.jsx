import React from 'react';
import { Button } from 'reactstrap';
import List from '../containers/list';
import * as ActionTypes from '../actions/types';

const Board = ({ board = {}, children, ...rest }) => {
  const mode = board.settings && board.settings.mode === 'vertical' ? 'vertical' : 'horizontal';

  return (
    <div className="plp-board">
      <BoardHeader mode={mode} board={board} {...rest} />
      <BoardContent mode={mode} board={board} {...rest} />
      {children}
    </div>
  );
};

const BoardHeader = ({ board = {}, mode, updateBoard, showModal, toggleEditMode, isEditing }) => {
  let input;

  const save = () => {
    if (!input.value.trim() || input.value === board.name) {
      toggleEditMode('board-name');
      return;
    }

    updateBoard({ id: board.id, name: input.value });
    toggleEditMode('board-name');
  };

  return (
    <header>
      { !isEditing ?
        <h6 onClick={() => toggleEditMode('board-name')}>{board.name}</h6> :
        <input
          type="text"
          autoFocus
          className="board-name"
          ref={node => { input = node; }}
          defaultValue={board.name}
          onBlur={save}
          onKeyPress={(e) => {
            if (e.key === 'Enter') {
              save();
            }
          }}
        />
      }
      <Button
        color="danger"
        size="sm"
        className="float-xs-right"
        onClick={() => showModal(ActionTypes.MODAL_DELETE_BOARD, board)}
      >Delete</Button>
      <BoardModeSwitcher
        mode={mode}
        board={board}
        updateBoard={updateBoard}
      />
    </header>
  );
};

const BoardModeSwitcher = ({ board = {}, updateBoard, mode }) => {
  const changeMode = (newMode) => {
    if (board.settings && newMode === board.settings.mode) {
      return;
    }

    const settings = {
      ...board.settings,
      mode: newMode
    };
    updateBoard({ id: board.id, settings });
  };

  return (
    <div className="btn-group float-xs-right" data-toggle="buttons">
      <label
        className={`btn btn-secondary btn-sm ${mode === 'horizontal' ? 'active' : ''}`}
        htmlFor="gridBoard-horizontal-mode"
      >
        <input
          type="radio"
          name="horizontal"
          id="gridBoard-horizontal-mode"
          onClick={() => changeMode('horizontal')}
          defaultChecked={mode === 'horizontal'}
        /> Horizontal
      </label>
      <label
        className={`btn btn-secondary btn-sm ${mode === 'vertical' ? 'active' : ''}`}
        htmlFor="gridBoard-vertical-mode"
      >
        <input
          type="radio"
          name="options"
          id="gridBoard-vertical-mode"
          onClick={() => changeMode('vertical')}
          defaultChecked={mode === 'vertical'}
        /> Vertical
    </label>
    </div>
  );
};

const BoardContent = ({ mode, lists, board, showModal }) => {
  const modeClass = mode === 'horizontal' ? 'gridBoard-horizontal' : 'gridBoard-vertical';

  return (
    <div className={`plp-board-content ${modeClass}`}>
      {lists.map((list, i) => (
        <List key={list.id} id={list.id} index={i} mode={mode} />
      ))}
      <div className="plp-cards-list-wrapper">
        <div className="plp-cards-list">
          <Button color="success" size="sm" className="btn-block" onClick={() => showModal(ActionTypes.MODAL_CREATE_LIST, board)}>+ Create a new list</Button>
        </div>
      </div>
    </div>
  );
};

export default Board;
