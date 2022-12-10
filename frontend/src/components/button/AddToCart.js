import { API } from "aws-amplify";
import { useAuthenticator } from '@aws-amplify/ui-react';
import { useEffect, useState } from "react";
import { Button, Spinner } from "react-bootstrap";
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';

const AddToCart = (props) => {
  const [isloading, setIsLoading] = useState(false);
  const [toCartState, setToCartState] = useState(false);
  const { authStatus } = useAuthenticator(context => [context.authStatus])
  const navigate = useNavigate();
  
  useEffect(() => {
    if (toCartState) {
      navigate("/cart");
    }
  }, [toCartState, navigate])

  const notify = () => {
    toast.warn('Please login before adding a product to cart', {
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

  const onAddToCart = async () => {
    if (authStatus !== 'authenticated') {
      notify();
      return;
    }
    setIsLoading(true);
    const productInCart = await API.get("shopApi", `/cart/${props.productId}`);
    if (productInCart.productId === props.productId) {
      API.put("shopApi", "/cart", {
        body: {
          productId: props.productId,
          quantity: productInCart.quantity + 1,
        },
      }).then(() => {
        setToCartState(true);
      });
    } else {
      await API.post("shopApi", "/cart", {
        body: {
          productId: props.productId,
          price: props.price,
          quantity: 1,
        },
      }).then(() => {
        setToCartState(true)
      });
    }
  }

  return (
    <Button
      size="sm" 
      variant="danger"
      disable={`isloading`}
      onClick={onAddToCart}
    >
      {isloading 
        ? <Spinner animation="border" variant="secondary"  size="sm" /> 
        : props.mark === "buyAgain" ? "Buy Again"  : "Add to cart"}
    </Button>
  );
}

export default AddToCart;