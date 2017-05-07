import React from 'react';
import DeleteBoardModal from '@/Home/Board/Header/DeleteBoardModal';
import DeleteListModal from '@/Home/Board/List/DeleteListModal';
import CardDetailModal from '@/Home/Board/List/Card/Detail';

const Modals = () =>
  <div>
    <DeleteBoardModal />
    <DeleteListModal />
    <CardDetailModal />
  </div>;

export default Modals;
