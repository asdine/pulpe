import React, { Component } from 'react';
import update from 'react/lib/update';

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
    const { className, onDrop, itemClassName } = this.props;
    const { items } = this.state;

    const childrenWithProps = React.Children.map(items,
     (child, i) => React.cloneElement(child, {
       className: itemClassName,
       key: items[i].props.id,
       id: items[i].props.id,
       index: i,
       findItem: this.findItem,
       moveItem: this.moveItem,
       onDrop
     })
    );

    return (
      <div className={className}>
        {childrenWithProps}
      </div>
    );
  }
}

export default DragDropContainer;
