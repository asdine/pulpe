import { Component } from 'react';
import { connect } from 'react-redux';
import * as actions from '../actions';
import { getBoardBySlug } from '../reducers';

class BoardActivate extends Component {
  componentDidMount() {
    const { slug, board, setActiveBoard, fetchBoards } = this.props;
    if (!slug) {
      return;
    }

    if (board && board.id) {
      setActiveBoard(board.id);
      return;
    }

    fetchBoards({ slug });
  }

  componentDidUpdate(prevProps) {
    const { slug, board, setActiveBoard, fetchBoards } = this.props;
    if (!slug) {
      return;
    }

    if (board && board.id) {
      setActiveBoard(board.id);
      return;
    }

    if (this.props.slug && this.props.slug === prevProps.slug) {
      return;
    }
    fetchBoards({ slug });
  }

  render() {
    return null;
  }
}

export default connect(
  (state, { params }) => ({
    slug: params.slug,
    board: getBoardBySlug(state, params.slug)
  }),
  actions,
)(BoardActivate);
