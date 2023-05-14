import {combineReducers, legacy_createStore as createStore} from 'redux';
import ReduserTimeTable from './ReduserTable';


let reduser = combineReducers({
    parametrTable:ReduserTimeTable,
})

let store = createStore(reduser);

export default store;