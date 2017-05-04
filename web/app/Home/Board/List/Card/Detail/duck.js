import { combineEpics } from 'redux-observable';
import { browserHistory } from 'react-router';
import { successOf } from '@/services/api/ajaxEpic';
import { getBoardSelector } from '@/Home/Board/duck';
import { getListSelector } from '@/Home/Board/List/duck';
import { UPDATE, DELETE } from '@/Home/Board/List/Card/duck';

const redirectCardModalOnNameUpdateSuccessEpic = (action$, store) => action$.ofType(successOf(UPDATE))
  .do((action) => {
    if (!action.originalAction.patch.name) {
      return;
    }

    const card = action.response.entities.cards[action.response.result];
    const board = getBoardSelector(store.getState());
    const list = getListSelector(store.getState(), card.listID);
    browserHistory.push(`/${board.owner.login}/${board.slug}/${list.slug}/${card.slug}`);
  })
  .ignoreElements();

const closeCardModalOnDeleteSuccessEpic = (action$, store) => action$.ofType(successOf(DELETE))
  .do(() => {
    const board = getBoardSelector(store.getState());
    browserHistory.push(`/${board.owner.login}/${board.slug}`);
  })
  .ignoreElements();

export const epics = combineEpics(/* eslint import/prefer-default-export: 0 */
  redirectCardModalOnNameUpdateSuccessEpic,
  closeCardModalOnDeleteSuccessEpic,
);
