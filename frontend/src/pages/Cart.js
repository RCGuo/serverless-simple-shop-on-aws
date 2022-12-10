import { API } from "aws-amplify";
import { useEffect, useState, useCallback  } from "react";
import { useAuthenticator } from '@aws-amplify/ui-react';
import { Col, Container, Row, Table, Spinner, Button, Image } from "react-bootstrap";
import { useNavigate } from 'react-router-dom';
import CartTableRow from "../components/cart/CartTableRow";
import './../App.css';
import emptyCart from './../resources/images/common/empty_cart.png'

const Cart = () => {
  const [isLoading, setIsLoading] = useState(true);
  const [isChecout, setIsCheckout] = useState(false);
  const [ordersInCart, setOrdersInCart] = useState([]);
  const [isTotalUpdating, setTotalUpdating] = useState(false);
  const [totalPrice, setTotalPrice] = useState(0);
  const [totalItemQty, setTotalItemQty] = useState(0);
  const { authStatus } = useAuthenticator(context => [context.authStatus])
  const navigate = useNavigate();
  const deliveryFee = 0;

  const getOrderTotal = useCallback( async () => {
    setTotalUpdating(true)
    await API.get("shopApi", "/cart")
    .then((ordersInCart) => {
      setOrdersInCart(ordersInCart);
      const totalPrice = ordersInCart.reduce((total, order) => {
        return total + order.price * order.quantity;
      }, deliveryFee);
      setTotalPrice(parseFloat(totalPrice.toFixed(2)));
      const totalQty = ordersInCart.reduce((total, order) => {
        return total + order.quantity;
      }, 0)
      setTotalItemQty(totalQty);
      setTotalUpdating(false)
    });
  }, []);

  useEffect(() => {
    if (isChecout) {
      navigate("/checkout");
    }
    const getOrdersInCart = async () => {
      try {
        await API.get("shopApi", "/cart")
        .then((ordersInCart) => {
          setOrdersInCart(ordersInCart);
        });
        setIsLoading(false);
      } catch (e) {
        console.log(e);
      }
    }

    if ( authStatus === 'authenticated' ) {
      getOrdersInCart();
      getOrderTotal();
    } else {
      setIsLoading(false);
    }
  }, [authStatus, getOrderTotal, isChecout, navigate]);

  const onCheckout = () => {
    setIsCheckout(true);
  }

  if ( isLoading ) {
    return (
      <Container className="text-center">
        <Row>
          <Col>
            <Spinner animation="border" variant="secondary" size="lg" />
          </Col>
        </Row>
      </Container>
    ); 
  }

  if (ordersInCart.length === 0 | authStatus !== 'authenticated') {
    return (
      <Container className="cart-container mt-5 pt-4">
        <Row className="text-center">
          <Col className="cart-title">
            <Image src={emptyCart}  />
          </Col>
        </Row>
      </Container>
    );
  }

  return (
    <Container className="cart-container mt-5 pt-4">
      <Row className="justify-content-center">
        <Col lg={9} className="cart-title">
          Shopping Cart
        </Col>
      </Row>
      <Row className="justify-content-center">
        <Col lg={9}>
          <Table bordered className="orders-in-cart">
            <thead>
              <tr>
                <th width="2%"></th>
                <th width="16%"></th>
                <th width="45%"></th>
                <th width="22%">Unit Price</th>
                <th width="15%">Quantity</th>
              </tr>
            </thead>
            <tbody>
              { !isLoading && 
                ordersInCart.map((data, index) => 
                  <CartTableRow 
                    order={data} 
                    key={data.productId} 
                    getOrderTotal={getOrderTotal} />) 
              }
            </tbody>
          </Table>
          <Table className="order-fee">
            <tbody>
              <tr className="text-end"> 
                <td id="fee-title">
                  {isTotalUpdating && <Spinner animation="border" variant="secondary" size="sm" />}
                  {' '} Subtotal ({totalItemQty} items): {' '}
                </td>
                <td id="product-price">
                  {totalPrice.toLocaleString('en-US', { style: 'currency', currency: 'usd' })}
                </td>
              </tr>
              <tr className="text-end">
                <td id="product-subtotal">
                  Shipping: 
                </td>
                <td id="product-price">
                  {deliveryFee.toLocaleString('en-US', { style: 'currency', currency: 'usd' })}
                </td>
              </tr>
              <tr className="text-end"> 
                <td id="product-subtotal">
                  Total: 
                </td>
                <td id="product-price">
                  {totalPrice.toLocaleString('en-US', { style: 'currency', currency: 'usd' })}
                </td>
              </tr>
            </tbody>
          </Table>
        </Col>
      </Row>
      <Row className="justify-content-center">
        <Col lg={9} className="text-end cart-checkout">
          <Button 
            onClick={onCheckout} 
            variant="warning" 
            disabled={isLoading}
            id="cart-checkout-btn">
            Proceed to checkout
          </Button>
        </Col>
      </Row>
    </Container>
  );
}

export default Cart;