import './TimeTable.css';
import axios from 'axios';
import React, { useEffect, useState } from 'react';

let ChoiseElem = (props) =>{
    return(
        <div>
            <label><input type="radio" id={props.id_el}/>{props.name}</label>
        </div>
    )
}

const TimeTable = () => {
    // let gettable = () => {
    //     axios.get("localhost:9888?group=ИКБО-02-20&week=1&day=2 ").then(response=>{
    //         debugger;
    //         console.log(response)
    //     })
    // }
  return (
        <div className='timetable_main'>
        <div className='timetable_panel'>
            <div className='timetable_panel_wrapper'>
                <form>
                <p className='timetable_type_info'><b>Введите номер группы, ФИО преподавателя или номер аудитории</b></p>
                <input type="text" className=''/>
                <button type="submit" className='change_type_info' name="type_of_info" value="day">Получить расписание</button>
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
                </form> 
            </div>
        </div>
        <div className='timetable_body'>
            <div className='timetable_body_wrapper'>
                <div>
                    <p className='timetable_body_header'>Расписание группы: </p>
                    {/* <button onClick={ gettable }></button> */}
                        {/* <table className='timetable_table'>
                                <tr className='timetable_column'>
                                    {getTimetable}
                                </tr>
                            
                        </table> */}
                </div>
            </div>
        </div>
    </div>
  );
}

export default TimeTable;