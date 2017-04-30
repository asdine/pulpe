import React from 'react';
import BoardList from './BoardList';

const Menu = () =>
  <div className="plp-left-menu" style={{ position: 'relative' }}>
    <h1>pulpe</h1>
    <BoardList />
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
