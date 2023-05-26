import axios from 'axios'

let initialState = {
    login: '',
    password: '',
    role: ''
}
 
const UserReduser = (state = initialState, action) =>{
    if (action.type === 'AUTORIZATION-USER'){
        var autorization = axios.get('http://localhost:9888/api/users', {
                params: {
                    username: action.login_user,
                    password: action.password_user
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(
                console.log(autorization) 
            )
            .catch(error => {
                console.error(error);
                });
        return state;
    } else if (action.type === 'REGISTRATION-USER'){
        var registration = axios.post('http://localhost:9888/api/users', {
                params: {
                    username: action.login_user,
                    password:action.password_user
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(
                console.log(registration) 
            )
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