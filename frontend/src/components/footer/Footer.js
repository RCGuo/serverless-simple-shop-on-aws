
import { Container, Row, Col } from 'react-bootstrap';
import './footer.css';

import TownDrawing from '../../resources/images/footer/line_drawing_of_town_illustration.jpg';

const Footer = () => {
  return (
    <Container fluid className='footer-container mt-3 mx-0 px-0 text-center'>
      <Row>
        <Col>
          <img src={TownDrawing} width="400px" alt="Town Drawing" />
        </Col>
      </Row>
      <Row className='bottom-bar mx-0 text-center'>
        <Col >
          <div>© 2022 SimpleShop</div>
          <div>Made by RCGuo</div>
        </Col>
      </Row>
    </Container>
  );
}

export default Footer;