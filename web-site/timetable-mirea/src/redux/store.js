import {combineReducers, legacy_createStore as createStore} from 'redux';
import ReduserTimeTable from './ReduserTable';
import UserReduser from './AutorizationReduser';


let reduser = combineReducers({
    parametrTable:ReduserTimeTable,
    userData:UserReduser,
})

let store = createStore(reduser);

export default store;