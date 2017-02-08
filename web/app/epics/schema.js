import { schema } from 'normalizr';

export const board = new schema.Entity('boards');
export const list = new schema.Entity('lists');
export const card = new schema.Entity('cards');

board.define({
  lists: [list],
  cards: [card]
});
