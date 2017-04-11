import { Component } from 'react';
import { connect } from 'react-redux';
import { showModal, hideModal } from '@/components/Modal/duck';
import { MODAL_CARD_DETAIL } from '@/Home/Board/List/Card/duck';

class CardRoute extends Component {
  componentDidMount() {
    this.props.showCardModal(this.props.params);
  }

  componentDidUpdate() {
    this.props.showCardModal(this.props.params);
  }

  componentWillUnmount() {
    this.props.hideModal();
  }

  render() {
    return null;
  }
}

export default connect(
  null,
  (dispatch) => ({
    showCardModal: (params) => {
      dispatch(showModal(MODAL_CARD_DETAIL, params));
    },
    hideModal: () => dispatch(hideModal())
  }),
)(CardRoute);
