import axios from 'axios';

let initialState={
    getTable:0,
    parametr:[],
    table:[]
}


const ReduserTimeTable = (state = initialState, action) =>{
    if (action.type === 'UPDATE-PARAMETR'){
        console.log("update")
        return state;
    }
    else if (action.type === 'GET-TABLE'){
        const response =axios.get('http://mirea-club.site/api/timetable', {
                params: {
                    group: 'ИКБО-02-20'
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => {
                const tableData = response.data.data[0]; // Получаем данные таблицы из response.data.data.weeks
                console.log(tableData); // Выводим данные таблицы в консоль
                // Далее можно обновить состояние, сохранить данные в store или передать их в компоненты для отображения
            })
            .catch(error => {
            console.error(error);
            // Обработка ошибки при получении данных
            });
                
        state.getTable = 1;
        console.log(state.table)
        return state;
    }
    else{
        return state;
    }
}


export const getTimeTable = () =>({type: 'GET-TABLE'})
export const updateCreator = (text) =>({type: 'UPDATE-PARAMETR', newText: text})
export default ReduserTimeTable;