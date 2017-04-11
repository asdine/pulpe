import React, { Component } from 'react';
import { DragSource, DropTarget } from 'react-dnd';

const itemSource = {
  beginDrag({ id, index }) {
    return { id, index };
  },

  endDrag(props, monitor) {
    const { id: droppedId, index } = monitor.getItem();
    const didDrop = monitor.didDrop();

    if (!didDrop) {
      props.moveItem(droppedId, index);
    }
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
  }
};


@DropTarget('Item', itemTarget, conn => ({
  connectDropTarget: conn.dropTarget(),
}))
@DragSource('Item', itemSource, (conn, monitor) => ({
  connectDragSource: conn.dragSource(),
  connectDragPreview: conn.dragPreview(),
  isDragging: monitor.isDragging()
}))
class Draggable extends Component {
  render() {
    const { connectDragSource, connectDragPreview, connectDropTarget,
            children, className } = this.props;

    const childrenWithProps = React.Children.map(children,
     (child) => React.cloneElement(child, {
       connectDragSource
     })
    );

    return connectDropTarget(connectDragPreview(
      <div className={className}>{childrenWithProps}</div>
    ));
  }
}

export default Draggable;
