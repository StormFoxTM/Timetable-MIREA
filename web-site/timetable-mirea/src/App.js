import Main from './pages/Main/Main';
import TimeTable from './pages/TimeTable/TimeTable';
import Authorization from './pages/Authorization/Authorization';
import Register from './pages/Register/Register';
import {Routes, Route} from 'react-router-dom'


function App() {
  return (
    <body>
      <Routes>
      <Route path='/' element={<Main/>}/>
      <Route path='/authorization' element={<Authorization/>}/>
      <Route path='/register' element={<Register/>}/>
      <Route path='/timeTable' element={<TimeTable/>}/>
      </Routes>
    </body>
  );
}

export default App;
