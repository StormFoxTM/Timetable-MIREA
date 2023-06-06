import './TimeTable.css';
import React from 'react';
import {getTimeTable, getDateTable} from './../../redux/ReduserTable'

const ChoiceMenu = (props) =>{
    let tableRef = React.createRef();
    let settable = () =>{
        props.dispatch(getTimeTable(tableRef.current.value));
    }
    let gettable = (t) =>{
      props.dispatch(getDateTable(t.target.value));
  }
    return (
        <div className='timetable_panel'>
            <div className='timetable_panel_wrapper'>
                <p className='timetable_type_info'><b>Введите номер группы, ФИО преподавателя или номер аудитории</b></p>
                <div>
                <input ref={tableRef} type="text" className='input_group'/>
                <button onClick={() => {settable()}}>Получить расписание</button>
                </div>
                <div>
                    <p className='timetable_type_info'><b>Быстрая настройка:</b></p>
                    <div>
                      <label><input type="radio" name="choice" value="today" onChange={gettable} id='Choice1'/>Сегодня</label> 
                    </div>
                    <div>
                      <label><input type="radio" name="choice" value="tomorrow" onChange={gettable} id='Choice2'/>Завтра</label>
                    </div> 
                    <div>
                      <label><input type="radio" name="choice" value="week" onChange={gettable} id='Choice3'/>Эта неделя</label>
                    </div>
                    <div>
                      <label><input type="radio" name="choice" value="next_week" onChange={gettable} id='Choice4'/>Следующая неделя</label>
                    </div> 
                <div>
                    <p className='timetable_type_info'><b>Выберите день</b></p>
                    <input type="date" name="calendar"/>
                </div>
                <button onClick={() => {gettable()}}>Получить расписание</button>
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
            <table className='timetable_table'>
              <thead>
                <tr className='timetable_column'>
                  {/* <th>Номер предмета</th>
                  <th>Название</th>
                  <th>Преподаватель</th>
                  <th>Аудитория</th> */}
                  {/* <th>&nbsp;</th> */}
                  {daysOfWeek.map(day => (
                    <th>{day}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {props.state.map(e => {
                  if (e.day !== null) {
                    return (
                      <td className='timetable_elem'>
                        {e.day.map(day => (
                          <tr>
                            <p>{day.subject_to_number}</p>
                            <p>{day.subject_title}</p>
                            <p>{day.name_lecturer}</p>
                            <p>{day.auditorium}</p>
                          </tr>
                        ))}
                      </td>
                    );
                  }
                })}
              </tbody>
            </table>
        );
      };
    return(
      <div className='timetable_body'>
          <div className='timetable_body_wrapper'>
            {table()}
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