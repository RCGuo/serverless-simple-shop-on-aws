import { API } from "aws-amplify";
import React, { useState, useEffect } from 'react';
import { useAuthenticator } from '@aws-amplify/ui-react';
import { Container, Row, Col } from 'react-bootstrap';
import { useLocation  } from 'react-router-dom';
import CardGallery from '../components/card/CardGallery';

const SearchView = () => {
  const {state} = useLocation();
  const [rsltNum, setRsltNum] = useState(0);
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

    const handleSearch = async () => {
      setLoadingState(true);
      try {
        await API.get("search", "/search/fuzzy-query", {
          queryStringParameters: {
            query: state,
            index: "product"
          }
        }).then((data) => {
          setProductList(data);
          if ( authStatus === 'authenticated' && data.length > 0 ) {
            getFavorList();
          }
          setRsltNum(data.length);
        });
        setLoadingState(false);
      } catch (e) {
        console.log("Product search error: ", e);
      }
    }

    handleSearch();
  }, [authStatus, state, setRsltNum]);

  return (
    <Container className="category-view mt-5">
      <Row>
        <Col className='category-view-title'>
          <h2><b>Search results</b></h2>
          <h6>About {rsltNum} results</h6>
        </Col>
      </Row>
      <Container className="g-0 mt-5 card-container">
        <CardGallery productList={productList} favoriteObject={favoriteObject} isLoading={isLoading}  mark="search" />
      </Container>
    </Container>
  );
}

export default SearchView;