import React from 'react';
import { connect } from 'react-redux';
import { getModalType, getModalProps, isEditing, getActiveBoard, getEditLevel, getBoards, getCardsByListID, getBoardByID } from '../reducers';
import * as actions from '../actions';
import { MainModal as BasicModal, CreateBoardModal, DeleteBoardModal, CreateListModal, DeleteListModal, DeleteCardModal, CreateCardModal } from '../components/Modal';
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
  (state) => {
    const id = getModalProps(state).id;
    const nextBoard = getBoards(state).find(b => b.id !== id) || {};
    return ({
      id,
      isOpen: getModalType(state) === ActionTypes.MODAL_DELETE_BOARD,
      redirectTo: nextBoard.slug
    });
  },
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

const CreateCard = connect(
  (state) => {
    const list = getModalProps(state);
    return ({
      list,
      board: getBoardByID(state, list.boardID),
      cards: getCardsByListID(state, list.boardID, list.id),
      isOpen: getModalType(state) === ActionTypes.MODAL_CREATE_CARD
    });
  },
  actions,
)(CreateCardModal);

const DeleteCard = connect(
  (state) => {
    const card = getModalProps(state);
    const board = getBoardByID(state, card.boardID) || {};
    return ({
      card,
      isOpen: getModalType(state) === ActionTypes.MODAL_DELETE_CARD,
      redirectTo: board.slug
    });
  },
  actions,
)(DeleteCardModal);

const Modals = () =>
  <div>
    <CreateBoard />
    <DeleteBoard />
    <CreateList />
    <DeleteList />
    <CreateCard />
    <DeleteCard />
  </div>;

export default Modals;
