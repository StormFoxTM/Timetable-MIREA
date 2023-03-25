import {NavLink} from 'react-router-dom'
import './header.css'

const Header = () => {
    return(
    <header>
        <header className="header-left">
            <NavLink to='/' className="main_header">Главная</NavLink>
            <NavLink to='/timeTable' className="timetable_header">Расписание</NavLink>
        </header>
        <header className="header-right">
            <div className="dropdown">
                <NavLink className="settings">Настройки</NavLink>
                <div id="popup" className="dropdown-content">
                    <NavLink to='/change_themes' className="background" href="/change_themes">Сменить тему</NavLink>
                </div>
            </div>
            <NavLink to='/authorization' className="login">Войти</NavLink>
        </header>
    </header>
    );
}

export default Header;