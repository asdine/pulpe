import { successType, FETCH_BOARD } from '../actions/types';
import { updateCard } from '../actions';

const fixCardPositionEpic = action$ => action$.ofType(successType(FETCH_BOARD))
  .filter(action => action.response.entities.cards)
  .mergeMap(action =>
    Object.keys(action.response.entities.cards)
      .filter(id => !action.response.entities.cards[id].position))
  .map((id, i) => updateCard({
    id,
    position: (i + 1) * (1 << 16)
  }));

export default fixCardPositionEpic;
