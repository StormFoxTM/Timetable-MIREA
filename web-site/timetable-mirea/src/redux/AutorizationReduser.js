import axios from 'axios'

let initialState = {
    login: '',
    password: '',
    role: '',
    group: '',
    email: ''
}

const UserReduser = (state = initialState, action) =>{
    if (action.type === 'AUTORIZATION-USER'){
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
            if (response.data === 'user' || response.data === 'admin'){
                state.role = response.data;
                state.login = action.login_user
            }
            console.log(state.role, response.data)
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
            state.group = action.group_user
        })
        .catch(error => {
            console.error(error);
        });
        return state;
    } else if (action.type === 'LOGOUT-USER'){
        state.login = '';
        state.role = '';
        state.password = '';
        state.email = '';
        return state;
    } else
        return state;
} 


export const registrationCreator = (login, password, group) =>({type: 'REGISTRATION-USER', login_user: login, password_user: password, group_user:group});
export const autorizationCreator = (login, password) =>({type: 'AUTORIZATION-USER', login_user: login, password_user: password});
export const logoutCreator = () =>({type: 'LOGOUT-USER'});
export default UserReduser;