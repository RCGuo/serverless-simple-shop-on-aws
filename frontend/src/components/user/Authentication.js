import React, { useEffect } from "react";
import { Authenticator, useAuthenticator } from '@aws-amplify/ui-react';
import { useNavigate, useLocation, Navigate } from "react-router";
import { Col, Row, Container } from "react-bootstrap";
import "@aws-amplify/ui-react/styles.css";

export const RequireAuth = ({ children }) => {
  const location = useLocation();
  const { route } = useAuthenticator((context) => [context.route]);

  if (route !== 'authenticated') {
    return <Navigate to="/auth" state={{ from: location }} replace />;
  }
  return children;
}

const formFieldsConfig = {
  signUp: {
    username: {
      placeholder: 'Email Address',
      isRequired: true,
      label: 'Email Address*',
      order: 1
    },
    password: {
      placeholder: 'Password',
      isRequired: true,
      label: 'New Password*',
      order: 2
    },
    confirm_password: {
      placeholder: 'Confirm Password',
      isRequired: true,
      label: 'Confirm Password*',
      order: 3
    },
    nickname: {
      placeholder: 'Nickname',
      isRequired: true,
      order: 4
    },
  },
  signIn: {
    username: {
      placeholder: 'Email Address',
      isRequired: true,
      label: 'Email Address*',
      order: 1
    }, 
  }
}

const Authentication = () => {
  const { route } = useAuthenticator((context) => [context.route]);
  const location = useLocation();
  const navigate = useNavigate();
  let from = location.state?.from?.pathname || '/home';

  useEffect(() => {
    if (route === 'authenticated') {
      navigate(from, { replace: true });
    }
  }, [route, navigate, from]);

  return (
    <Container className="mt-5 pt-5">
      <Row>
        <Col>
          <Authenticator formFields={formFieldsConfig} usernameAttributes='email'>
            {/* {() => ( setLoggedIn(true) )} */}
          </Authenticator>
        </Col>
      </Row>
    </Container>
  );
}

export default Authentication;