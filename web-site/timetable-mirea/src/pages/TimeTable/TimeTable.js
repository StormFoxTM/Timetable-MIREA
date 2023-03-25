import './TimeTable.css';

function TimeTable() {
  return (
    <body>
      <div className='timetable_main'>
        <div className='timetable_panel'>
            <div className='timetable_panel_wrapper'>
                <form>
                    <p className='timetable_type_info'><b>Вид расписания:</b></p>
                    <p>
                        <button type="submit" className='change_type_info' name="type_of_info" value="day">День</button>
                        <button type="submit" className='change_type_info' name="type_of_info" value="week">Неделя</button>
                    </p>
                </form> 
            </div>
        </div>
        <div className='timetable_body'>
            <div className='timetable_body_wrapper'>
                <div>
                    <p className='timetable_body_header'>Расписание группы: </p>
                        <table className='timetable_table'>
                                <tr className='timetable_column'>
                                   
                                </tr>
                            
                        </table>
                </div>
            </div>
        </div>
    </div>
    </body>
  );
}

export default TimeTable;
