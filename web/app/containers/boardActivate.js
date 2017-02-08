import { Component } from 'react';
import { connect } from 'react-redux';
import * as actions from '../actions';

class BoardActivate extends Component {
  componentDidMount() {
    const { setActiveBoard, params } = this.props;
    if (!params.id) {
      return;
    }
    setActiveBoard(params.id);
  }

  componentDidUpdate(prevProps) {
    const { setActiveBoard, params } = this.props;
    if (params.id !== prevProps.id) {
      setActiveBoard(params.id);
    }
  }

  render() {
    return null;
  }
}

export default connect(
  null,
  {
    setActiveBoard: actions.setActiveBoard
  }
)(BoardActivate);
