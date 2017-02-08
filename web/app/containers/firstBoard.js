import React from 'react';
import { connect } from 'react-redux';
import { hashHistory } from 'react-router';
import { getFirstBoardID } from '../reducers';

const mapStateToProps = (state) => ({
  id: getFirstBoardID(state)
});

class FirstBoard extends React.Component {
  componentWillReceiveProps(nextProps) {
    const { id } = nextProps;
    if (id) {
      hashHistory.push(`/b/${id}`);
    }
  }

  render() {
    return null;
  }
}

export default connect(mapStateToProps)(FirstBoard);
