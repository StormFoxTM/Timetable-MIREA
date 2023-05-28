import axios from 'axios'

let initialState = {
    login: '',
    password: '',
    role: ''
}
 
const UserReduser = (state = initialState, action) =>{
    if (action.type === 'AUTORIZATION-USER'){
        // state.login = action.login_user
        axios.get('http://localhost:8080/api/users', {
            params: {
                username: action.login_user,
                password: action.password_user
            },
            headers: {
                'Content-Type': 'application/json',
            }
        })
        .then( response => {
            console.log(response.data) 
        })
        .catch(error => {
            console.error(error);
        });
        return state;
    } else if (action.type === 'REGISTRATION-USER'){
        axios.post('http://localhost:8080/api/users', {
            username: action.login_user,
            password: action.password_user
            }, {
            headers: {
                'Content-Type': 'application/json',
            }
        })
        .then(response => {
            console.log(response.data) 
        })
        .catch(error => {
            console.error(error);
        });
        return state;
    } 
    else
        return state;
} 


export const registrationCreator = (login, password) =>({type: 'REGISTRATION-USER', login_user: login, password_user: password});
export const autorizationCreator = (login, password) =>({type: 'AUTORIZATION-USER', login_user: login, password_user: password});
export default UserReduser;