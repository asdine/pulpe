import React, { Component } from 'react';
import { Modal as BootstrapModal } from 'reactstrap';
import './style.scss';

class Modal extends Component {
  render() {
    const { isOpen, toggle, children } = this.props;

    return (
      <BootstrapModal
        isOpen={isOpen}
        toggle={toggle}
        backdrop="static"
        modalClassName="plp-modal"
        onClick={toggle}
      >
        <ModalContent>
          {children}
        </ModalContent>
      </BootstrapModal>
    );
  }
}

const modalContentClick = (e) => {
  e.stopPropagation();
};

const ModalContent = ({ children }) =>
  <div onClick={modalContentClick}>
    {children}
  </div>;

export default Modal;
