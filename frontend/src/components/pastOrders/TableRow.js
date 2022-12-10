import { API } from "aws-amplify";
import { useState } from "react";
import { Spinner, Image } from "react-bootstrap";
import { FaRegTrashAlt, FaShoppingCart } from "react-icons/fa";
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './pastOrders.css';

import visa from "./../../resources/images/creditCardsIcons/visa_electron.png";
import master from "./../../resources/images/creditCardsIcons/mastercard.png";
import jcb from "./../../resources/images/creditCardsIcons/jcb.png";
import discover from "./../../resources/images/creditCardsIcons/discover.png";
import american_express from "./../../resources/images/creditCardsIcons/american_express.png";
import FavoriteIcon from "../button/FavoriteIcon";

const cardIcon = {
  visa: visa,
  mastercard: master,
  jcb: jcb,
  discover: discover,
  amex: american_express,
}

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return `${date.getMonth()+1}/${date.getDate()}/${date.getFullYear()} 
          ${date.getHours()}:${date.getMinutes() < 10 ? '0' : ''}
          ${date.getMinutes()}`;
}

const PastOrderTableRow = ({ order }) => {
  return (
    <tr className="text-center">
      <td>{formatDate(order.orderDate)}</td>
      <td className="truncate">{order.orderId}</td>
      <td>
        <div><Image src={cardIcon[order.paymentMethod.brand]}></Image></div>
      </td>
      <td>{(order.total/100).toLocaleString('en-US', { style: 'currency', currency: 'usd' })}</td>
    </tr>   
  );
}

const FavoriteTableRow = ({ product, getFavorites, mark }) => {
  const [isRemoving, setIsRemoving] = useState(false);
  const [isAdding, setIsAdding] = useState(false);

  const onRemoveFavorite = async ({productId}) => {
    setIsRemoving(true);
    await API.post("shopApi", "/product/favorites", {
      body: {
        productId: productId,
        favorite: false
      },
    });
    getFavorites();
  }

  const notify = () => {
    toast.success('Added to Cart', {
      position: "top-center",
      autoClose: 4000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
      theme: "light",
    });
  }

  const onAddToCart = async (e, props) => {
    setIsAdding(true);
    const productInCart = await API.get("shopApi", `/cart/${props.productId}`);

    if (productInCart.productId === props.productId) {
      API.put("shopApi", "/cart", {
        body: {
          productId: props.productId,
          quantity: productInCart.quantity + 1,
        },
      }).then(() => {
        setIsAdding(false);
        notify();
      });
    } else {
      await API.post("shopApi", "/cart", {
        body: {
          productId: props.productId,
          price: props.price,
          quantity: 1,
        },
      }).then(() => {
        setIsAdding(false);
        notify();
      });
    }
  }

  return (
    <tr className="text-start">
      <td className="text-center" width="20%">
        <Image id="product-image" fluid={true} src={product.imageFile} thumbnail={true} />
      </td>
      <td>
        <div id="favorite-name">
          {product.name}
        </div>
        <div>
        {product.company}
        </div>
      </td>
      <td id="favorite-price">{product.price.toLocaleString('en-US', { style: 'currency', currency: 'usd' })}</td>
      <td> 
        <div className="mb-4 text-center">          
            {isAdding
              ? <Spinner animation="border" variant="secondary" size="sm" />
              : <button id="favorite-shopping">
                  <FaShoppingCart onClick={(e) => onAddToCart(e, product)} size={23} />
                </button>}
        </div>
        {/* FavoriteIcon */}
        
            <div className="text-center">    
            { mark === "for_favorite" ?      
                isRemoving 
                  ? <Spinner animation="border" variant="secondary" size="sm" />
                  : <button id="favorite-trash">
                      <FaRegTrashAlt size={23} onClick={() => onRemoveFavorite({productId:product.productId})} />
                    </button>
              : <FavoriteIcon />
            }
            </div>
          
      </td>
    </tr>
  );
}

export {PastOrderTableRow, FavoriteTableRow};