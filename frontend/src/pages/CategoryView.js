import { API } from "aws-amplify";
import React, { useState, useEffect } from 'react';
import { useAuthenticator } from '@aws-amplify/ui-react';
import { Container, Row, Col } from 'react-bootstrap';
import { useParams } from "react-router-dom";
import CardGallery from '../components/card/CardGallery';
import { Categories } from './../components/category/categoryConfig';

const CategoryView = () => {
  let { category } = useParams();
  const [isLoading, setLoadingState] = useState(true);
  const [productList, setProductList] = useState([]);
  const [favoriteObject, setFavoriteObject] = useState({});
  const { authStatus } = useAuthenticator(context => [context.authStatus]);
  const categoryTitle = Categories[category] ? Categories[category] : "Unknown category";

  useEffect(() => {
    const getProductList = async () => {
      try {
        await API.get("shopApi", "/product", {
          queryStringParameters: {
            category: category,
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
  }, [category, authStatus]);

  return (
    <Container className="category-view mt-5">
      <Row>
        <Col className='category-view-title'>
          <h2><b>{categoryTitle}</b></h2>
        </Col>
      </Row>

      <Container className="g-0 mt-5 card-container">
        <CardGallery productList={productList} favoriteObject={favoriteObject} isLoading={isLoading} mark="category"/>
      </Container>
    </Container>
  );
}

export default CategoryView;