import { API, Auth } from "aws-amplify";
import { useEffect, useState, useCallback } from "react";
import { useAuthenticator } from '@aws-amplify/ui-react';
import { RequireAuth } from "../components/user/Authentication";
import { Card, Col, Container, Row, Table, Tab, Tabs, Spinner } from "react-bootstrap";
import { useLocation, useNavigate } from "react-router-dom";
import { FavoriteTableRow, PastOrderTableRow } from "./../components/pastOrders/TableRow";
import { FaRegKissWinkHeart } from "react-icons/fa";
import { AiOutlineHistory } from "react-icons/ai";
import { VscAccount } from 'react-icons/vsc';
import { AiOutlineMail, AiOutlineSmile }from 'react-icons/ai';
import './../App.css';

import avatar from './../resources/images/common/avatar.jpg';

const PastOrders = () => {
  const location = useLocation();
  const [userInfo, setUserInfo] = useState({email:"", nickname:""});
  const [isOrderLoading, setOrderLoading] = useState(true);
  const [isFavorLoading, setFavorLoading] = useState(true);
  const [pastOrders, setPastOrders] = useState([]);
  const [favorites, setFavorites] = useState([]);
  const [currentTab, setCurrentTab] = useState('/past-order');
  const { authStatus } = useAuthenticator(context => [context.authStatus]);
  const navigate = useNavigate();

  const getFavorites = useCallback(() => {
    const getFavoritesData = async () => {
      try {
        await API.get("shopApi", "/product/favorites")
        .then((resp) => {
          API.post("shopApi", '/product/batch-fetch', {
            body: resp,
          }).then((data) => {
            setFavorites(data);
            setFavorLoading(false);
          });
        });
      } catch (e) {
        console.log(e);
      }
    }
    getFavoritesData();
  }, [])

  useEffect(() => {
    setCurrentTab(location.pathname);
    const getUserInfo = async () => {
      await Auth.currentAuthenticatedUser().then(
        data => {
          setUserInfo({
            email: data.attributes.email,
            nickname: data.attributes.nickname,
          });
        }
      );
    }
    const getOrderList = async () => {
      try {
        await API.get("shopApi", "/order/past-orders")
        .then((data) => {
          setPastOrders(data);
          setOrderLoading(false);
        });
        
      } catch(e) {
        console.log(e);
      }
    }

    if ( authStatus === 'authenticated' ) {
      if(currentTab === "/past-orders"){
        setOrderLoading(true);
        getOrderList();
      } else if (currentTab === "/favorites") {
        setFavorLoading(true);
        getFavorites();
      } else if (currentTab === "/account") {
        getUserInfo();
      }
    }
  },[currentTab, location.pathname, authStatus, getFavorites]);

  const handleSelect = (e) => {
    setCurrentTab(e);
    navigate(e);
  }

  return (
    <>
      <RequireAuth>
        <Container className="past-order-container mt-5 pt-4 px-5">
          <Row className="justify-content-center">
            <Col lg={10}>
              <Card>
                <Card.Body>
                  <Tabs defaultActiveKey="/past-orders" activeKey={currentTab} onSelect={handleSelect}>
                    <Tab eventKey="/past-orders" title={<><AiOutlineHistory />{'  '}Past Order</>}>
                      <Table bordered hover className="mt-3">
                        <thead className="grey-background">
                          <tr className="text-center">
                            <th width="25%">Order date</th>
                            <th width="25%">Order ID</th>
                            <th width="30%">Payment method</th>
                            <th width="20%">Total price</th>
                          </tr>
                        </thead>
                        <tbody>
                          {isOrderLoading ?
                            <tr>
                              <td colSpan="4"><Spinner animation="border" variant="secondary" size="lg" /></td>
                            </tr>
                            : pastOrders.length === 0 ?       
                              <tr className="text-center">
                                <td colSpan="4">Empty</td>
                              </tr> 
                              : pastOrders.map((data, index) => 
                                  <PastOrderTableRow order={data} key={index} />)
                          }
                          <tr >
                            <td colSpan="4" className="text-end">Total orders: {pastOrders.length === 0 ? 0 : pastOrders.length}</td>
                          </tr>
                        </tbody>
                      </Table>
                    </Tab>
                    <Tab eventKey="/favorites" title={<><FaRegKissWinkHeart />{'  '}My favorites</>} className="favorite-tab">
                      <Table className="mt-3" bordered hover>
                        <tbody>
                          {isFavorLoading ?
                            <tr>
                              <td colSpan="4"><Spinner animation="border" variant="secondary" size="lg" /></td>
                            </tr>
                            : favorites.length === 0 ?
                              <tr className="text-center">
                                <td colSpan="4">Empty</td>
                              </tr> 
                              : favorites.map((data, index) => 
                                  <FavoriteTableRow product={data} getFavorites={getFavorites} mark="for_favorite" key={data.productId} />)
                          }
                        </tbody>
                      </Table>
                    </Tab>
                    <Tab eventKey="/account" title={<><VscAccount />{'  '}Account</>} className="account-tab">
                      <Container className="mt-3 px-4">
                        <Row>
                          <Col><img src={avatar} style={{width:"10%"}} alt="avatar" /></Col>
                        </Row>
                        <Row className="mt-3">
                          <Col>
                            <p><AiOutlineSmile />{'  '}{userInfo.nickname}</p>
                            <p><AiOutlineMail />{'  '}{userInfo.email}</p>
                          </Col>
                        </Row>
                      </Container>
                    </Tab>
                  </Tabs>
                </Card.Body>
              </Card>  
            </Col>
          </Row>
        </Container>
      </RequireAuth>
    </>
  );
};

export default PastOrders;