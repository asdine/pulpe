import { normalize } from 'normalizr';
import { Observable } from 'rxjs/Observable';
import { requestType, successType, failureType } from '../actions/types';

const ajaxEpic = (type, req, schema) =>
  action$ => action$.ofType(requestType(type))
    .mergeMap(req)
    .map(response => {
      if (schema) {
        return {
          type: successType(type),
          response: normalize(response, schema)
        };
      }

      return response;
    })
    .catch(error => Observable.of({
      type: failureType(type),
      payload: error.xhr.response,
      error
    }));

export default ajaxEpic;
