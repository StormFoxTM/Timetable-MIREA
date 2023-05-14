import './TimeTable.css';
import {getTimeTable} from './../../redux/ReduserTable'

let ChoiseElem = (props) =>{
    return(
        <div>
            <label><input type="radio" id={props.id_el}/>{props.name}</label>
        </div>
    )
}

let ChoiceMenu = (props) =>{
    let settable = () =>{
        props.dispatch(getTimeTable());
    }
    return (
        <div className='timetable_panel'>
            <div className='timetable_panel_wrapper'>
                <p className='timetable_type_info'><b>Введите номер группы, ФИО преподавателя или номер аудитории</b></p>
                <div>
                <input type="text" className='input_group'/>
                <button onClick={() => {settable()}}>Получить расписание</button>
                </div>
                    <p className='timetable_type_info'><b>Быстрая настройка:</b></p>
                    <ChoiseElem id_el='Choice1' name='Сегодня'/>
                    <ChoiseElem id_el='Choice2' name='Завтра'/>
                    <ChoiseElem id_el='Choice3' name='Эта неделя'/>
                    <ChoiseElem id_el='Choice4' name='Следующая неделя'/>
                    <ChoiseElem id_el='Choice5' name='Месяц'/>
                <div>
                    <p className='timetable_type_info'><b>Выберите период</b></p>
                    <input type="date" name="calendar"/>
                </div>
            </div>
        </div>
        )
}

let Tabletime = (props) =>{
    let table = () => {
        console.log(props)
    }
    return(
        <div className='timetable_body'>
            <div className='timetable_body_wrapper'>
                <div>
                    <p className='timetable_body_header'>Расписание группы: </p>
                        {table}
                </div>
            </div>
        </div>
    )
}

const TimeTable = (props) => {
    if (props.state.getTable){
        return(
            <div className='timetable_main'>
                <ChoiceMenu dispatch={props.dispatch}/>
                <Tabletime state={props.state.table}/>
            </div>
          );
    }
    return(
        <div className='timetable_main'>
            <ChoiceMenu dispatch={props.dispatch}/>
            <div className='timetable_body'>
            <p>Расписание не получено, кликните на кнопку в меню справа, чтобы получить расписание вашей группы на неделю</p>
        </div>
        </div>
      );
  
}

export default TimeTable;