import React, { Component } from 'react';
import { connect } from 'react-redux';
import { setActiveBoard } from '@/Home/duck';

class BoardRoute extends Component {
  componentDidMount() {
    this.props.setActiveBoard(this.props.owner, this.props.boardSlug);
  }

  componentDidUpdate() {
    this.props.setActiveBoard(this.props.owner, this.props.boardSlug);
  }

  render() {
    return <div>{this.props.children}</div>;
  }
}

export default connect(
  (state, { params }) => ({
    boardSlug: params.board,
    owner: params.owner,
  }), {
    setActiveBoard,
  },
)(BoardRoute);
