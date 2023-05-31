
import set from './Setting.module.css';
import { useTheme } from '../../hooks/UseTheme'

function Setting() {
  const { setTheme } = useTheme()
    const LightTheme = () => {
        setTheme('light')
    }
    const DarkTheme = () => {
        setTheme('dark')
    }
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
        <button onClick={() => {DarkTheme()}}>Темная тема</button>
        <button onClick={() => {LightTheme()}}>Светлая тема</button>
      </div>
    </div>
  </div>
      );
    }
    
export default Setting;