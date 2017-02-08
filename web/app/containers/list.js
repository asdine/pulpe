import { connect } from 'react-redux';
import List from '../components/List';
import { getListByID, getCardsByListID } from '../reducers';
import * as actions from '../actions';

const mapStateToProps = (state, { id }) => {
  const list = getListByID(state, id);
  return {
    list,
    cards: getCardsByListID(state, list.boardID, list.id)
  };
};

export default connect(
  mapStateToProps,
  actions
)(List);
