import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { getFirstBoardSlug } from '../reducers';

const mapStateToProps = (state) => ({
  slug: getFirstBoardSlug(state)
});

class FirstBoard extends React.Component {
  componentWillReceiveProps(nextProps) {
    const { slug } = nextProps;
    if (slug) {
      browserHistory.push(`/${slug}`);
    }
  }

  render() {
    return null;
  }
}

export default connect(mapStateToProps)(FirstBoard);
