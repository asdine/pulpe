import React from 'react';
import { Editor, Html } from 'slate';
import Portal from 'react-portal';
import rules from './rules';

const schema = {
  nodes: {
    code: props => <pre {...props.attributes}>{props.children}</pre>,
    paragraph: props => <p {...props.attributes}>{props.children}</p>,
    quote: props => <blockquote {...props.attributes}>{props.children}</blockquote>,
  },
  marks: {
    bold: props => <strong>{props.children}</strong>,
    code: props => <code>{props.children}</code>,
    italic: props => <em>{props.children}</em>,
    underlined: props => <u>{props.children}</u>,
  }
};

const html = new Html({ rules });

class RichEditor extends React.Component {
  constructor(props) {
    super(props);

    let content = '<p></p>';

    if (props.content.length > 0) {
      content = props.content[0] !== '<' ? `<p>${props.content}</p>` : props.content;
    }
    this.state = {
      state: html.deserialize(content)
    };

    this.renderMarkButton = this.renderMarkButton.bind(this);
    this.onClickMark = this.onClickMark.bind(this);
    this.onChange = this.onChange.bind(this);
    this.onOpen = this.onOpen.bind(this);
  }

  componentDidMount() {
    this.updateMenu();
  }

  componentDidUpdate() {
    this.updateMenu();
  }

  onChange(state) {
    this.props.onSave(html.serialize(state));
    this.setState({ state });
  }

  onClickMark(e, type) {
    e.preventDefault();
    let { state } = this.state;

    state = state
      .transform()
      .toggleMark(type)
      .apply();

    this.setState({ state });
  }

  onOpen(portal) {
    this.setState({ menu: portal.firstChild });
  }

  hasMark(type) {
    const { state } = this.state;
    return state.marks.some(mark => mark.type === type);
  }

  updateMenu() {
    const { menu, state } = this.state;
    if (!menu) return;

    if (state.isBlurred || state.isCollapsed) {
      menu.removeAttribute('style');
      return;
    }

    if (menu.hasAttribute('style')) {
      // ignore changes
      return;
    }

    const selection = window.getSelection();
    const range = selection.getRangeAt(0);
    const rect = range.getBoundingClientRect();
    menu.style.opacity = 1;
    menu.style.position = 'absolute';
    menu.style.display = 'block';
    menu.style.top = `${rect.top + window.scrollY - menu.offsetHeight}px`;
    menu.style.left = `${rect.left + window.scrollX - menu.offsetWidth / 2 + rect.width / 2}px`;
    menu.style.zIndex = 1060;
  }

  renderMenu() {
    return (
      <Portal isOpened onOpen={this.onOpen}>
        <div className="hover-menu bg-inverse">
          {this.renderMarkButton('bold', 'format_bold')}
          {this.renderMarkButton('italic', 'format_italic')}
          {this.renderMarkButton('underlined', 'format_underlined')}
          {this.renderMarkButton('code', 'code')}
        </div>
      </Portal>
    );
  }

  renderMarkButton(type, icon) {
    const isActive = this.hasMark(type);
    const onMouseDown = e => this.onClickMark(e, type);

    return (
      <span className="button" onMouseDown={onMouseDown} data-active={isActive}>
        <span className="material-icons">{icon}</span>
      </span>
    );
  }

  renderEditor() {
    return (
      <div className="editor">
        <Editor
          schema={schema}
          state={this.state.state}
          onChange={this.onChange}
        />
      </div>
    );
  }

  render() {
    return (
      <div>
        {this.renderMenu()}
        {this.renderEditor()}
      </div>
    );
  }
}

export default RichEditor;
