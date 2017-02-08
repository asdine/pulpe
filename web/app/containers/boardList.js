import { connect } from 'react-redux';
import BoardList from '../components/BoardList';
import { getBoards, getActiveBoard } from '../reducers';
import * as actions from '../actions';

const mapStateToProps = (state) => ({
  boards: getBoards(state),
  activeBoard: getActiveBoard(state)
});

export default connect(
  mapStateToProps,
  actions,
)(BoardList);
