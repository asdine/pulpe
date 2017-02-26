import { Collapse, Navbar, NavbarToggler, Nav, NavItem, NavLink } from 'reactstrap';
import React from 'react';
import BoardDetail from '../containers/board';
import BoardList from '../containers/boardList';
import Modals from '../containers/modal';

const Home = ({ children }) =>
  <div className="wrapper container-fluid">
    <MobileMenu />
    <Menu />
    <BoardDetail />
    {children}
    <Modals />
  </div>;

const Menu = () =>
  <div className="plp-left-menu">
    <h1>pulpe</h1>
    <BoardList />
  </div>;

class MobileMenu extends React.Component {
  constructor(props) {
    super(props);

    this.toggleNavbar = this.toggleNavbar.bind(this);
    this.state = {
      collapsed: true
    };
  }

  toggleNavbar() {
    this.setState({
      collapsed: !this.state.collapsed
    });
  }
  render() {
    return (
      <div className="plp-top-menu">
        {/* <h1>pulpe</h1>*/}
        <Navbar light>
          <NavbarToggler onClick={this.toggleNavbar} />
          <h1>pulpe</h1>
          <Collapse className="navbar-toggleable-md" isOpen={!this.state.collapsed}>
            <Nav navbar>
              <NavItem>
                <NavLink href="/test/">Tests</NavLink>
              </NavItem>
              <NavItem>
                <NavLink href="/paris/">Paris</NavLink>
              </NavItem>
              <NavItem>
                <NavLink href="/paris/">+ Create a board</NavLink>
              </NavItem>
            </Nav>
          </Collapse>
        </Navbar>
      </div>
    );
  }
}

export default Home;
