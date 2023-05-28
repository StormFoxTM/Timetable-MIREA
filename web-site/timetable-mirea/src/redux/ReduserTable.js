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
        const response = axios.get('http://mirea-club.site/api/timetable', {
                params: {
                    group: action.param
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => {
                const tableData = response.data; // Получаем данные таблицы из response.data.data.weeks
                console.log(tableData); // Выводим данные таблицы в консоль
            })
            .catch(error => {
            console.error(error);
            });
                
        state.getTable = 1;
        console.log(state.table)
        return state;
    }
    else{
        return state;
    }
}


export const getTimeTable = (text) =>({type: 'GET-TABLE', param: text})
export const updateCreator = (text) =>({type: 'UPDATE-PARAMETR', newText: text})
export default ReduserTimeTable;