import React from 'react';

const BLOCK_TAGS = {
  blockquote: 'quote',
  p: 'paragraph',
  pre: 'code'
};

// Add a dictionary of mark tags.
const MARK_TAGS = {
  em: 'italic',
  strong: 'bold',
  u: 'underlined',
  code: 'code'
};

const rules = [
  {
    deserialize(el, next) {
      const type = BLOCK_TAGS[el.tagName];
      if (!type) return;
      return {
        kind: 'block',
        type,
        nodes: next(el.children)
      };
    },
    serialize(object, children) {
      if (object.kind !== 'block') return;
      switch (object.type) { // eslint-disable-line
        case 'code': return <pre><code>{children}</code></pre>;
        case 'paragraph': return <p>{children}</p>;
        case 'quote': return <blockquote>{children}</blockquote>;
      }
    }
  },
  // Add a new rule that handles marks...
  {
    deserialize(el, next) {
      const type = MARK_TAGS[el.tagName];
      if (!type) return;
      return {
        kind: 'mark',
        type,
        nodes: next(el.children)
      };
    },
    serialize(object, children) {
      if (object.kind !== 'mark') return;
      switch (object.type) { // eslint-disable-line
        case 'bold': return <strong>{children}</strong>;
        case 'italic': return <em>{children}</em>;
        case 'underlined': return <u>{children}</u>;
        case 'code': return <code>{children}</code>;
      }
    }
  }
];

export default rules;
