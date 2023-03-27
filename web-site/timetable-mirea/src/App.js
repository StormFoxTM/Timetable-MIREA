import TimeTable from './pages/TimeTable/TimeTable';
import Authorization from './pages/Register/Authorization';
import Register from './pages/Register/Register';
import Setting from './pages/Setting/Setting';
import {Routes, Route} from 'react-router-dom'
import './App.css';

function App() {
  return (
    <body>
      <Routes>
      <Route path='/' element={<TimeTable/>}/>
      <Route path='/authorization' element={<Authorization/>}/>
      <Route path='/register' element={<Register/>}/>
      <Route path='/setting' element={<Setting />}/>
      </Routes>
    </body>
  );
}

export default App;
