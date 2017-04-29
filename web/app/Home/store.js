import 'rxjs/add/operator/mergeMap';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/mapTo';
import 'rxjs/add/operator/ignoreElements';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/do';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/delay';
import 'rxjs/add/observable/of';
import 'rxjs/add/observable/empty';
import 'rxjs/add/observable/dom/ajax';
import { combineReducers } from 'redux';
import { combineEpics } from 'redux-observable';
import ModalReducer from '@/components/Modal/duck';
import SubReducer from '@/components/Sub/duck';
import configureStore from '@/store';
import BoardListReducer, { epics as BoardListEpics } from './Menu/BoardList/duck';
import BoardReducer, { epics as BoardEpics } from './Board/duck';
import { epics as HeaderEpics } from './Board/Header/duck';
import ListReducer, { epics as ListEpics } from './Board/List/duck';
import CardReducer, { epics as CardEpics } from './Board/List/Card/duck';
import { epics as CardDetailEpics } from './Board/List/Card/Detail/duck';
import HomeReducer from './duck';

const reducers = combineReducers({
  ...HomeReducer,
  ...BoardListReducer,
  ...BoardReducer,
  ...ListReducer,
  ...CardReducer,
  ...ModalReducer,
  ...SubReducer,
});

const epics = combineEpics(
  BoardListEpics,
  BoardEpics,
  HeaderEpics,
  ListEpics,
  CardEpics,
  CardDetailEpics,
);

export default configureStore(reducers, epics);
