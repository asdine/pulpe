import React from 'react';
import BoardDetail from '../containers/board';
import BoardList from '../containers/boardList';
import Modals from '../containers/modal';

const Home = ({ children }) =>
  <div className="wrapper container-fluid">
    <Menu />
    <BoardDetail />
    {children}
    <Modals />
  </div>;

const Menu = () =>
  <div className="plp-left-menu">
    <h1>pulpe</h1>
    <BoardList />
  </div>;

export default Home;
