import React, { Component } from 'react';
import { DragSource, DropTarget, DragLayer } from 'react-dnd';
import { getEmptyImage } from 'react-dnd-html5-backend';
import List from './index';
import Card from './Card';

const itemSource = {
  beginDrag(props) {
    const { id, index, list, card } = props;
    return { id, index, list, card };
  },

  canDrag(props = {}) {
    return !props.locked;
  },

  isDragging(props, monitor) {
    return monitor.getItem().id === props.id;
  },

  endDrag(props, monitor) {
    if (!monitor.didDrop()) {
      props.onDrop(props.card, props.index, true);
      return;
    }
    const { card, index } = monitor.getDropResult();
    if (!card) {
      return;
    }

    props.onDrop(props.card, index);
  }
};

const itemTarget = {
  hover(props, monitor) {
    const { id: draggedId, list: fromList } = monitor.getItem();
    const { id: overId, list, cards = [] } = props;

    if (props.locked) {
      if (cards.length === 0 || fromList.id !== list.id) {
        const draggedCard = cards.find((c) => c.id === draggedId);
        if (!draggedCard) {
          props.moveToList({ id: draggedId, listID: list.id });
        }
      }
      return;
    }

    if (draggedId !== overId) {
      const { index: overIndex } = props.findItem(overId);
      props.moveItem(draggedId, list.id, overIndex);
    }
  },

  drop(props, monitor) {
    if (props.locked) {
      const { list, card } = monitor.getItem();
      return { list, card, index: props.cards.length - 1 };
    }

    const { list = {}, card = {}, index } = props;
    return { list, card, index };
  }
};


@DropTarget('Card', itemTarget, conn => ({
  connectDropTarget: conn.dropTarget(),
}))
@DragSource('Card', itemSource, (conn, monitor) => ({
  connectDragSource: conn.dragSource(),
  connectDragPreview: conn.dragPreview(),
  isDragging: monitor.isDragging()
}))
class Draggable extends Component {
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
            children, className, isDragging } = this.props;

    const childrenWithProps = React.Children.map(children,
      (child) => React.cloneElement(child, {
        connectDragSource,
        isDragging
      })
    );

    return connectDragSource(connectDropTarget(
      <div className={className}>{childrenWithProps}</div>
    ));
  }
}

export default Draggable;

function getItemStyles(props) {
  const { initialOffset, currentOffset } = props;
  if (!initialOffset || !currentOffset) {
    return {
      display: 'none',
    };
  }

  const { x, y } = currentOffset;

  const transform = `translate(${x}px, ${y}px) rotate(6deg)`;
  return {
    transform,
    WebkitTransform: transform,
  };
}

@DragLayer(monitor => ({
  item: monitor.getItem(),
  itemType: monitor.getItemType(),
  initialOffset: monitor.getInitialSourceClientOffset(),
  currentOffset: monitor.getSourceClientOffset(),
  isDragging: monitor.isDragging(),
}))
export class DraggablePreview extends Component { // eslint-disable-line react/no-multi-comp
  componentWillMount() {
    this.setState({});
    this.setPreviewFromProps(this.props);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.item !== this.props.item) {
      this.setPreviewFromProps(nextProps);
    }
  }

  setPreviewFromProps(props) {
    const { item, lists = [], board = {}, itemType } = props;

    if (!item) {
      return;
    }

    switch (itemType) {
      case 'Card': {
        const card = item.card;

        this.setState({
          preview: <Card id={card.id} card={card} board={board} list={lists.find(l => l.id === card.listID)} />
        });
        break;
      }
      case 'List': {
        this.setState({
          preview: <List id={item.id} />
        });
        break;
      }
      default: {
        this.setState({});
      }
    }
  }

  render() {
    const { isDragging, itemType } = this.props;
    const { preview } = this.state;

    if (!isDragging || !preview || itemType !== 'Card') {
      return null;
    }

    return (React.cloneElement(preview, {
      isDragged: isDragging,
      style: getItemStyles(this.props)
    }));
  }
}
