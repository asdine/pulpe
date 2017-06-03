import React, { Component } from 'react';
import { connect } from 'react-redux';
import { DragSource, DropTarget } from 'react-dnd';
import { getEmptyImage } from 'react-dnd-html5-backend';
import { showModal } from '@/components/Modal/duck';
import Editable from '@/components/Editable';
import { getBoardSelector } from '@/Home/Board/duck';
import Draggable from './Draggable';
import { patchList, updateList, deleteList, getListSelector, MODAL_DELETE_LIST } from './duck';
import { getCardsByListIDSelector, patchCard, dropCard, createCard } from './Card/duck';
import DragDropContainer from './DragDropContainer';
import Card from './Card';

const itemSource = {
  beginDrag({ id, index }) {
    return { id, index };
  },

  isDragging(props, monitor) {
    return monitor.getItem().id === props.id;
  },

  endDrag(props, monitor) {
    const { id, index: idx } = monitor.getItem();

    if (!monitor.didDrop()) {
      props.moveItem(id, idx);
      return;
    }

    const { id: droppedOnId, index } = monitor.getDropResult();

    props.onDrop(id, droppedOnId, index);
  },
};

const itemTarget = {
  hover(props, monitor) {
    const { id: draggedId } = monitor.getItem();
    const { id: overId } = props;

    if (draggedId !== overId) {
      const { index: overIndex } = props.findItem(overId);
      props.moveItem(draggedId, overIndex);
    }
  },

  drop(props) {
    const { id, index } = props;
    return { id, index };
  }
};


@DropTarget('List', itemTarget, conn => ({
  connectDropTarget: conn.dropTarget(),
}))
@DragSource('List', itemSource, (conn, monitor) => ({
  connectDragSource: conn.dragSource(),
  connectDragPreview: conn.dragPreview(),
  isDragging: monitor.isDragging()
}))
class List extends Component {
  componentDidMount() {
    // Use empty image as a drag preview so browsers don't draw it
    // and we can draw whatever we want on the custom drag layer instead.
    this.props.connectDragPreview(getEmptyImage(), {
      // IE fallback: specify that we'd rather screenshot the node
      // when it already knows it's being dragged so we can hide it with CSS.
      captureDraggingState: true,
    });
  }

  render() {
    const { connectDragSource, connectDropTarget,
            isDragging, isDragged, ...rest } = this.props;

    const { preview } = rest;

    return connectDropTarget(
      <div className="board-list-item">
        <div className={`list-wrapper ${isDragged ? 'dragged' : ''} ${isDragging ? 'shadow' : ''}`}>
          {!preview ?
            connectDragSource(<div><Header {...rest} /></div>) :
            <Header {...rest} />
          }
          <Body {...rest} />
          <Footer {...rest} />
        </div>
      </div>
    );
  }
}

const listPreview = (props) =>
  <div className="board-list-item dragged" style={props.style}>
    <div className="list-wrapper">
      <Header {...props} preview="true" />
      <Body {...props} preview="true" />
      <Footer {...props} preview="true" />
    </div>
  </div>;

const Header = ({ list = {}, onChangeName, index }) =>
  <div className="list-top">
    <Editable
      className="list-top-edit"
      value={list.name}
      onSave={(value) => onChangeName({ id: list.id, name: value })}
    >
      <h3>{ list.name || `#${index + 1}` }</h3>
    </Editable>
  </div>;

const Body = ({ board = {}, list = {}, cards = [], preview, moveToList, onDropCard }) => {
  const children = cards.map((card) => (
    <Card key={card.id} id={card.id} card={card} board={board} list={list} />
  ));

  if (preview) {
    return <div>{children}</div>;
  }

  return (
    <DragDropContainer moveToList={moveToList} onDrop={onDropCard}>
      {children}
    </DragDropContainer>
  );
};

const Footer = (props) => {
  const { list = {}, cards = [], moveToList, preview } = props;

  if (preview) {
    return (
      <FooterActions {...props} />
    );
  }

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
  <div className="list-bottom">
    <Editable
      onSave={(name) => onCreateCard({
        listID: list.id,
        name,
      })}
    >
      <button
        className="btn btn-secondary btn-sm btn-new-card"
      >+ Add a new card</button>
    </Editable>

    <button
      className="btn btn-secondary btn-sm btn-delete-list"
      onClick={() => onDelete(list.id, cards)}
    >
      Delete
    </button>
  </div>;

const connector = connect(
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
    onCreateCard: (data) => {
      dispatch(createCard(data));
    },
    onDelete: (id, cards) =>
      cards.length > 0 ?
        dispatch(showModal(MODAL_DELETE_LIST, id)) :
        dispatch(deleteList(id)),
    moveToList: (patch) => {
      dispatch(patchCard(patch));
    },
    onDropCard: (card, index, canceled) => {
      dispatch(dropCard(card, index, canceled));
    }
  })
);
export default connector(List);

export const ListPreview = connector(listPreview);
