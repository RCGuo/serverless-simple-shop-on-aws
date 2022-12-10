import { Spinner } from 'react-bootstrap';
import { ProductCard } from './Card';
import './card.css';

const CardGallery = ({ productList, favoriteObject, isLoading, mark }) => {
  return (
    <>
      { 
        isLoading === true ? 
          <Spinner animation="border" variant="secondary" size="lg" />
          : 
          mark === "search" ?
            productList.map( (product, index) => {
              let favoriteStatus = !!favoriteObject[product.productId];
              return <ProductCard {...product._source} favoriteStatus={favoriteStatus} key={product._source.productId}/>;
            })
            : 
            productList.map( (product, index) => {
              let favoriteStatus = !!favoriteObject[product.productId];
              return <ProductCard {...product} favoriteStatus={favoriteStatus} key={product.productId}/>;
            })
      }
    </>
  );
}

export default CardGallery;