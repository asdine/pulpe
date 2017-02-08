import React from 'react';
import { findDOMNode } from 'react-dom';
import { Link } from 'react-router';
import { Button } from 'reactstrap';
import { DragSource, DropTarget, DragDropContext } from 'react-dnd';
import HTML5Backend from 'react-dnd-html5-backend';
import { Compact } from './Card';
import * as ActionTypes from '../actions/types';

class List extends React.Component {
  constructor(props) {
    super(props);

    this.onDrag = this.onDrag.bind(this);
    this.moveCard = this.moveCard.bind(this);
    this.onEndDrag = this.onEndDrag.bind(this);
    this.toggleEditMode = this.toggleEditMode.bind(this);

    this.state = {
      userIsDragging: false,
      isEditing: false
    };
  }

  onDrag(cardID, hoveredID, monitor, component) {
    this.setState({
      userIsDragging: true
    });

    const { cards = [] } = this.props;

    const card = cards.find((c) => c.id === cardID);
    const hovered = cards.find((c) => c.id === hoveredID);
    if (!card || !hovered) {
      return;
    }

    if (this.shouldSwitch(card, hovered, monitor, component)) {
      // Time to actually perform the action
      this.moveCard(card, hovered);
    }
  }

  onEndDrag(id) {
    this.setState({
      userIsDragging: false
    });

    const { cards = [], updateCard } = this.props;

    const card = cards.find((c) => c.id === id);

    updateCard({ id, position: card.position });
  }

  shouldSwitch(card, hovered, monitor, component) {
    // Don't replace items that aren't from the same list
    if (card.listID !== hovered.listID) {
      return false;
    }

    const { mode } = this.props;

    // Determine rectangle on screen
    const hoverBoundingRect = findDOMNode(component).getBoundingClientRect(); // eslint-disable-line

    // Get middle
    const hoverMiddle = mode === 'horizontal' ?
      (hoverBoundingRect.bottom - hoverBoundingRect.top) / 2 :
      (hoverBoundingRect.right - hoverBoundingRect.left) / 2;

    // Determine mouse position
    const clientOffset = monitor.getClientOffset();

    // Get pixels to the top
    const hoverClient = mode === 'horizontal' ?
      clientOffset.y - hoverBoundingRect.top :
      clientOffset.x - hoverBoundingRect.left;

    // Only perform the move when the mouse has crossed half of the items height or width
    // When dragging downwards, only move when the cursor is below 50%
    // When dragging upwards, only move when the cursor is above 50%

    // Dragging downwards
    if (card.position < hovered.position && hoverClient < hoverMiddle) {
      return false;
    }

    // Dragging upwards
    if (card.position > hovered.position && hoverClient > hoverMiddle) {
      return false;
    }

    return true;
  }

  moveCard(card, hovered) {
    const { cards = [], updateCardPosition } = this.props;

    const hoveredIndex = cards.findIndex((c) => c.id === hovered.id);

    if (hovered.position < card.position) {
      if (hoveredIndex === 0) {
        updateCardPosition({ id: card.id, position: hovered.position / 2 });
      } else {
        const previous = cards[hoveredIndex - 1];
        updateCardPosition({
          id: card.id,
          position: previous.position + ((hovered.position - previous.position) / 2)
        });
      }
    } else if (hovered.position > card.position) {
      if (hoveredIndex === cards.length - 1) {
        updateCardPosition({ id: card.id, position: hovered.position + (1 << 16) });
      } else {
        const next = cards[hoveredIndex + 1];
        updateCardPosition({
          id: card.id,
          position: hovered.position + ((next.position - hovered.position) / 2)
        });
      }
    }
  }

  toggleEditMode() {
    this.setState({
      isEditing: !this.state.isEditing
    });
  }

  render() {
    const { cards = [], list = {}, index, showModal, updateList, deleteList } = this.props;
    let input;

    const save = () => {
      if (!input.value.trim() || input.value === list.name) {
        this.toggleEditMode();
        return;
      }

      updateList({ id: list.id, name: input.value });
      this.toggleEditMode();
    };

    return (
      <div className="plp-cards-list-wrapper">
        <div className="plp-cards-list">
          { !this.state.isEditing ?
            <div className="plp-list-top" onClick={this.toggleEditMode}>{ list.name || `#${index + 1}` }</div> :
            <input
              type="text"
              className="plp-list-top-edit"
              defaultValue={list.name}
              autoFocus
              ref={node => { input = node; }}
              onBlur={save}
              onKeyPress={(e) => {
                if (e.key === 'Enter') {
                  save();
                }
              }}
            />
          }
          {cards.map((card) => (
            <Draggable
              key={card.id}
              id={card.id}
              name={card.name}
              onDrag={this.onDrag}
              onEndDrag={this.onEndDrag}
              userIsDragging={this.state.userIsDragging}
            />
        ))}
          <div className="plp-list-bottom">
            <Link to={`/b/${list.boardID}/${list.id}/newcard`}>
              <Button color="secondary" size="sm" className="btn-new-card">+ Add a new card</Button>
            </Link>
            <Button
              color="secondary"
              size="sm"
              className="btn-delete-list"
              onClick={() =>
                cards.length > 0 ?
                  showModal(ActionTypes.MODAL_DELETE_LIST, list) :
                  deleteList(list)
              }
            >Delete
          </Button>
          </div>
        </div>
      </div>
    );
  }
}

const source = {
  beginDrag(props) {
    return { id: props.id };
  },

  endDrag(props, monitor) {
    if (!monitor.didDrop()) {
      return;
    }

    const item = monitor.getItem();
    props.onEndDrag(item.id);
  }
};

const target = {
  hover(props, monitor, component) {
    const { onDrag, id } = props;
    const item = monitor.getItem();

    // Don't replace items with themselves
    if (item.id === id) {
      return;
    }

    onDrag(item.id, id, monitor, component);
  }
};

const Draggable = DropTarget('Card', target, connect => ({
  connectDropTarget: connect.dropTarget()
}))(DragSource('Card', source, (connect, monitor) => ({
  connectDragSource: connect.dragSource(),
  isDragging: monitor.isDragging()
}))(Compact));

export default DragDropContext(HTML5Backend)(List);
