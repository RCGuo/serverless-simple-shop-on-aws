import React from 'react';
import { Nav, NavItem, Container, Row, Col } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import './banner.css';

import ad50Off from "../../resources/images/advert/50-off-template-for-advertising.jpg";
import fashionWeek from "../../resources/images/advert/fashion_week_banner.jpg";
import specialOffer from "../../resources/images/advert/special_offer_banner.jpg";

const AdBanner = () => {
  return (
    <Container className='mt-5'>
      <Row>
        <Col>
          <h4>AD Banner</h4>
        </Col>
      </Row>
      <Row>
        <Col>
          <Nav>
            <NavItem>
              <Nav.Link as={Link} to={"#fashionWeek"}>
                <img className="img-fluid" src={fashionWeek} alt="fashion week" />
              </Nav.Link>
            </NavItem>
          </Nav> 
        </Col>
        <Col>
          <Nav>
            <NavItem>
              <Nav.Link as={Link} to={"#specialOffer"}>
                <img className="img-fluid" src={specialOffer} alt="special offer" />
              </Nav.Link>
            </NavItem>
          </Nav> 
        </Col>
        <Col>
          <Nav>
            <NavItem>
              <Nav.Link as={Link} to={"#ad50Off"}>
                <img className="img-fluid" src={ad50Off} alt="ad 50Off" />
              </Nav.Link>
            </NavItem>
          </Nav> 
        </Col>
      </Row>
    </Container>
  );
}

export default AdBanner;