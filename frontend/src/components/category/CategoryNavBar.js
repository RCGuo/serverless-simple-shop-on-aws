import React from "react";
import { Nav, Container, Row, Col } from 'react-bootstrap';
import { Categories } from "./categoryConfig";
import './category.css';

const CategoryNavBar = () => {
  return (
    <Container fluid className="category-nav">
      <Row>
        <Col>
          <Nav className="justify-content-center">
            {Object.keys(Categories).map((key, index) =>
              <Nav.Item key={key}>
                <Nav.Link href={`/category/${key}`}>{Categories[key]}</Nav.Link>
              </Nav.Item>
            )}
          </Nav>
        </Col>
      </Row>
    </Container>
  );
}

export default CategoryNavBar;