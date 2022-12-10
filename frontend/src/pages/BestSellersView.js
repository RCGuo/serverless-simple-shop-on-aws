import { API } from "aws-amplify";
import React, { useState, useEffect } from 'react';
import { useAuthenticator } from '@aws-amplify/ui-react';
import { Container, Row, Col } from 'react-bootstrap';
import CardGallery from '../components/card/CardGallery';

const BestSellersView = () => {
  const [isLoading, setLoadingState] = useState(true);
  const [productList, setProductList] = useState([]);
  const [favoriteObject, setFavoriteObject] = useState({});
  const { authStatus } = useAuthenticator(context => [context.authStatus]);

  useEffect(() => {
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

    const handleBest = async () => {
      try {
        await API.get("search", "/search/purchased-ranking", {
          queryStringParameters: {
            size: 20,
            index: "product-sold-counter",
            sortField: "counter",
            sortDirect: "desc",
          }
        }).then((data) => {
          let productIdlist = data.map((obj, i) => {
            return {
              productId: obj._id,
            };
          });
          API.post("shopApi", '/product/batch-fetch', {
            body: productIdlist,
          }).then((data) => {
            setProductList(data);
            if ( authStatus === 'authenticated' && data.length > 0 ) {
              getFavorList();
            }
            setLoadingState(false);
          });
        });
      } catch (e) {
        console.log(e);
      }
    }

    handleBest();
  }, [authStatus]);

  return (
    <Container className="topic-view mt-5">
      <Row className='mt-3'>
        <Col className='topic-view-title'>
          <h2><b>Top 20 best sellers</b></h2>
        </Col>
      </Row>
      <Container className="g-0 mt-5 card-container">
        <CardGallery productList={productList} favoriteObject={favoriteObject} isLoading={isLoading} mark="best"/>
      </Container>
    </Container>
  );
}

export default BestSellersView;