import { useEffect, useState } from "react";
import { Container, Row, Col, Button } from "react-bootstrap";
import { Link } from 'react-router-dom';


const CheckoutCompleted = () => {
  const [message, setMessage] = useState("");

  useEffect(() => {
    const redirect_status = new URLSearchParams(window.location.search).get(
      "redirect_status"
    );
    
    if (redirect_status === 'succeeded') {
      setMessage("Thank you for your order!");
    } else if (redirect_status === 'failed') {
      setMessage("Something went wrong");
    }
  }, [])

  return (
    <Container className="checkout-completed-container mt-5 pt-4">
      <Row className="text-center">
        <Col>
          <h1>{message}</h1>
          <Button className="mt-3" as={Link} to="/past-orders" variant="warning">View past order</Button>
        </Col>
      </Row>
    </Container>
    
  );
}

export default CheckoutCompleted;