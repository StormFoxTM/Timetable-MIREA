import a from "./Register.module.css";
import {NavLink} from "react-router-dom"
import Elem from './Elem';

const Authorization = () => {
  return (
    <div className={a.login_block}>
    <div className={a.container_login}>
            <form>
              <Elem name='Логин' type='text'/>
              <Elem name='Пароль' type='password'/>
              <div>
                <button className={a.activ_button}><p className={a.text_form}>Войти</p></button>
              </div>
                <div className={a.container_or}><p className={a.text_form}>или</p></div>
                <div><NavLink to='/register'>
                  <button className={a.registration_button}><p className={a.text_form}>Зарегистрироваться</p></button>
                  </NavLink></div>
            </form>
        </div>
    </div>
  );
} 

export default Authorization;