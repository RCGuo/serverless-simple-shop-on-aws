import { Container, Row, Col, Image } from 'react-bootstrap';
import { Categories } from "./categoryConfig";
import { Link } from 'react-router-dom';
import './category.css';

import elects from "../../resources/images/featuredCategory/elects.jpg";
import fashion from "../../resources/images/featuredCategory/fashion.jpg";
import sport from "../../resources/images/featuredCategory/sport.jpg";
import toys from "../../resources/images/featuredCategory/toys.jpg";
import pet from "../../resources/images/featuredCategory/pet.jpg";
import house from "../../resources/images/featuredCategory/house.jpg";

const images = {
  elects: elects,
  fashion: fashion,
  sport: sport,
  toys: toys,
  pet: pet,
  house: house,
} 

const FeaturedCategory = () => {
  return (
    <Container className="mt-5">
      <Row className='mt-3'>
        <Col className='featured-category-title'>
          <h2><b>Featured categories</b></h2>
        </Col>
      </Row>
      <Row>
        <Col className="my-3 featured-category-container">
          <ul>
            { Object.keys(images).map((key, index) => {
              return (
                <li className='bounce' key={key} as={Link} to={`/category/${key}`}>
                  <Link to={`/category/${key}`}>
                  <div className='circle-image-container'>
                    <Image fluid src={images[key]} alt={Categories[key]} />
                  </div>
                  <div className='mt-1'>
                    <h5>{Categories[key]}</h5>
                  </div>
                  </Link>
                </li>
              );
            })}
          </ul>
        </Col>
      </Row>
    </Container>
  );
}

export default FeaturedCategory;