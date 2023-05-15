import Elem from './Elem';
import r from './Register.module.css';
import {NavLink} from "react-router-dom"


const Register = () => {
    return (
        <div className={r.login_block}>
        <div className={r.container_register}>
                <form>
                    <Elem name='Адрес электронной почты' type='email'/>
                    <Elem name='Логин' type='text'/>
                    <Elem name='Пароль' type='password'/>
                    <Elem name='Подтверждение пароля' type='password'/>
                    <div>
                        <button className={r.activ_button}><p className={r.text_form}>Зарегистрироваться</p></button>
                    </div>
                    <div className={r.container_or}><p className={r.text_form}>или</p></div>
                    <div><NavLink to='/authorization'>
                  <button className={r.registration_button}><p className={r.text_form}>Войти</p></button>
                  </NavLink></div>
                </form>
            </div>
        </div>
    );
}

export default Register;