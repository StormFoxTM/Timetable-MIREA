import a from "./Register.module.css";
import React from 'react';
import {NavLink, useNavigate} from "react-router-dom";
import {autorizationCreator} from '../../redux/AutorizationReduser'

const Authorization = (props) => {
    let adduser = () =>{
      let navigate = useNavigate();
      props.dispatch(autorizationCreator(newLogElem.current.value, newPassElem.current.value));
      navigate('/mypage')
  }
  let newLogElem = React.createRef();
  let newPassElem = React.createRef();
  return (
    <div className={a.login_block}>
    <div className={a.container_login}>
            
            <div className={a.form_reg}>
        <p className={a.text_form}>Логин </p><input ref={newLogElem} className={a.form_input} type='text'/>
        </div>
        <div className={a.form_reg}>
        <p className={a.text_form}>Пароль </p><input ref={newPassElem} className={a.form_input} type='password'/>
        </div>
              
              <div>
                <button onClick={adduser} className={a.activ_button}><p className={a.text_form}>Войти</p></button>
              </div>
                <div className={a.container_or}><p className={a.text_form}>или</p></div>
                <div><NavLink to='/register'>
                  <button className={a.registration_button}><p className={a.text_form}>Зарегистрироваться</p></button>
                  </NavLink></div>
        </div>
    </div>
  );
} 

export default Authorization;