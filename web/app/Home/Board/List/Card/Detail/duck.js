import { combineEpics } from 'redux-observable';
import { browserHistory } from 'react-router';
import { successOf } from '@/services/api/ajaxEpic';
import { getBoardSelector } from '@/Home/Board/duck';
import { getListSelector } from '@/Home/Board/List/duck';
import { UPDATE, DELETE, updateCard } from '@/Home/Board/List/Card/duck';

export const DOMAIN = 'pulpe/home/board/list/card/detail';

// types
export const SAVE = `${DOMAIN}/SAVE`;

// action creators
export const saveCard = (patch) => ({
  type: SAVE,
  patch
});

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

const saveCardEpic = action$ => action$.ofType(SAVE)
  .debounceTime(5000)
  .map(({ patch }) => updateCard(patch));

export const epics = combineEpics(/* eslint import/prefer-default-export: 0 */
  redirectCardModalOnNameUpdateSuccessEpic,
  closeCardModalOnDeleteSuccessEpic,
  saveCardEpic
);
