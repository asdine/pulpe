import React from 'react';

function PlaceholderPlugin() {
  return {
    schema: {
      rules: [
        BLOCK_RENDER_RULE,
      ]
    }
  };
}

export default PlaceholderPlugin;

const BLOCK_RENDER_RULE = {
  match: (node) => node.kind === 'block',
  render: (props) => (
    <div {...props.attributes} style={{ position: 'relative' }}>
      {props.children}
      <Placeholder
        node={props.node}
        parent={props.state.document}
        state={props.state}
      >
        {props.editor.props.placeholder}
      </Placeholder>
    </div>
  )
};


class Placeholder extends React.Component {
  shouldComponentUpdate = (props) => (
    props.children !== this.props.children ||
    props.className !== this.props.className ||
    props.node !== this.props.node ||
    props.style !== this.props.style ||
    props.state.document.length !== this.props.state.document.length
    )


  isVisible = () => {
    const { state } = this.props;
    return state.document.isEmpty && state.document.nodes.size <= 1 && state.document.nodes.first() === this.props.node;
  }

  render = () => {
    const isVisible = this.isVisible();
    if (!isVisible) return null;

    const { children, className } = this.props;
    let { style } = this.props;

    if (typeof children === 'string' && style == null && className == null) {
      style = { opacity: '0.333' };
    } else if (style == null) {
      style = {};
    }

    const styles = {
      position: 'absolute',
      top: '0px',
      right: '0px',
      bottom: '0px',
      left: '0px',
      pointerEvents: 'none',
      ...style
    };

    return (
      <span contentEditable={false} className={className} style={styles}>
        {children}
      </span>
    );
  }

}
