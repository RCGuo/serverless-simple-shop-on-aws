import { useEffect, useState } from 'react';
import { Card, Row, Col } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import AddToCart from '../button/AddToCart';
import FavoriteIcon from '../button/FavoriteIcon';
import StarRating from '../starRating/StarRating';

const TopicCard = (props) => {
  return (
    <Card as={Link} to={props.link} className='mx-4 topic-card' key={props.productId}>
      <Card.Img src={props.img} />
      <Card.Body className='px-0'>
        <Card.Title><h4>{props.title}</h4></Card.Title>
      </Card.Body>
    </Card>
  );
}

const ProductCard = (props) => {
  const [favorState, setFavorState] = useState(!!props.favoriteStatus);

  useEffect(() => {
    setFavorState(!!props.favoriteStatus);
  }, [props.favoriteStatus])

  return (
    <Card className='product-card'>
      <Card.Img className='pb-2' variant="top" src={props.imageFile} />
      <Card.Body className='p-0'>
        <Card.Title className='my-0'>
          <Row>
            <Col md={10} className="pt-1">{props.name}</Col>
            <Col md={2} className="px-0" id="product-favorite">
              <FavoriteIcon productId={props.productId} favorite={favorState} />
            </Col>
          </Row>
          <div><StarRating rating={props.rating} /></div>
        </Card.Title>
        <Card.Text as="div">
          <div className='mt-3' id="product-price">
            {props.price.toLocaleString('en-US', { style: 'currency', currency: 'usd' })}
          </div>
        </Card.Text>
        <AddToCart productId={props.productId} price={props.price} mark={"addToCart"}></AddToCart>
      </Card.Body>
    </Card>
  );
}

export {
  TopicCard,
  ProductCard,
}