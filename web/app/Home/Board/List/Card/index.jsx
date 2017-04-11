import React from 'react';
import { browserHistory } from 'react-router';

const Card = ({ card = {}, list = {}, board = {} }) =>
  <div className="card" onClick={() => browserHistory.push(`/${board.slug}/${list.slug}/${card.slug}`)}>
    <div className="card-block">
      <h3 className="card-title">{ card.name }</h3>
    </div>
  </div>;

export default Card;
