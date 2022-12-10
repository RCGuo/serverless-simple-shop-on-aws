import React from 'react';
import { Routes, Route} from "react-router-dom";
import Authentication from "./components/user/Authentication";
import Home from "./pages/Home";
import CategoryView from "./pages/CategoryView";
import TopicView from "./pages/TopicView";
import PastOrders from "./pages/PastOrders";
import Landing from "./pages/Landing";
import Cart from "./pages/Cart";
import Checkout from "./pages/Checkout";
import CheckoutCompleted from "./pages/CheckouteCompleted";
import NotFound from "./pages/NotFound";
import SearchView from './pages/SearchView';
import BestSellersView from './pages/BestSellersView';

const RoutePaths = () => {
  return (
    <Routes>
      <Route path="/auth"               element={<Authentication />} /> 
      <Route path="/home"               element={<Home />} />
      <Route path="/"                   element={<Landing />} />
      <Route path="/past-orders"        element={<PastOrders />} />
      <Route path="/favorites"          element={<PastOrders />} />
      <Route path="/account"            element={<PastOrders />} />
      <Route path="/cart"               element={<Cart />} />
      <Route path="/checkout"           element={<Checkout />} />
      <Route path="/checkout-complete"  element={<CheckoutCompleted />} />
      <Route path="/category/:category" element={<CategoryView />} />
      <Route path="/topic/:topic"       element={<TopicView />} />
      <Route path="/search"             element={<SearchView />} />
      <Route path="/best-sellers"       element={<BestSellersView />} />
      <Route path="*"                   element={<NotFound />} />
    </Routes>
  );
}

export default RoutePaths;