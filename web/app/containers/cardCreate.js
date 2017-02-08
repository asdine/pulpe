import React from 'react';
import { connect } from 'react-redux';
import { LargeCreate } from '../components/Card';
import { getCardsByListID } from '../reducers';
import * as actions from '../actions';

const mapStateToProps = (state, { params }) => ({
  card: {
    boardID: params.id,
    listID: params.listID,
  },
  cards: getCardsByListID(state, params.id, params.listID)
});

class CardCreate extends React.Component {
  componentDidMount() {
    const { fetchBoard, setActiveBoard, params } = this.props;
    fetchBoard(params.id);
    setActiveBoard(params.id);
  }

  componentDidUpdate(prevProps) {
    const { fetchBoard, setActiveBoard, params } = this.props;
    if (prevProps.params.id !== params.id) {
      fetchBoard(params.id);
      setActiveBoard(params.id);
    }
  }

  render() {
    return (
      <LargeCreate {...this.props} />
    );
  }
}

export default connect(
  mapStateToProps,
  actions
)(CardCreate);
