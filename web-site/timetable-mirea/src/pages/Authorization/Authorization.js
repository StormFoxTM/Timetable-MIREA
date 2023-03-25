import './Authorization.css';
import {NavLink} from 'react-router-dom'

function Authorization() {
  return (
    <div className='login_block'>
    <div className='container_login'>
            <div className='login_main_text'>Вход в аккаунт</div>
            <form method="POST">
                <p>Введите логин или email: <input className='form-input' /></p>
                <p>Введите пароль: <input type='password' className='form-input' /></p>
                <button type="submit" id="enter">Вход</button>
                <NavLink to='/register' className="registration_button">Зарегистрироваться</NavLink>
            </form>
        </div>
    </div>
  );
}

export default Authorization;