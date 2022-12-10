import { useState } from 'react';
import { useStripe, useElements, PaymentElement} from '@stripe/react-stripe-js';
import { Button, Form, Spinner } from 'react-bootstrap';

const CheckoutForm = (props) => {
  const stripe = useStripe();
  const elements = useElements();
  const [message, setMessage] = useState(null);
  const [isLoading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    if (!stripe || !elements) {
      setLoading(false);
      return;
    }

    const { error } = await stripe.confirmPayment({
      elements,
      confirmParams: {
        return_url: `${window.location.origin}/checkout-complete`,
      }
    })
    
    if (error.type === "card_error" || error.type === "validation_error") {
      setMessage(error.message);
    } else {
      setMessage("An unexpected error occurred.");
    }
    setLoading(false);
  };

  return (
    <Form id="payment-form" onSubmit={handleSubmit}>
      <PaymentElement id="payment-element" />
      <Button className='mt-3' disabled={isLoading || !stripe || !elements} type="submit" id="submit" variant="warning">
        <span id="button-text">
          {isLoading 
            ? <Spinner animation="border" variant="secondary" size="sm" /> 
            : `Pay now (${props.totalAmount.toLocaleString('en-US', { style: 'currency', currency: 'usd' })})`}
        </span>
      </Button>
      {message && <div id="payment-message">{message}</div>}
    </Form>
  );
};

export default CheckoutForm;