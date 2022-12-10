import { API } from "aws-amplify";
import React, { useEffect, useState } from "react";
import { Container, Image, Spinner } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import './bestSeller.css';

import { Pagination, Navigation } from "swiper";
import { Swiper, SwiperSlide } from 'swiper/react';
import 'swiper/css';
import 'swiper/css/pagination';
import 'swiper/css/navigation';

const BestSellerBar = () => {
  const [isLoading, setLoadingState] = useState(true);
  const [top10, setTop10] = useState([]);
  
  useEffect(() => {
    const getTop10 = async () => {
      try {
        await API.get("search", "/search/purchased-ranking", {
          queryStringParameters: {
            size: 10,
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
            setTop10(data);
            setLoadingState(false);
          });
        });
      } catch (e) {
        console.log(e);
      }
    }
    getTop10();
  }, [])

  return (
    <Container className="best-seller p-0 mt-3">
    { isLoading 
      ? <Spinner animation="border" variant="secondary" size="lg" />
      : 
      <Swiper
        slidesPerView={5}
        slidesPerGroup={2}
        pagination={{
          clickable: true,
        }}
        navigation={true}
        scrollbar={{draggable: true}}
        modules={[Pagination, Navigation]}
      >
        <SwiperSlide className="best-heading-slide text-start px-4">
          <div>
            <b>Best Sellers</b>
            <div id="subtitle">
              <h6>a product that is extremely popular and has sold in very large numbers</h6>
            </div>
          </div>
        </SwiperSlide>
          { top10.map((product) => {
            return (
              <SwiperSlide key={product.name}>
                <Link as={Link} to={"/best-sellers"}>
                  <Image fluid src={product.imageFile} alt={product.name} />
                </Link>
              </SwiperSlide>
            );
          } ) }
      </Swiper>
    }
    </Container>
  );

}

export default BestSellerBar;