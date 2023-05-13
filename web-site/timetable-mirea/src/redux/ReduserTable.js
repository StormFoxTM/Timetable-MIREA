import axios from 'axios';

let initialState={
    getTable:0,
    parametr:[]
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
            });
        let timeTable = response
        state.rigth = 1;
        console.log(timeTable);
        return state;
    }
    else{
        return state;
    }
}


export const getTimeTable = () =>({type: 'GET-TABLE'})
export const updateCreator = (text) =>({type: 'UPDATE-PARAMETR', newText: text})
export default ReduserTimeTable;