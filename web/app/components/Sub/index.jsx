import React, { Component } from 'react';
import { connect } from 'react-redux';
import { registerSub, subIsOpened, addSub, closeAllSubs } from './duck';

const onClick = (e) => {
  e.stopPropagation();
};

class Sub extends Component {
  componentDidMount() {
    this.props.onEnter(this.props.name);
  }

  render() {
    const { children, onEnter, ...rest } = this.props;

    const childrenWithProps = React.Children.map(children,
      (child) => React.cloneElement(child, {
        ...child.props,
        ...rest,
      })
    );

    return (
      <div>
        {childrenWithProps}
      </div>
    );
  }
}

export default connect(
  (state, { name }) => ({
    isOpened: subIsOpened(state, name)
  }),
  (dispatch) => ({
    onEnter: (name) => dispatch(registerSub(name)),
    open: (name) => dispatch(addSub(name)),
    closeAll: () => dispatch(closeAllSubs())
  }),
)(Sub);

export const SubOpened = ({ isOpened, children }) => {
  if (!isOpened) {
    return null;
  }

  return (
    <div onClick={onClick}>
      { children }
    </div>
  );
};

export const SubClosed = (props) => {
  const { isOpened, children, ...rest } = props;

  if (isOpened) {
    return null;
  }

  const childrenWithProps = React.Children.map(children,
    (child) => React.cloneElement(child, {
      ...child.props,
      ...rest,
    })
  );

  return (
    <div>
      { childrenWithProps }
    </div>
  );
};


export const SubOpener = ({ name, open, closeAll, children }) =>
  <div
    onClick={(e) => {
      e.stopPropagation();
      closeAll();
      open(name);
    }}
  >{ children }</div>;
