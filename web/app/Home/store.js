import { combineEpics } from 'redux-observable';
import BoardListReducer, { epics as BoardListEpics } from '@/Home/Menu/BoardList/duck';
import BoardReducer, { epics as BoardEpics } from '@/Home/Board/duck';
import { epics as HeaderEpics } from '@/Home/Board/Header/duck';
import ListReducer, { epics as ListEpics } from '@/Home/Board/List/duck';
import CardReducer, { epics as CardEpics } from '@/Home/Board/List/Card/duck';
import { epics as CardDetailEpics } from '@/Home/Board/List/Card/Detail/duck';
import HomeReducer from './duck';

export const reducers = {
  ...HomeReducer,
  ...BoardListReducer,
  ...BoardReducer,
  ...ListReducer,
  ...CardReducer,
};

export const epics = combineEpics(
  BoardListEpics,
  BoardEpics,
  HeaderEpics,
  ListEpics,
  CardEpics,
  CardDetailEpics,
);
