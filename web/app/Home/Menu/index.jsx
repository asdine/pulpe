import Logo from '@/assets/img/logo-regular.svg';
import React from 'react';
import BoardList from './BoardList';

const Menu = () =>
  <div className="left-menu" style={{ position: 'relative' }}>
    <img src={Logo} className="left-menu__logo" alt="Pulpe logo" height="24px" width="30px" />
    <nav>
      <div className="left-menu__board-panel">
        <header className="left-menu__title">My boards</header>
        <BoardList />
      </div>
    </nav>
    <div
      className="card"
      style={{
        position: 'absolute',
        bottom: 0
      }}
    >
      <div className="card-block">
        <a className="btn btn-secondary" href="/logout" role="button">Sign out</a>
      </div>
    </div>
  </div>;

export default Menu;
