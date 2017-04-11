import { Observable } from 'rxjs/Observable';

class Client {
  constructor(url) {
    this.url = url;
  }

  getBoards = (filters = {}) => get(`${this.url}/boards${
    filters ?
      Object.keys(filters)
        .map((k, i) => `${i > 0 ? ',' : '?'}${k}=${filters[k]}`)
        .reduce((a, c) => a + c, '')
      : ''}`)

  getBoard = (id) => get(`${this.url}/boards/${id}`)

  createBoard = ({ name }) => post(`${this.url}/boards`, {
    name
  })

  deleteBoard = (id) => del(`${this.url}/boards/${id}`)

  updateBoard = ({ id, patch }) => update(`${this.url}/boards/${id}`, patch)

  createList = ({ boardID, name }) => post(`${this.url}/boards/${boardID}/lists`, {
    name
  })

  updateList = ({ id, patch }) => update(`${this.url}/lists/${id}`, patch)

  deleteList = (id) => del(`${this.url}/lists/${id}`)

  getCard = (id) => get(`${this.url}/cards/${id}`)

  createCard = ({ id, listID, name, description, position }) => post(`${this.url}/lists/${listID}/cards`, {
    id,
    name,
    description,
    position
  })

  deleteCard = (id) => del(`${this.url}/cards/${id}`)

  updateCard = ({ id, patch }) => update(`${this.url}/cards/${id}`, patch)
}

const get = (url) => Observable.ajax.getJSON(url);

const post = (url, data) => Observable.ajax({
  url,
  method: 'POST',
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json'
  },
  responseType: 'json',
  body: JSON.stringify(data)
}).map(response => response.response);

const del = (url) => Observable.ajax({
  url,
  method: 'DELETE'
});

const update = (url, data) => Observable.ajax({
  url,
  method: 'PATCH',
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json'
  },
  responseType: 'json',
  body: JSON.stringify(data)
}).map(response => response.response);

const client = new Client('/v1');

export default client;
