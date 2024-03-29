import React from 'react';
import r from './Register.module.css';
import {NavLink, useNavigate} from "react-router-dom";
import {registrationCreator} from '../../redux/AutorizationReduser'

const Register = (props) => {
    let newLogElem = React.createRef();
    let newPassCheckElem = React.createRef();
    let newPassElem = React.createRef();
    let RefGroup = React.createRef();
    let navigate = useNavigate();

    let Adduser = () =>{
        if (newPassElem.current.value===newPassCheckElem.current.value){
            props.dispatch(registrationCreator(newLogElem.current.value, newPassElem.current.value, RefGroup.current.value));
            navigate('/authorization');
        }
        else{
            console.log("error")
        }
    }
    
    return (
        <div className={r.login_block}>
        <div className={r.container_register}>
                
                <div className={r.form_reg}>
        <p className={r.text_form}>Адрес электронной почты </p><input className={r.form_input} type='email'/>
        </div>
        <div className={r.form_reg}>
        <p className={r.text_form}>Логин </p><input ref={newLogElem} className={r.form_input} type='text'/>
        </div>
        <div className={r.form_reg}>
        <p className={r.text_form}>Группа </p><input ref={RefGroup} className={r.form_input} type='text'/>
        </div>
        <div className={r.form_reg}>
        <p className={r.text_form}>Пароль </p><input ref={newPassElem} className={r.form_input} type='password'/>
        </div>
        <div className={r.form_reg}>
        <p className={r.text_form}>Подтверждение пароля </p><input ref={newPassCheckElem} className={r.form_input} type='password'/>
        </div>
                    <div>
                        <button onClick={Adduser} className={r.activ_button}><p className={r.text_form}>Зарегистрироваться</p></button>
                    </div>
                    <br></br>
                    <div><NavLink to='/authorization'>
                  <button className={r.registration_button}><p className={r.text_form}>Войти</p></button>
                  </NavLink></div>
                
            </div>
        </div>
    );
}

export default Register;