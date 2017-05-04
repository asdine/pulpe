import { Observable } from 'rxjs/Observable';

class Client {
  constructor(url) {
    this.url = url;
  }

  register = (payload) => post(`${this.url}/register`, payload)

  login = (payload) => post(`${this.url}/login`, payload)

  getBoards = (filters = {}) => get(`${this.url}/user/boards${
    filters ?
      Object.keys(filters)
        .map((k, i) => `${i > 0 ? ',' : '?'}${k}=${filters[k]}`)
        .reduce((a, c) => a + c, '')
      : ''}`)

  getBoard = (owner, slug) => get(`${this.url}/boards/${owner}/${slug}`)

  createBoard = (payload) => post(`${this.url}/user/boards`, payload)

  deleteBoard = (id) => del(`${this.url}/boards/${id}`)

  updateBoard = ({ id, patch }) => update(`${this.url}/boards/${id}`, patch)

  createList = ({ boardID, type, ...rest }) => post(`${this.url}/boards/${boardID}/lists`, rest)

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


const client = new Client('/api');

export default client;
