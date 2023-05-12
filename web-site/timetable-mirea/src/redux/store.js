import {combineReducers, legacy_createStore as createStore} from 'redux';
import ReduserTimeTable from './ReduserTable';
import Postgresql from './Posgresql'

let reduser = combineReducers({
    parametrTable:ReduserTimeTable,
    userData:Postgresql,
})

let store = createStore(reduser);

export default store;