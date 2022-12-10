import { API } from "aws-amplify";
import { useAuthenticator } from '@aws-amplify/ui-react';
import { useEffect, useState } from "react";
import { OverlayTrigger, Tooltip } from "react-bootstrap";
import { FaRegKissWinkHeart, FaRegMeh } from 'react-icons/fa';
import { toast } from 'react-toastify';
import './button.css';

const FavoriteIcon = (props) => {
  const [favoriteState, setFavoriteState] = useState(false);
  const { authStatus } = useAuthenticator(context => [context.authStatus])

  useEffect(() => {
    setFavoriteState(!!props.favorite);
  }, [props.favorite])

  const notify = () => {
    toast.warn('Please login before adding a product to favorite.', {
      position: "top-center",
      autoClose: 4000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
      theme: "light",
    });
  }

  const onChangeFavorite = async () => {
    if (authStatus !== 'authenticated') {
      notify();
      return;
    }

    setFavoriteState(!favoriteState)
    await API.post("shopApi", "/product/favorites", {
      body: {
        productId: props.productId,
        favorite: !favoriteState
      },
    });
  }

  const renderTooltip = (props) => (
    <Tooltip {...props}>Click to sign in and save</Tooltip>
  );

  return (
    authStatus !== 'authenticated' ? 
      <OverlayTrigger 
        placement="top"
        delay={{ show: 250, hide: 200 }}
        overlay={renderTooltip}>
        <span>{favoriteState === true 
          ? <FaRegKissWinkHeart id="kiss-face" onClick={onChangeFavorite} /> 
          : <FaRegMeh id="meh-face" onClick={onChangeFavorite} />}</span>
      </OverlayTrigger>
    : <span>{favoriteState === true 
        ? <FaRegKissWinkHeart id="kiss-face" onClick={onChangeFavorite} /> 
        : <FaRegMeh id="meh-face" onClick={onChangeFavorite} />}</span>
  );
}

export default FavoriteIcon;