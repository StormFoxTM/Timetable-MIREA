import {NavLink} from 'react-router-dom'
import set from './Setting.module.css';

function Setting() {
  return (
  <div className={set.set_main}>
    <div className={set.container}>
      <div className={set.table}>
      <div className={set.tabl}><h1 className={set.text}>Настройки отображения расписания</h1></div>
      <div className={set.group}><p className={set.text}>Номер группы</p>
      <input></input></div>
      <div className={set.tabl}><button>Изменить номер группы</button></div>
      </div>
      <div className={set.tem}>
        <p className={set.text}>Настройки темы</p>
        <button>Сменить тему</button>
      </div>
    </div>
  </div>
      );
    }
    
export default Setting;