import React, { Component } from 'react';
import { connect } from 'react-redux';
import { setActiveBoardID } from '@/Home/duck';

class BoardRoute extends Component {
  componentDidMount() {
    this.props.setActiveBoardID(this.props.boardSlug);
  }

  componentDidUpdate() {
    this.props.setActiveBoardID(this.props.boardSlug);
  }

  render() {
    return <div>{this.props.children}</div>;
  }
}

export default connect(
  (state, { params }) => ({
    boardSlug: params.board,
  }), {
    setActiveBoardID,
  },
)(BoardRoute);
