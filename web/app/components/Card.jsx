import React from 'react';
import { hashHistory } from 'react-router';
import { Button, ModalBody } from 'reactstrap';
import * as ActionTypes from '../actions/types';

export const Compact = ({
  connectDragSource, connectDropTarget,
  id, name, isDragging, userIsDragging }) =>
  connectDragSource(connectDropTarget(
    <div
      className={`card ${userIsDragging ? '' : 'card-hover'}`}
      style={{ opacity: isDragging ? 0 : 1 }}
      onClick={() => hashHistory.push(`/c/${id}`)}
    >
      <div className="card-block">
        <h3 className="card-title">{ name }</h3>
      </div>
    </div>
  ));

export const Large = (props) => {
  const { card, toggle, disableAllEditModes, showModal, setEditLevel } = props;

  if (card === undefined) {
    return null;
  }

  return (
    <div
      className="plp-modal-large-card"
      onClick={() => {
        disableAllEditModes();
        setEditLevel(0);
      }}
      onFocus={(e) => {
        if (['submit', 'button'].indexOf(e.target.type) > -1) {
          return;
        }
        setEditLevel(1);
      }}
    >
      <div className="modal-header">
        <Name {...props} />
        <div className="modal-options clearfix">
          <Button className="close" data-dismiss="modal" aria-label="Close" onClick={toggle}>
            <span aria-hidden="true">&times;</span>
          </Button>
          <Button
            color="danger"
            size="sm"
            className="float-xs-right"
            onClick={() => {
              setEditLevel(0);
              showModal(ActionTypes.MODAL_DELETE_CARD, card);
            }}
          >
            Delete
          </Button>
        </div>
      </div>
      <ModalBody>
        <Description {...props} />
      </ModalBody>
    </div>
  );
};

const Name = ({ card, isEditingName, saveName,
                setEditLevel, disableAllEditModes, toggleEditMode }) => {
  let input;

  return (
    <div className="modal-title">
      { !isEditingName ?
        <h5
          onClick={(e) => {
            e.stopPropagation();
            disableAllEditModes();
            toggleEditMode('card-name');
          }}
        >{ card.name }</h5> :
        <input
          type="text"
          autoFocus
          className="modal-title-edit"
          ref={node => { input = node; }}
          defaultValue={card.name}
          onBlur={() => saveName(input)}
          onClick={e => e.stopPropagation()}
          onKeyPress={(e) => {
            if (e.key === 'Enter') {
              saveName(input);
              setEditLevel(0);
            }
          }}
        />
      }
    </div>
  );
};

const Description = ({ card, toggleEditMode, disableAllEditModes, saveDesc, isEditingDesc }) => {
  let textarea;

  return (
    <div>
      { !isEditingDesc &&
        <div
          className="large-card-description"
          onClick={(e) => {
            e.stopPropagation();
            toggleEditMode('card-desc');
          }}
        >
          {card.description || <div className="large-card-description__no-description">Click here to add content</div>}
        </div>
      }
      { isEditingDesc &&
        <div className="large-card-description-edit">
          <textarea
            name="large-card-description-edit"
            defaultValue={card.description}
            autoFocus
            onClick={e => e.stopPropagation()}
            ref={node => {
              textarea = node;
            }}
          />
          <div className="large-card-description-edit__footer">
            <button type="button" className="btn btn-secondary cancel-btn" onClick={disableAllEditModes}>Cancel</button>
            <button type="button" className="btn btn-primary save-btn" onClick={() => saveDesc(textarea)}>Save</button>
          </div>
        </div>
      }
    </div>
  );
};

export const LargeCreate = ({ card, cards, createCard, ...rest }) => (
  <LargeForm
    card={card}
    {...rest}
    onSave={(input, textarea) => {
      const name = input.value.trim();
      const description = textarea.value.trim();

      if (!name) {
        return;
      }

      const newCard = {
        boardID: card.boardID,
        listID: card.listID,
        name,
        description
      };

      newCard.position = cards.length > 0 ?
          cards[cards.length - 1].position + (1 << 16) :
          1 << 16;

      createCard(newCard);

      return hashHistory.push(`/b/${card.boardID}`);
    }}
  />
  );

const LargeForm = ({ card, toggle, onSave }) => {
  let input;
  let textarea;

  if (card === undefined) {
    return null;
  }

  return (
    <div>
      <div className="modal-header clearfix">
        <div className="row">
          <div className="col-xs-8">
            <div className="form-group">
              <input
                type="text"
                className="form-control"
                placeholder="Card name"
                defaultValue={card.name}
                ref={node => {
                  input = node;
                  return input && input.focus();
                }}
              />
            </div>
          </div>
          <div className="col-xs-3 offset-xs-1 clearfix">
            <button type="button" className="close" data-dismiss="modal" aria-label="Close" onClick={toggle}>
              <span aria-hidden="true">&times;</span>
            </button>
            <button
              type="button"
              className="btn btn-secondary btn-sm float-xs-right"
              onClick={(e) => {
                e.preventDefault();
                onSave(input, textarea);
              }}
            >
              Save
            </button>
          </div>
        </div>
      </div>
      <ModalBody>
        <div className="form-group">
          <label htmlFor="card-content">Content</label>
          <textarea
            className="form-control"
            id="card-content"
            rows="3"
            defaultValue={card.description}
            ref={node => {
              textarea = node;
            }}
          />
        </div>
      </ModalBody>
    </div>
  );
};
