import React from 'react';
import { browserHistory } from 'react-router';

const Card = (props) => {
  const { card = {}, list = {}, board = {}, isDragging, isDragged, style } = props;

  return (
    <div style={style} className={`card-item ${isDragged ? 'dragged' : ''} ${isDragging ? 'shadow' : ''}`} onClick={() => browserHistory.push(`/${board.owner.login}/${board.slug}/${list.slug}/${card.slug}`)}>
      <div className="card-wrapper">
        <h3 className="card-title">{ card.name }</h3>
      </div>
    </div>
  );
};

export default Card;
