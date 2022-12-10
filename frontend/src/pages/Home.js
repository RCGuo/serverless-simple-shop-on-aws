import React from 'react'
import { Container, Row, Col } from 'react-bootstrap';
import Hero from '../components/hero/Hero';
import SloganBar from '../components/sloganBar/SloganBar';
import BestSellerBar from '../components/bestSeller/BestSellerBar';
import TopicBanner from '../components/banner/TopicBanner';
import AdBanner from '../components/banner/AdBanner';
import CategoryBoeard from '../components/category/FeaturedCategory';

const Home = () => {
  return (
    <Container>
      <Row>
        <Col className="text-center">
          <Container className="shop-home">
            <Hero />
            <SloganBar />
            <BestSellerBar />
            <TopicBanner />
            <AdBanner />
            <CategoryBoeard />
          </Container> 
        </Col>
      </Row>
    </Container>
  );
}

export default Home;