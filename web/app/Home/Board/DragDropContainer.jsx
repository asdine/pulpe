import React, { Component } from 'react';
import update from 'react/lib/update';
import Draggable from './Draggable';

class DragDropContainer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      items: [...this.props.children]
    };

    this.findItem = this.findItem.bind(this);
    this.moveItem = this.moveItem.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      items: [...nextProps.children]
    });
  }

  findItem(id) {
    const index = this.state.items.findIndex(child => child.props.id === id);
    return {
      item: this.state.items[index],
      index
    };
  }

  moveItem(id, atIndex) {
    const { item, index } = this.findItem(id);
    this.setState(update(this.state, {
      items: {
        $splice: [
          [index, 1],
          [atIndex, 0, item],
        ],
      },
    }));
  }

  render() {
    const { className, itemClassName, onDrop } = this.props;
    const { items } = this.state;

    return (
      <div className={className}>
        {items.map((child, i) => (
          <Draggable
            className={itemClassName}
            key={items[i].props.id}
            id={items[i].props.id}
            index={i}
            findItem={this.findItem}
            moveItem={this.moveItem}
            onDrop={onDrop}
          >
            {child}
          </Draggable>
        ))}
      </div>
    );
  }
}

export default DragDropContainer;
