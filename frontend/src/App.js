import React, { useEffect } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import NavBar from './components/navBar/NavBar';
import RoutePaths from './Routes';
import Footer from './components/footer/Footer';
import CategoryNavBar from './components/category/CategoryNavBar';
import { ToastContainer } from 'react-toastify';
import './App.css';

const App = () => {
  useEffect(() => {
    document.title = "Serverless Simple Shop";
  }, []);

  return (
    <>
      <CategoryNavBar />
      <Container className="App">
        <Row className='align-items-center'>
          <Col>
            <NavBar />
            <RoutePaths />
          </Col>
        </Row>
        <ToastContainer />
      </Container>
      <Footer />
    </>
  );
}

export default App;