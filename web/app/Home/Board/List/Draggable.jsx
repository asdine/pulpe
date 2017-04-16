import React, { Component } from 'react';
import { DragSource, DropTarget, DragLayer } from 'react-dnd';
import { getEmptyImage } from 'react-dnd-html5-backend';

const itemSource = {
  beginDrag(props) {
    const { id, index, list } = props;
    return { id, index, list };
  },

  isDragging(props, monitor) {
    return monitor.getItem().id === props.id;
  },
};

const itemTarget = {
  hover(props, monitor) {
    const { id: draggedId } = monitor.getItem();
    const { id: overId, list } = props;

    if (draggedId !== overId) {
      const { index: overIndex } = props.findItem(overId);
      props.moveItem(draggedId, list.id, overIndex);
    }
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
  card: monitor.getItem(),
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
    if (nextProps.items !== this.props.items || nextProps.card !== this.props.card) {
      this.setPreviewFromProps(nextProps);
    }
  }

  setPreviewFromProps(props) {
    const { items, card = {} } = props;

    if (!items || !card) {
      return;
    }

    const preview = items.find(item => item.props.id === card.id);

    if (!preview) {
      return;
    }

    this.setState({
      preview
    });
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
