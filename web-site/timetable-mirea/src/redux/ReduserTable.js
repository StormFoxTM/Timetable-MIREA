import axios from 'axios';
import {getDay, getISOWeek} from 'date-fns';

let initialState={
    getTable:0,
    parametr:'',
    data: '',
    data_type: '',
    table:[]
}

function determineType(value) {
    if (/^[А-Я]{1,3}-\d{1,3}(-\d{1,2})?$/i.test(value)) {
    // if (/\w\w\w\w-\d\d-\d\d/i.test(value)) {
      return "auditorium";
    } else if (/^[А-Я]{4}-\d{2}-\d{2}$/i.test(value)) {
    // } else if (/\w\w\w\w-\d\d-\d\d/i.test(value)) {
      return "group";
    } else{
      return "lecturer";
    } 
  }

function determineTypeDate(value){
    if (value === 'today' || value === 'tomorrow') {
        return 'day';
    }
    else if (value === 'week' || value === 'next_week'){
        return 'week';
    }
}

function determineDate(value) {
    if (value === 'today' ) {
        const currentDate = new Date();
        const dayOfWeek = getDay(currentDate);
        return dayOfWeek;
    } else if (value === 'tomorrow') {
        const currentDate = new Date();
        const dayOfWeek = getDay(currentDate) + 1;
        return dayOfWeek;
    } else if (value === 'week' || value === 'next_week') {
        const date = new Date(); // текущая дата
        const year = date.getFullYear();
        const month = date.getMonth();

        // определяем первый день в феврале текущего года
        const firstDayOfFebruary = new Date(year, 1, 1);
        const firstWeekOfYear = getISOWeek(firstDayOfFebruary);

        let weekNumber;
        if (month >= 1 && month <= 6) {
        // месяц между февралем и июлем
        const firstDayOfCurrentYear = new Date(year, 0, 1);
        const daysSinceFirstWeek = (date - firstDayOfCurrentYear) / (1000 * 60 * 60 * 24);
        weekNumber = Math.ceil((daysSinceFirstWeek + firstDayOfCurrentYear.getDay() + 1) / 7) - 5;
        } else {
        // месяц между сентябрем и февралем следующего года
        const firstDayOfWorkingYear = new Date(year, 8, 1);
        const daysSinceFirstWeek = (date - firstDayOfWorkingYear) / (1000 * 60 * 60 * 24);
        weekNumber = Math.ceil((daysSinceFirstWeek + firstWeekOfYear) / 7);
        }
        if (value === 'week')
            return weekNumber % 2 + 1;
        else
            return (weekNumber + 1) % 2 + 1;
    }
}

function ReduserTimeTable (state = initialState, action){
    if (action.type === 'GET-TABLE'){
        const paramName = determineType(action.param);
        const paramValue = action.param; 
        state.parametr = paramValue;
        const response = axios.get('http://mirea-club.site/api/timetable', {
                params: {
                    [paramName]: paramValue,
                    'week':'1'
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
    } else if (action.type === 'GET-DATE'){
        const paramDate = determineDate(action.param_ch);
        const paramTypeDate = determineTypeDate(action.param_ch);
        state.data = paramDate;
        state.data_type = paramTypeDate;
        const paramName = determineType(state.parametr);
        const paramValue = state.parametr; 
        if (paramTypeDate === 'day'){
            const response = axios.get('http://mirea-club.site/api/timetable', {
                params: {
                    [paramName]: paramValue,
                    'week': '1',
                    [paramTypeDate]:paramDate
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
        }
        else{
            const response = axios.get('http://mirea-club.site/api/timetable', {
                params: {
                    [paramName]: paramValue,
                    [paramTypeDate]:paramDate
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
        }
        
                
        state.getTable = 0;
        return state;
    }
    else if (action.type === 'GET-TIMETABLE'){
        const paramDate = state.data;
        const paramTypeDate = state.data_type;
        const paramName = determineType(state.parametr);
        const paramValue = state.parametr; 
        if (paramTypeDate === 'day'){
            const response = axios.get('http://mirea-club.site/api/timetable', {
                params: {
                    [paramName]: paramValue,
                    'week': '1',
                    [paramTypeDate]:paramDate
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => {
                const tableData = response.data.weeks;
                state.table = tableData
                console.log(tableData)
            })
            .catch(error => {
            console.error(error);
            });
        }
        else{
            const response = axios.get('http://mirea-club.site/api/timetable', {
                params: {
                    [paramName]: paramValue,
                    [paramTypeDate]:paramDate
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => {
                const tableData = response.data.weeks;
                state.table = tableData
                console.log(tableData)
            })
            .catch(error => {
            console.error(error);
            });
        }
        state.getTable = 1;
        return state;
    } else{
        return state;
    }
}


export const getTimeTable = (text) =>({type: 'GET-TABLE', param: text})
export const getDateTable = (chose) =>({type: 'GET-DATE', param_ch: chose})
export const getResponse = () =>({type: 'GET-TIMETABLE'})
export default ReduserTimeTable;