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
      <div className="left-menu__account-panel">
        <header className="left-menu__title">My account</header>
        <ul className="left-menu__account-list">
          <li className="left-menu__item">
            <a className="account-list__sign-out" href="/logout" role="button">Sign out</a>
          </li>
        </ul>
      </div>
    </nav>
  </div>;

export default Menu;
