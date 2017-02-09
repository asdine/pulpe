import { Observable } from 'rxjs/Observable';

class Client {
  constructor(url) {
    this.url = url;
  }

  allBoards = () => get(`${this.url}/boards`)

  getBoard = (id) => get(`${this.url}/boards/${id}`)

  createBoard = ({ name }) => post(`${this.url}/boards`, {
    name
  })

  deleteBoard = (id) => del(`${this.url}/boards/${id}`)

  updateBoard = ({ id, toUpdate }) => update(`${this.url}/boards/${id}`, toUpdate)

  createList = ({ boardID, name }) => post(`${this.url}/lists`, {
    boardID,
    name
  })

  updateList = ({ id, toUpdate }) => update(`${this.url}/lists/${id}`, toUpdate)

  deleteList = (id) => del(`${this.url}/lists/${id}`)

  getCard = (id) => get(`${this.url}/cards/${id}`)

  createCard = ({ id, boardID, listID, name, description, position }) => post(`${this.url}/cards`, {
    id,
    boardID,
    listID,
    name,
    description,
    position
  })

  deleteCard = (id) => del(`${this.url}/cards/${id}`)

  updateCard({ id, name, description, position }) {
    const patch = {};
    if (name !== undefined) {
      patch.name = name;
    }

    if (description !== undefined) {
      patch.description = description;
    }

    if (position !== undefined) {
      patch.position = position;
    }

    return update(`${this.url}/cards/${id}`, patch);
  }
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
