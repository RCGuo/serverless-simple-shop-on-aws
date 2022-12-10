import React from 'react'
import { Container, Row, Col, Button, Image } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import screenshot from './../resources/images/screenshot/simple_shop_demo_screenshot.png';


const Landing = () => {
  return (
    <Container className='landing-container mt-5'>
      <Row>
        <Col className="text-center">
          <div className="lander">
            <h1>Simple Shop Demo</h1>
          </div>
        </Col>
      </Row>
      <Row className='mt-5'>
        <Col>
          <div>
            <p>This is cloud-native serverless demo application, with basic E-commerce features such as a shopping cart, product searching, product gallery, favorites feature, and stripe online payment.</p>
            <p>It is a full-stack application implements event-driven microservices architecture with using AWS serverless services <i>for demonstration purpose only</i>. Its backend is written in Golang with a frontend UI written in ReactJS. Such an architecture reduces operational costs, and management overheads, and enables high scalability and availability without additional effort from the developer.</p>
          </div>
        </Col>
      </Row>
      <Row className='mt-4'>
        <Col>
          <div className='entrance text-center'>
            <Link to={"/home"}>
              <Button
                href="/home"
                variant="warning" 
                id="entrance">
                Enter Simple Shop Demo 
              </Button>
            </Link >
          </div>
        </Col>
      </Row>
      <Row className='mt-5'>
        <Col>
          <div>
            <Image fluid src={screenshot} />
          </div>
        </Col>
      </Row>
    </Container>
  );
}

export default Landing;