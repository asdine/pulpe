import React from 'react';
import { connect } from 'react-redux';
import { getModalType, getModalProps, isEditing, getActiveBoard, getEditLevel } from '../reducers';
import * as actions from '../actions';
import { MainModal as BasicModal, CreateBoardModal, DeleteBoardModal, CreateListModal, DeleteListModal, DeleteCardModal } from '../components/Modal';
import * as ActionTypes from '../actions/types';

export const MainModal = connect(
  (state) => ({
    editLevel: getEditLevel(state),
    isEditing: isEditing(state, 'card-name') || isEditing(state, 'card-desc'),
    board: getActiveBoard(state)
  }),
  actions
)(BasicModal);

const CreateBoard = connect(
  (state) => ({
    isOpen: getModalType(state) === ActionTypes.MODAL_CREATE_BOARD
  }),
  actions,
)(CreateBoardModal);

const DeleteBoard = connect(
  (state) => ({
    id: getModalProps(state).id,
    isOpen: getModalType(state) === ActionTypes.MODAL_DELETE_BOARD
  }),
  actions,
)(DeleteBoardModal);

const CreateList = connect(
  (state) => ({
    board: getModalProps(state),
    isOpen: getModalType(state) === ActionTypes.MODAL_CREATE_LIST
  }),
  actions,
)(CreateListModal);

const DeleteList = connect(
  (state) => ({
    list: getModalProps(state),
    isOpen: getModalType(state) === ActionTypes.MODAL_DELETE_LIST
  }),
  actions,
)(DeleteListModal);

const DeleteCard = connect(
  (state) => ({
    card: getModalProps(state),
    isOpen: getModalType(state) === ActionTypes.MODAL_DELETE_CARD
  }),
  actions,
)(DeleteCardModal);

const Modals = () =>
  <div>
    <CreateBoard />
    <DeleteBoard />
    <CreateList />
    <DeleteList />
    <DeleteCard />
  </div>;

export default Modals;
