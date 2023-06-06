import './TimeTable.css';
import React from 'react';
import {getTimeTable, getDateTable, getResponse} from './../../redux/ReduserTable'

const ChoiceMenu = (props) =>{
    let tableRef = React.createRef();
    let settable = () =>{
        props.dispatch(getTimeTable(tableRef.current.value));
    }
    let gettable = (t) =>{
      props.dispatch(getDateTable(t.target.value));
  }
  let gettable_but = () =>{
    props.dispatch(getResponse());
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
                {/* <div>
                    <p className='timetable_type_info'><b>Выберите день</b></p>
                    <input type="date" name="calendar"/>
                </div> */}
                <button onClick={() => {gettable_but()}}>Получить расписание</button>
                </div>
            </div>
        </div>
        )
}


const Tabletime = (props) =>{
    let table = () => {
      console.log(props.state.data_type)
      if (props.state.data_type === 'day'){
          return(
            <table className='timetable_table_day'>
              <tbody>
                {props.state.table.map(e => {
                  if (e.day !== null) {
                    return (
                      <td valign="top" className='timetable_elem'>
                        {e.day.map(day => (
                          <div className='timetable_elem_tr'>
                            <div className='timetable_elem_tr_in_day'>
                          {/* <tr className='timetable_elem_tr'> */}
                          <div className='num'><p>{day.subject_to_number}</p></div>
                          <p>{day.subject_title}</p>
                          <div className='ni'><div className='lec'><p>{day.name_lecturer}</p></div>
                          <p>{day.auditorium}</p></div>
                          {/* </tr> */}
                          </div>
                          </div>
                        ))}
                      </td>
                    );
                  }
                })}
              </tbody>
            </table>
        );
        } else{
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
                    <th className='th_head'>{day}</th>
                  ))}

                </tr>
              </thead>
              <tbody>
                {props.state.table.map(e => {
                  if (e.day !== null) {
                    return (
                      <td valign="top" className='timetable_elem'>
                        {e.day.map(day => (
                          <div className='timetable_elem_tr'>
                            <div className='timetable_elem_tr_in'>
                          {/* <tr className='timetable_elem_tr'> */}
                          <div className='num'><p>{day.subject_to_number}</p></div>
                          <p>{day.subject_title}</p>
                          <div className='ni'><div className='lec'><p>{day.name_lecturer}</p></div>
                          <p>{day.auditorium}</p></div>
                          {/* </tr> */}
                          </div>
                          </div>
                        ))}
                      </td>
                    );
                  }
                })}
              </tbody>
            </table>
        );
        }
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
                <Tabletime state={props.state}/>
            </div>
          );
    }
    return(
        <div className='timetable_main'>
            <ChoiceMenu dispatch={props.dispatch}/>
            <div className='timetable_body'>
            <h1>Расписание не получено, кликните на кнопку в меню справа, чтобы получить расписание вашей группы на неделю</h1>
        </div>
        </div>
      );
  
}

export default TimeTable;