import React from 'react';
import { connect } from 'react-redux';
import { getActiveBoard, getDevStore } from '../reducers';
import * as actions from '../actions';

const mapStateToProps = (state) => ({
  store: getDevStore(state),
  board: getActiveBoard(state)
});

let Dev = (props) => {
  const { children, inline, ...rest } = props;
  const childrenWithProps = React.Children.map(children,
     (child) => React.cloneElement(child, {
       ...child.props,
       ...rest,
     })
  );

  return inline ?
    <span>{ childrenWithProps }</span> :
    <div>{ childrenWithProps }</div>;
};

Dev = connect(
  mapStateToProps,
  actions
)(Dev);

export default Dev;
