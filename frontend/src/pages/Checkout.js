import { API } from "aws-amplify";
import { useAuthenticator } from '@aws-amplify/ui-react';
import { useEffect, useState, useRef } from "react";
import { Container, Row, Col, Spinner, Table } from "react-bootstrap";
import { RequireAuth } from "../components/user/Authentication";
import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';
import { useNavigate } from "react-router";
import CheckoutForm from "../components/checkout/CheckoutForm";
import { HiClipboard } from 'react-icons/hi';
import { TiTick } from 'react-icons/ti';
import './../App.css';

const stripePromise = loadStripe(process.env.REACT_APP_STRIPE_PUBLISHABLE_KEY);
const Checkout = () => {
  const [clientSecret, setClientSecret] = useState("");
  const [ordersInCart, setOrdersInCart] = useState("");
  const [totalAmount, setTotalAmount] = useState(0);
  const [isLoading, setLoading] = useState(false);
  const { authStatus } = useAuthenticator(context => [context.authStatus]);
  const visaRef = useRef(null);
  const masterRef = useRef(null);
  const cvcFailedRef = useRef(null);
  const blockedRef = useRef(null);
  const [isActiveVisa, setIsActiveVisa] = useState(false);
  const [isActiveMaster, setIsActiveMaster] = useState(false);
  const [isActiveCvcFailed, setIsActiveCvcFailed] = useState(false);
  const [isActiveBlocked, setIsActiveBlocked] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    // Create PaymentIntent as soon as the page loads
    const getOrdersInCart = async () => {
      setLoading(true);
      try {
        const orderItemsInCartData = await API.get("shopApi", "/cart");
        setOrdersInCart(orderItemsInCartData);
        if ( ! orderItemsInCartData | orderItemsInCartData.length === 0 ) {
          navigate("/cart", { replace: true });
        } else {
          await API.post("shopApi", "/checkout/create-payment-intent", {
            body: { items: orderItemsInCartData },
          }).then((data) => {
            setClientSecret(data.clientSecret);
            if (!isNaN(data.total) && data.total> 0) {
              setTotalAmount(data.total/100);
            }
          });
        }
      } catch (e) {
        console.log(e);
      }
      setLoading(false);
    }

    if ( authStatus === 'authenticated' ) {
      getOrdersInCart();
    }
  }, [navigate, authStatus]);

  const appearance = {
    theme: 'stripe',
  };
  const options = {
    clientSecret,
    appearance,
  };
  
  if ( isLoading ) return (
    <Container className="checkout-container mt-5 pt-4">
      <Row className="text-center">
        <Col lg={9}>
          <Spinner animation="border" variant="secondary" size="lg" />
        </Col>
      </Row>
    </Container>
  );
  
  const copyToClipboard = (e, refObj) => {
    const code = refObj.current.innerText;
    navigator.clipboard.writeText(code);
    e.target.focus();
  };

  const setCopiedActive = (e, setActive) => {
    setActive(true);
    const timer = setTimeout(() => {
      setActive(false);
    }, 2000);
    return () => { 
      clearTimeout(timer); 
    }
  }

  return (
    <RequireAuth>
     { clientSecret !== "" && 
      <Container className="checkout-container mt-5 pt-4">
        <Row className="justify-content-center">
          <Col lg={9} className="checkout-title">
            Checkout
          </Col>
        </Row>
        <Row className="justify-content-center">
          <Col lg={9} className="checkout-form">
            <Elements options={options} stripe={stripePromise} key={clientSecret}>
              <CheckoutForm ordersInCart={ordersInCart} totalAmount={totalAmount}/>
            </Elements>
          </Col>
        </Row>
        <Row className="mt-5 pt-5 justify-content-center">
          <Col lg={9} className="credit-card-testing">
            <h5>Credit Cards for Testing</h5>
            <Table striped bordered hover >
              <thead className="">
                <tr className="">
                  <th width="">Brand</th>
                  <th width="">Number</th>
                  <th width="">Description</th>
                  <th width="">CVC</th>
                  <th width="">Date</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Visa</td>
                  <td>
                    <button className="text-nowrap" id="visa" onClick={(e)=>copyToClipboard(e, visaRef)} ref={visaRef}>
                      <span className="card-number">
                        <span className="visa">4242</span><span>4242</span><span>4242</span><span>4242</span>
                      </span>
                      <div>
                        {isActiveVisa ? <TiTick onClick={(e)=>{setCopiedActive(e, setIsActiveVisa)}} />
                        : <HiClipboard onClick={(e)=>{setCopiedActive(e, setIsActiveVisa)}} />}
                      </div>    
                    </button>
                  </td>
                  <td>Valid Visa card</td>
                  <td>Any 3 digits</td>
                  <td>Any future date</td>
                </tr>
                <tr>
                  <td>Mastercard</td>
                  <td>                    
                    <button id="master" onClick={(e)=>copyToClipboard(e, masterRef)} ref={masterRef}>
                      <span className="card-number">
                        <span className="master">5555</span><span>5555</span><span>5555</span><span>4444</span>
                      </span>
                      <div>
                      {isActiveMaster ? <TiTick onClick={(e)=>{setCopiedActive(e, setIsActiveMaster)}} />
                        : <HiClipboard onClick={(e)=>{setCopiedActive(e, setIsActiveMaster)}} />}
                      </div> 
                    </button>
                  </td>
                  <td>Valid Master card</td>
                  <td>Any 3 digits</td>
                  <td>Any future date</td>
                </tr>
                <tr>
                  <td>CVC check fails</td>
                  <td>
                    <button id="cvcFailed" onClick={(e)=>copyToClipboard(e, cvcFailedRef)} ref={cvcFailedRef}>
                      <span className="card-number">
                        <span className="master">4000</span><span>0000</span><span>0000</span><span>0101</span>
                      </span>
                      <div>
                      {isActiveCvcFailed ? <TiTick onClick={(e)=>{setCopiedActive(e, setIsActiveCvcFailed)}} />
                        : <HiClipboard onClick={(e)=>{setCopiedActive(e, setIsActiveCvcFailed)}} />}
                      </div> 
                    </button>
                  </td>
                  <td>If you provide a CVC number, the CVC check fails.</td>
                  <td>Any 3 digits</td>
                  <td>Any future date</td>
                </tr>
                <tr>
                  <td>Always blocked</td>
                  <td>
                    <button id="blocked" onClick={(e)=>copyToClipboard(e, blockedRef)} ref={blockedRef}>
                      <span className="card-number">
                        <span className="master">4100</span><span>0000</span><span>0000</span><span>0019</span>
                      </span>
                      <div>
                      {isActiveBlocked ? <TiTick onClick={(e)=>{setCopiedActive(e, setIsActiveBlocked)}} />
                        : <HiClipboard onClick={(e)=>{setCopiedActive(e, setIsActiveBlocked)}} />}
                      </div> 
                    </button>
                  </td>
                  <td>The charge has a risk level of “highest” Radar always blocks it.</td>
                  <td>Any 3 digits</td>
                  <td>Any future date</td>
                </tr>
              </tbody>
            </Table>
          </Col>
        </Row>
      </Container> 
    }
    </RequireAuth>
  );
}

export default Checkout;