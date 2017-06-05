import React, { Component } from 'react';

class Editable extends Component {
  constructor(props) {
    super(props);

    this.state = {
      isEditing: false,
      input: null,
    };

    this.onSave = this.onSave.bind(this);
    this.onKeyPress = this.onKeyPress.bind(this);
    this.onRef = this.onRef.bind(this);
    this.toggle = this.toggle.bind(this);
  }

  componentWillMount() {
    this.editor = this.props.editor || DefaultEditor;
  }

  componentDidUpdate(prevProps, prevState = {}) {
    if (this.state.isEditing && !prevState.isEditing) {
      document.addEventListener('keydown', this.onKeyPress);
    } else if (!this.state.isEditing && prevState.isEditing) {
      document.removeEventListener('keydown', this.onKeyPress);
    }
  }

  componentWillUnmount() {
    document.removeEventListener('keydown', this.onKeyPress);
  }

  onSave() {
    const value = this.input.value.trim();

    if (value && value !== this.props.value) {
      this.props.onSave(value);
    }

    this.toggle();
  }

  onKeyPress(e) {
    if (e.key === 'Enter') {
      this.onSave();
      return;
    }

    if (e.keyCode === 27) { // esc
      this.toggle();
    }
  }

  onRef(node) {
    this.input = node;
  }

  toggle() {
    const newState = !this.state.isEditing;
    this.setState({
      isEditing: newState
    });

    if (this.props.onToggle) {
      this.props.onToggle(newState);
    }
  }

  render() {
    const { value, children, className, editorClassName, childrenClassName, style } = this.props;
    const { isEditing } = this.state;

    const Editor = this.editor;

    return (
      <div className={className}>
        { !isEditing ?
          <div className={childrenClassName} onClick={this.toggle}>{ children }</div> :
          <Editor
            autoFocus
            style={style}
            className={editorClassName}
            defaultValue={value}
            onBlur={this.onSave}
            onKeyPress={this.onKeyPress}
            onRef={this.onRef}
          />
        }
      </div>
    );
  }
}

const DefaultEditor = ({ onRef, ...rest }) => (
  <input
    type="text"
    ref={onRef}
    {...rest}
  />
);

export default Editable;
