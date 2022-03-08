import React from "react";
import logo from '../../../assets/icons/app_white_24dp.svg';
import './Header.scss';

export function Header() {

    return (
        <div className="header">
            <img src={logo} alt="logo" className="logo"></img>
            <span>Discord Sound Board</span>
        </div>
    )
}