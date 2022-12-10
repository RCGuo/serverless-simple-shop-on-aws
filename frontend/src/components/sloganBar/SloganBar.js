import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import './sloganBar.css';

const SloganBar = () => {
  return (
    <Container className="slogan-bar-container" id="main-search">
      <Row>
        <Col>
          <h3 className="slogan-bar white">Online shopping<span className="orange">{` Simple `}</span></h3>
        </Col>
      </Row>
    </Container>
  );
}

export default SloganBar;