import { normalize } from 'normalizr';
import { Observable } from 'rxjs/Observable';

export const requestOf = (type) => `${type}_REQUEST`;
export const successOf = (type) => `${type}_SUCCESS`;
export const failureOf = (type) => `${type}_FAILURE`;

const ajaxEpic = (type, req, schema) =>
  action$ => action$.ofType(requestOf(type))
    .mergeMap(action => req(action)
      .map(response => {
        if (schema) {
          return {
            type: successOf(type),
            response: normalize(response, schema),
            originalAction: action
          };
        }

        return {
          type: successOf(type),
          response,
          originalAction: action
        };
      })
      .catch(error => Observable.of({
        type: failureOf(type),
        payload: error.xhr.response,
        error
      })));

export default ajaxEpic;
