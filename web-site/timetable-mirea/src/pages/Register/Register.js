import './Register.css';

function Register() {
  return (
    <div className='login_block'>
    <div className='container_login'>
            <div className='login_main_text'>Регистрация</div>
            <form method="POST">
                <p>Введите логин: <input className='form-input' /></p>
                <p>Введите email: <input  type='email' className='form-input' /></p>
                <p>Введите пароль: <input type='password' className='form-input' /></p>
                <p>Повторите пароль: <input type='password' className='form-input' /></p>
                <button type="submit" id="enter">Регистрация</button>
            </form>
        </div>
    </div>
  );
}

export default Register;