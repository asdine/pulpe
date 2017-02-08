import { connect } from 'react-redux';
import { getActiveBoard } from '../reducers';
import Modal from '../components/Modal';

const mapStateToProps = (state) => ({
  board: getActiveBoard(state)
});

export default connect(mapStateToProps)(Modal);
