
// const Pool = require("pg");


let initialState={
    getTable:0,
    parametr:[]
}


const Postgresql = (state = initialState, action) =>{
    if (action.type === 'Registration'){
        console.log("Registration")
        // connectDb(state)
        return state;
    }
    else if (action.type === 'Authorization'){
        console.log("Authorization")
        // connectDb(state);
        return state;
    }
    else {
        // connectDb(state);
        return state;
    }
}
 
// const connectDb = async (state) => {
//     try {
//         const pool = new Pool({
//             user: "admin",
//             host: "postgres",
//             database: "postgres",
//             password: "admin",
//             port: "5432",
//         });
 
//         await pool.connect()
//         const res = await pool.query('SELECT * FROM clients')
//         console.log(res)
//         await pool.end()
//     } catch (error) {
//         console.log(error)
//     }
// }

export default Postgresql;