import React from 'react';
import { Container, Row, Col, CardGroup } from 'react-bootstrap';
import { TopicCard }  from '../card/Card'
import './banner.css';

const topicImages = [
  {
    title: "outerwear",
    link: "/topic/winter",
    img: process.env.REACT_APP_IMAGE_CDN + "/winter/womens_brown_coat.jpg",
  },
  {
    title: "boots",
    link: "/topic/winter",
    img: process.env.REACT_APP_IMAGE_CDN + "/winter/womens_elle_boot.jpg",
  },
  {
    title: "sweaters",
    link: "/topic/winter",
    img: process.env.REACT_APP_IMAGE_CDN + "/winter/blue_christmas_sweater.jpg",
  },
  {
    title: "pajamas & socks",
    link: "/topic/winter",
    img: process.env.REACT_APP_IMAGE_CDN + "/winter/womens_pajama_sets_long_sleeve.jpg",
  },
];

const TopicBanner = () => {

  const cardList = topicImages.map((topic) => {
    return (
      <TopicCard key={topic.title} {...topic}/>
    );
  });

  return (
    <Container className='mt-4 pt-3 pb-4 topic-banner'>
      <Row>
        <Col>
          <div className='topic-banner-title'>
            <h2><b>Warm up this winter</b></h2>
          </div>
        </Col>
      </Row>
      <CardGroup className='mt-2 topic-card'>
        {cardList}
      </CardGroup>
    </Container>
  );
}

export default TopicBanner;