import React from 'react';
import { connect } from 'react-redux';
import { showModal } from '@/components/Modal/duck';
import Editable from '@/components/Editable';
import { getBoardSelector } from '@/Home/Board/duck';
import Draggable from './Draggable';
import { patchList, updateList, deleteList, getListSelector, MODAL_DELETE_LIST } from './duck';
import { getCardsByListIDSelector, patchCard, dropCard, MODAL_CREATE_CARD } from './Card/duck';
import DragDropContainer from './DragDropContainer';
import Card from './Card';

const List = (props) =>
  <div className="plp-cards-list-wrapper">
    <div className="plp-cards-list">
      <Header {...props} />
      <Body {...props} />
      <Footer {...props} />
    </div>
  </div>;

const Header = ({ connectDragSource, list = {}, onChangeName, index }) =>
  connectDragSource(
    <div className="plp-list-top">
      <Editable
        className="plp-list-top-edit"
        value={list.name}
        onSave={(value) => onChangeName({ id: list.id, name: value })}
      >
        <h3>{ list.name || `#${index + 1}` }</h3>
      </Editable>
    </div>
  );

const Body = ({ board = {}, list = {}, cards = [], moveToList, onDrop }) =>
  <DragDropContainer moveToList={moveToList} onDrop={onDrop}>
    {cards.map((card) => (
      <Card key={card.id} id={card.id} card={card} board={board} list={list} />
    ))}
  </DragDropContainer>;

const Footer = (props) => {
  const { list = {}, cards = [], moveToList } = props;

  return (
    <Draggable
      locked="true"
      id="addcard"
      list={list}
      cards={cards}
      moveToList={moveToList}
    >
      <FooterActions {...props} />
    </Draggable>
  );
};


const FooterActions = ({ list, onCreateCard, onDelete, cards }) =>
  <div className="plp-list-bottom">
    <button
      className="btn btn-secondary btn-sm btn-new-card"
      onClick={() => onCreateCard(list)}
    >+ Add a new card</button>
    <button
      className="btn btn-secondary btn-sm btn-delete-list"
      onClick={() => onDelete(list.id, cards)}
    >
      Delete
    </button>
  </div>;

export default connect(
  (state, { id }) => ({
    list: getListSelector(state, id),
    board: getBoardSelector(state),
    cards: getCardsByListIDSelector(state, id)
  }),
  (dispatch) => ({
    onChangeName: (patch) => {
      dispatch(patchList(patch));
      dispatch(updateList(patch));
    },
    onCreateCard: (list) => {
      dispatch(showModal(MODAL_CREATE_CARD, list));
    },
    onDelete: (id, cards) =>
      cards.length > 0 ?
        dispatch(showModal(MODAL_DELETE_LIST, id)) :
        dispatch(deleteList(id)),
    moveToList: (patch) => {
      dispatch(patchCard(patch));
    },
    onDrop: (card, index, canceled) => {
      dispatch(dropCard(card, index, canceled));
    }
  })
)(List);
