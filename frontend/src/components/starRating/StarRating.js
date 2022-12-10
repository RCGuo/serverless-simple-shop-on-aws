import React from 'react';
import { BsStarFill, BsStarHalf, BsStar } from 'react-icons/bs';
import "./starRating.css";


const StarFill = () => {
  return (
    <BsStarFill size={12} color={"#ffd700"} />
  );
}

const StarFillHalf = () => {
  return (
    <BsStarHalf size={12} color={"#ffd700"} />
  );
}

const StarEmpty = () => {
  return (
    <BsStar size={12} />
  );
}

const calculateStar = (base, rating) => {
  return (
    rating >= base ? <StarFill /> : rating >= base - 0.5 ? <StarFillHalf /> : <StarEmpty />
  ); 
}

const StarRating = (props) => {
  return (
    <span>
      {calculateStar(1, props.rating)}
      {calculateStar(2, props.rating)}
      {calculateStar(3, props.rating)}
      {calculateStar(4, props.rating)}
      {calculateStar(5, props.rating)}
    </span>
  );
}

export default StarRating;