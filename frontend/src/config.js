const Config = {
  ecommerceApiGateway: {
    REGION: process.env.REACT_APP_REGION,
    API_URL: process.env.REACT_APP_ECOMMERCE_API_URL,
  },
  opensearchApiGateway: {
    REGION: process.env.REACT_APP_REGION,
    API_URL: process.env.REACT_APP_OPENSEARCH_API_URL,
  },
  cognito: {
    REGION: process.env.REACT_APP_REGION,
    USER_POOL_ID: process.env.REACT_APP_USER_POOL_ID,
    APP_CLIENT_ID: process.env.REACT_APP_CLIENT_ID,
  }
};

export default Config;