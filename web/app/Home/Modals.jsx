import React from 'react';
import CreateBoardModal from '@/Home/Menu/BoardList/CreateBoardModal';
import DeleteBoardModal from '@/Home/Board/Header/DeleteBoardModal';
import DeleteListModal from '@/Home/Board/List/DeleteListModal';
import CreateCardModal from '@/Home/Board/List/Card/CreateCardModal';
import CardDetailModal from '@/Home/Board/List/Card/Detail';

const Modals = () =>
  <div>
    <CreateBoardModal />
    <DeleteBoardModal />
    <DeleteListModal />
    <CreateCardModal />
    <CardDetailModal />
  </div>;

export default Modals;
