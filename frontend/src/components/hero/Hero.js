import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import Image from 'react-bootstrap/Image'
import image from './../../resources/images/hero/we_are_open_sign.jpg';
import './hero.css';

const Hero = () => {
  return (
    <Container className="px-0">
      <Row>
        <Col>
          <Image src={image} fluid={true} />
        </Col>
      </Row>
    </Container>
  );
}

export default Hero;