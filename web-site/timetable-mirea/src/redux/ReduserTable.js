import axios from 'axios';

let initialState={
    getTable:0,
    parametr:[],
    table:[]
}


const ReduserTimeTable = (state = initialState, action) =>{
    if (action.type === 'GET-TABLE'){
        const response = axios.get('http://mirea-club.site/api/timetable', {
                params: {
                    group: action.param
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => {
                const tableData = response.data.weeks;
                state.table = tableData
            })
            .catch(error => {
            console.error(error);
            });
                
        state.getTable = 1;
        return state;
    }
    else{
        return state;
    }
}


export const getTimeTable = (text) =>({type: 'GET-TABLE', param: text})
export default ReduserTimeTable;