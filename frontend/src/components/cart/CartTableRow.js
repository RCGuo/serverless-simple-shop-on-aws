import { API } from "aws-amplify";
import { useEffect, useState } from "react";
import { Dropdown, Image, Spinner } from "react-bootstrap";
import { Link } from "react-router-dom";
import './cart.css';

const CartTableRow = (props) => {
  const [productInfo, setProductInfo] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isRemoving, setIsRemoving] = useState(false);
  const [quantity, setQuantity] = useState(props.order.quantity);
  const [qtyLoading, setQtyLoading] = useState(false);

  useEffect(() => {
    const getProductByID = async () => {
      try {
        await API.get("shopApi", "/product", {
          queryStringParameters: {
            productId: props.order.productId
          }
        }).then((data) => {
            setProductInfo(data);
            setIsLoading(false);
          });
      } catch (e) {
        console.log(e);
      }
    }
    getProductByID(props.order);
  }, [props.order]);

  const onRemove = async () => {
    setIsRemoving(true);
    await API.del("shopApi", "/cart", {
      body: {
        productId: props.order.productId,
      },
    });
    props.getOrderTotal();
  }

  const onQuantityUpdate = async (newQty, oldQty) => {
    if ( parseInt(newQty, 10) !== parseInt(oldQty, 10) ) {
      setQtyLoading(true);
      setQuantity(parseInt(newQty, 10));
      await API.put("shopApi", "/cart", {
        body: {
          productId: props.order.productId,
          quantity: parseInt(newQty, 10),
        },
      });
      props.getOrderTotal();
      setQtyLoading(false);
    }
  }

  if (isLoading) {
    return(
      <tr>
        <td colSpan="4">
          <Spinner animation="border" variant="secondary" size="lg" />
        </td>
      </tr>
    );
  }
  return (
    <tr className="cart-row-product-detail" aria-hidden="true">
      <td></td>
      <td>
        <div style={{width: "60px"}}>
          <Image src={productInfo.imageFile} fluid={true}/>
        </div>
      </td>
      <td>
        <div id="product-name">{productInfo.name}</div>
        <div id="product-stock">In Stock</div>
        <div id="product-shipping">FREE Shipping </div>
        <div className="cart-row-product-delete">
          {isRemoving 
            ? <Spinner animation="border" variant="secondary" size="sm" />
            : <Link onClick={() => onRemove()}>Delete</Link>}
        </div>
      </td>
      
      <td id="product-price">
        {props.order.price.toLocaleString('en-US', { style: 'currency', currency: 'usd' })}
      </td>
      <td id="product-quantity">
         <Dropdown onSelect={(newQty) => onQuantityUpdate(newQty, quantity)}>
          <Dropdown.Toggle variant="secondary" id="dropdown-basic" disabled={qtyLoading}>
            Qty: {quantity}
          </Dropdown.Toggle>
          <Dropdown.Menu className="w-100">
            <Dropdown.Item eventKey="1">1</Dropdown.Item>
            <Dropdown.Item eventKey="2">2</Dropdown.Item>
            <Dropdown.Item eventKey="3">3</Dropdown.Item>
            <Dropdown.Item eventKey="4">4</Dropdown.Item>
            <Dropdown.Item eventKey="5">5</Dropdown.Item>
            <Dropdown.Item eventKey="6">6</Dropdown.Item>
            <Dropdown.Item eventKey="7">7</Dropdown.Item>
            <Dropdown.Item eventKey="8">8</Dropdown.Item>
            <Dropdown.Item eventKey="9">9</Dropdown.Item>
            <Dropdown.Item eventKey="10">10</Dropdown.Item>
          </Dropdown.Menu>
        </Dropdown>
        </td>
    </tr>   
  );
}

export default CartTableRow;

