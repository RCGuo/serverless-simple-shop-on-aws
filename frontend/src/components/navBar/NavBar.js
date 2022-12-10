import { Auth } from "aws-amplify";
import { useState, useEffect, useCallback } from 'react';
import { useAuthenticator } from '@aws-amplify/ui-react';
import { Nav, Navbar, NavItem, NavDropdown, Container, Row, Button, Form } from 'react-bootstrap';
import { Link, useNavigate } from 'react-router-dom';
import { FaShoppingCart } from 'react-icons/fa';
import { VscAccount } from 'react-icons/vsc';
import { BsSearch } from 'react-icons/bs';
import { AiOutlineSmile }from 'react-icons/ai';
import './navbar.css';

import brandIcon from './../../resources/images/logo/home_logo.png';

const NavBar = (props) => {
  const [classState, setClass] = useState('');
  const [userInfo, setUserInfo] = useState({email:"", nickname:""});
  const [navOffsetState, setNavOffsetTop] = useState(0);
  const [queryString, setQueryString] = useState("");
  const { user, authStatus } = useAuthenticator(context => [context.authStatus]);
  const navigate = useNavigate();

  useEffect(() => {
    if ( authStatus === 'authenticated' ) {
      setUserInfo({
        email: user.attributes.email,
        nickname: user.attributes.nickname,
      });
    }
    const handleScroll = () => {
      if ( window.pageYOffset > navOffsetState ) {
        setClass("fixed-top justify-content-center");
      } else {
        setClass('');
      }
    };
    window.addEventListener("scroll", handleScroll);
  }, [navOffsetState, user, authStatus]);

  const navcb = useCallback((node) => {
      if (node !== null) {
        setNavOffsetTop(node.offsetTop);
      }
  }, []);

  const handleLogout = async (e) => {
    e.preventDefault();
    try {
      await Auth.signOut(e)
      .then(() => {
        // navigate("/auth");
      });
    } catch (error) {
      console.log('error signing out: ', error);
    }
  };

  const handleLogIn = () => {
    navigate("/auth")
  }

  const userDropdown = () => {
    return (
      <>
        <NavDropdown 
          title={
            authStatus !== 'authenticated' ?  
              <span id="account-title"><VscAccount size={20}/>{' '}
                <span className="orange">Account</span></span>
            : <span id="account-title"><AiOutlineSmile style={{fontSize:"1.4rem"}}/>{'  '}Hi!{'   '}{userInfo.nickname}</span>
          } 
          id="account-nav-dropdown"
        >
          <NavDropdown.Item as={Link} to="/past-orders">
            Past order
          </NavDropdown.Item>
          <NavDropdown.Item as={Link} to="/favorites">
            My favorites
          </NavDropdown.Item>
          <NavDropdown.Item as={Link} to="/account">Account setting</NavDropdown.Item>
          <NavDropdown.Divider />
          { authStatus !== 'authenticated' ?  
            <NavDropdown.Item>
              <Nav.Link className="p-0" as="span" onClick={handleLogIn}>Sign in</Nav.Link>
            </NavDropdown.Item>
          : <NavDropdown.Item>
              <Nav.Link className="p-0" as="span" onClick={(e) => handleLogout(e)}>Sign out</Nav.Link>
            </NavDropdown.Item>}
        </NavDropdown>
      </>
    );
  }
  
  const showLoggedInBar = () => {
    return (
      <>
        <NavItem className="pr-1">
          <Nav.Link eventKey="past" as={Link} to={"/best"}>
            <span className="orange text-nowrap">Best sellers</span>
          </Nav.Link>
        </NavItem>
        <NavItem className="pr-1">
          <Nav.Link eventKey="past" as={Link} to={"/past"}>
            <span className="orange text-nowrap">Past orders</span>
          </Nav.Link>
        </NavItem>
        {userDropdown()}
        <NavItem>
          <Nav.Link eventKey="cart" as={Link} to={"/cart"}>
            <div className="shopping-cart-icon-container">
              <FaShoppingCart fontWeight="bold" color="white" size={18} />
            </div>
          </Nav.Link>
        </NavItem>
      </>
    );
  };
  
  const showLoggedOutBar = () => {
    return (
      <>
        <NavItem className="pr-1">
          <Nav.Link eventKey="past" as={Link} to={"/best-sellers"}>
            <span className="orange text-nowrap">Best sellers</span>
          </Nav.Link>
        </NavItem>
        {userDropdown()}
        <NavItem>
          <Nav.Link eventKey="cart" as={Link} to={"/cart"}>
            <div className="shopping-cart-icon-container">
              <FaShoppingCart fontWeight="bold" color="white" size={18} />
            </div>
          </Nav.Link>
        </NavItem>
      </>
    );
  };

  const handleOnChange = (e) => {
    e.preventDefault();
    setQueryString(e.target.value.trim())
  }

  const handleSubmit = (e) => {
    e.preventDefault();
    if( queryString ){
      navigate("/search", {
        state: queryString,
      });
    }
  }

  return (
    <Container className="px-4">
      <Row>
          <Navbar className={classState} expand="lg" bg="light" id="main-nav" ref={navcb}>
            <Container>
              <Navbar.Brand href="/home">
                <span className="orange">
                  <img src={brandIcon} 
                       className="d-inline-block align-top" 
                       alt="SimpleShop logo" /> 
                  {' '}SimpleShop
                </span>
              </Navbar.Brand>
              <Navbar.Toggle aria-controls="basic-navbar-nav"/>
              <Nav>
                <Form className="d-flex search" style={{width:"23rem"}} onSubmit={handleSubmit}>
                  <Form.Control
                    type="search"
                    placeholder={ queryString === "" ? "Search" : queryString }
                    aria-label="Search"
                    onChange={handleOnChange}
                  />
                  <Button 
                    className="search-bar-icon-btn" 
                    variant="warning" 
                    type="submit"
                  >
                    <BsSearch />
                  </Button>
                </Form>
              </Nav>
              <Navbar.Collapse>
                <Nav className="ms-auto">
                  { props.isAuthenticatedState ? showLoggedInBar() : showLoggedOutBar() }
                </Nav>
              </Navbar.Collapse>
            </Container>
          </Navbar>
      </Row>
    </Container>
  );
}

export default NavBar;