import React from 'react';
import set from './../Setting/Setting.module.css';
import { useHistory } from "react-router-dom";
import { logoutCreator } from '../../redux/AutorizationReduser'
import Authorization from './Authorization';

function My_page(props) {
    
    let Logout = () =>{
        props.dispatch(logoutCreator());
        
    }
    return (
    <div className={set.set_main}>
      <div className={set.container}>
        <div className={set.table}>
            <div className={set.tabl}><h1 className={set.text}>Личный кабинет</h1></div>
            <div className={set.group}><p className={set.text}>Логин: </p>
                <p className={set.text}> Логин.props</p></div>
            <div className={set.group}><p className={set.text}>Адрес электронной почты: </p>
                <p className={set.text}> Адрес.props</p></div>
        </div>
        <div className={set.table}>
        <div className={set.tabl}><h1 className={set.text}>Смена пароля</h1></div>
        <div className={set.group}><p className={set.text}>Текущий пароль</p>
        <input></input></div>
        <div className={set.group}><p className={set.text}>Новый пароль</p>
        <input></input></div>
        <div className={set.tabl}><button>Изменить пароль</button></div>
        </div>
        <div className={set.account}>
        <div className={set.tabl}><button onClick={Logout}>Выход из аккаунта</button></div>
        {/* <div className={set.tabl}><button>Удалить аккаунт</button></div> */}
        </div>
      </div>
    </div>
    );
}

export default My_page;