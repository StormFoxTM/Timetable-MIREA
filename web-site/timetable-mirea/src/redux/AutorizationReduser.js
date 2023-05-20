import axios from 'axios'

let initialState = {
    login: '',
    password: ''
}
 
const UserReduser = (state = initialState, action) =>{
    if (action.type === 'AUTORIZATION-USER'){
        console.log(action.login_user);
        console.log(action.password_user);
        var autorization = axios.get('http://mirea-club.site/api/users', {
                params: {
                    login: action.login_user,
                    password: action.password_user
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(
                console.log(autorization) 
            );
        return state;
    } else if (action.type === 'REGISTRATION-USER'){
        console.log(action.login_user)
        console.log(action.password_user)
        var registration = axios.post('http://mirea-club.site/api/users', {
                params: {
                    login: action.login_user,
                    password:action.password_user
                },
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(
                console.log(registration) 
            );
        return state;
    } 
    else
        return state;
} 


export const registrationCreator = (login, password) =>({type: 'REGISTRATION-USER', login_user: login, password_user: password});
export const autorizationCreator = (login, password) =>({type: 'AUTORIZATION-USER', login_user: login, password_user: password});
export default UserReduser;