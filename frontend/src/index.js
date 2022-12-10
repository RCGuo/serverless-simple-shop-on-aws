import React from 'react';
import ReactDOM from 'react-dom/client';
import { Amplify, Auth } from 'aws-amplify';
import { Authenticator } from '@aws-amplify/ui-react';
import { BrowserRouter } from "react-router-dom";
import reportWebVitals from './reportWebVitals';
import App from './App';
import config from "./config";
import 'bootstrap/dist/css/bootstrap.min.css';
import './index.css';
import ScrollToTop from './ScrollToTop';

Amplify.configure({
  Auth: {
    mandatorySignIn: true,
    region: config.cognito.REGION,
    userPoolId: config.cognito.USER_POOL_ID,
    userPoolWebClientId: config.cognito.APP_CLIENT_ID,
    signUpVerificationMethod: 'code',
  },
  API: {
    endpoints: [
      {
        name: "shopApi",
        endpoint: config.ecommerceApiGateway.API_URL,
        region: config.ecommerceApiGateway.REGION,
        custom_header: 
          async () => {
            try {
              const session = await Auth.currentSession();
              const token = session.getIdToken().getJwtToken();
              return { Authorization: `Bearer ${token}` };
            } catch (error) {
              // console.info('User session not present');
            }
          }
      },
      {
        name: "search",
        endpoint: config.opensearchApiGateway.API_URL,
        region: config.opensearchApiGateway.REGION,
      },
    ]
  }
});

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <Authenticator.Provider>
      <BrowserRouter>
        <ScrollToTop />
        <App />
      </BrowserRouter>
    </Authenticator.Provider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();