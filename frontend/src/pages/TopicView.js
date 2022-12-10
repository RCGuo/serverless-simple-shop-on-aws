import { API } from "aws-amplify";
import React, { useState, useEffect } from 'react';
import { useAuthenticator } from '@aws-amplify/ui-react';
import { useParams } from "react-router-dom";
import { Container, Row, Col } from 'react-bootstrap';
import CardGallery from '../components/card/CardGallery';
import { Topics } from './../components/category/categoryConfig';

const TopicView = () => {
  let { topic } = useParams();
  const [isLoading, setLoadingState] = useState(true);
  const [productList, setProductList] = useState([]);
  const [favoriteObject, setFavoriteObject] = useState({});
  const { authStatus } = useAuthenticator(context => [context.authStatus]);
  const topicTitle = Topics[topic] ? Topics[topic] : "Unknown topic";

  useEffect(() => {
    const getProductList = async () => {
      try {
        await API.get("shopApi", "/product", {
          queryStringParameters: {
            topic: topic,
          }
        })
        .then((data) => {
          setProductList(data);
          if ( authStatus === 'authenticated' && data.length > 0 ) {
            getFavorList();
          }
          setLoadingState(false);
        });
      } catch (e) {
        console.log(e);
      }
    }

    const getFavorList = async () => {
      try {
        await API.get("shopApi", "/product/favorites")
        .then((data) => {
          const favoritesMapping = {};
          data.forEach(
            (item, i) => {
              favoritesMapping[item.productId] = true;
            });
            setFavoriteObject(favoritesMapping);
        });
      } catch (e) {
        console.log(e);
      }
    }
    getProductList();
  }, [topic, authStatus]);

  return (
    <>
    <Container className="topic-view mt-5">
      <Row className='mt-3'>
        <Col className='topic-view-title'>
           <h2><b>{topicTitle}</b></h2>
        </Col>
      </Row>

      <Container className="g-0 mt-5 card-container">
        <CardGallery productList={productList} favoriteObject={favoriteObject} isLoading={isLoading} mark="topic"/>
      </Container>
    </Container> 
    </>
  );
}

export default TopicView;