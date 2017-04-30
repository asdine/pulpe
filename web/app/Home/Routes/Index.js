import { Component } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { getBoards } from '@/Home/Menu/BoardList/duck';

class BoardIndex extends Component {
  componentDidMount() {
    this.redirectToFirstBoard();
  }

  componentDidUpdate() {
    this.redirectToFirstBoard();
  }

  redirectToFirstBoard() {
    const { boards = [] } = this.props;

    if (boards.length > 0) {
      browserHistory.push(`/${boards[0].owner.login}/${boards[0].slug}`);
    }
  }

  render() {
    return null;
  }
}

export default connect(
  (state) => ({
    boards: getBoards(state),
  })
)(BoardIndex);
