import './TimeTable.css';
import axios from 'axios';


let ChoiseElem = (props) =>{
    return(
        <div>
            <label><input type="radio" id={props.id_el}/>{props.name}</label>
        </div>
    )
}

let ChoiceMenu = () =>{
    return (
        
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
        )
}

let Tabletime = () =>{
    return(
        <div className='timetable_body'>
            <div className='timetable_body_wrapper'>
                <div>
                    <p className='timetable_body_header'>Расписание группы: </p>
                        
                </div>
            </div>
        </div>
    )
}


const TimeTable = () => {
    const rigth = 0;
    async function gettable(){
        try {
            const response = await axios.get('http://mirea-club.site/api/timetable', {
                params: {
                    group: 'ИКБО-02-20'
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            });
            return response.data;
        } catch (error) {
            return 'error';
        }
    };
    async function settable(){
        const timetable = await gettable();
        rigth = 1;
        console.log(timetable);
    }
    if (rigth){
        return(
            <div className='timetable_main'>
                <ChoiceMenu/>
                <Tabletime/>
            </div>
          );
    }
    return(
        <div className='timetable_main'>
            <ChoiceMenu/>
            <div className='timetable_body'>
            <p>Расписание не получено, кликните на кнопку, чтобы получить расписание вашей группы на неделю</p>
            <button onClick={settable()}>Получить</button>
        </div>
        </div>
      );
  
}

export default TimeTable;