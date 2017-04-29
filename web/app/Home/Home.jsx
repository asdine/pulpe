import React from 'react';
import Board from './Board';
import Menu from './Menu';
import Modals from './Modals';

const Home = ({ children }) =>
  <div className="wrapper container-fluid">
    <Menu />
    <Board />
    {children}
    <Modals />
  </div>;

export default Home;
