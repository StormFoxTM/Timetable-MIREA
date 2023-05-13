import {NavLink} from 'react-router-dom'
import './header.css'

const Heder_elem = (props) => {
    return(
        <NavLink to={props.k_ref}
        className={props.name_class}>
            {props.name}
        </NavLink>
    );
}

const Header = () => {
    return(
    <header className="header">
        <div className="header_left">
            <Heder_elem k_ref='/' name='Расписание' name_class='main_header' />
        </div>
        <div className="header_right">
            <div className="dropdown">
                <Heder_elem k_ref='/setting' name='Настройки' name_class='settings' />
            </div>
            <div id="popup" className="dropdown_content">
                <Heder_elem k_ref='/authorization' name='Войти' name_class='login' />
            </div>
            
        </div>
    </header>
    );
}

export default Header;