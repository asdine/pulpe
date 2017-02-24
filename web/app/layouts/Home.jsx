import React from 'react';
import BoardDetail from '../containers/board';
import BoardList from '../containers/boardList';
import Modals from '../containers/modal';

const Home = ({ children }) =>
  <div className="wrapper container-fluid">
    <MobileMenu />
    <Menu />
    <BoardDetail />
    {children}
    <Modals />
  </div>;

const Menu = () =>
  <div className="plp-left-menu hidden-sm-down">
    <h1>pulpe</h1>
    <BoardList />
  </div>;

const MobileMenu = () =>
  <div className="plp-top-menu hidden-md-up">
    <h1>pulpe</h1>
  </div>;

export default Home;
