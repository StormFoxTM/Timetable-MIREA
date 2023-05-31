import './TimeTable.css';
import React from 'react';
import {getTimeTable} from './../../redux/ReduserTable'

const ChoiseElem = (props) =>{
    return(
        <div>
            <label><input type="radio" id={props.id_el}/>{props.name}</label>
        </div>
    )
}

const ChoiceMenu = (props) =>{
    let tableRef = React.createRef();
    let settable = () =>{
        props.dispatch(getTimeTable(tableRef.current.value));
    }
    return (
        <div className='timetable_panel'>
            <div className='timetable_panel_wrapper'>
                <p className='timetable_type_info'><b>Введите номер группы, ФИО преподавателя или номер аудитории</b></p>
                <div>
                <input ref={tableRef} type="text" className='input_group'/>
                <button onClick={() => {settable()}}>Получить расписание</button>
                </div>
                    <p className='timetable_type_info'><b>Быстрая настройка:</b></p>
                    <ChoiseElem id_el='Choice1' name='Сегодня'/>
                    <ChoiseElem id_el='Choice2' name='Завтра'/>
                    <ChoiseElem id_el='Choice3' name='Эта неделя'/>
                    <ChoiseElem id_el='Choice4' name='Следующая неделя'/>
                    <ChoiseElem id_el='Choice5' name='Месяц'/>
                <div>
                    <p className='timetable_type_info'><b>Выберите день</b></p>
                    <input type="date" name="calendar"/>
                </div>
            </div>
        </div>
        )
}


const Tabletime = (props) =>{
    let table = () => {
        // Дни недели для заголовков
        const daysOfWeek = ['Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота'];
      
        return (
          <div className='timetable'>
            <table>
              <thead>
                <tr>
                  {/* <th>Номер предмета</th>
                  <th>Название</th>
                  <th>Преподаватель</th>
                  <th>Аудитория</th> */}
                  <th>&nbsp;</th>
                  {daysOfWeek.map(day => (
                    <th>{day}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {props.state.map(e => {
                  if (e.day !== null) {
                    return (
                      <tr>
                        {e.day.map(day => (
                          <th>
                            <p>{day.subject_to_number}</p>
                            <p>{day.subject_title}</p>
                            <p>{day.name_lecturer}</p>
                            <p>{day.auditorium}</p>
                          </th>
                        ))}
                      </tr>
                    );
                  }
                })}
              </tbody>
            </table>
          </div>
        );
      };
    return(
        <div className='timetable_body'>
            <div className='timetable_body_wrapper'>
                <div>
                        {table()}
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