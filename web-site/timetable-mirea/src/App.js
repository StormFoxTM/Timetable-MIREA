import TimeTable from './pages/TimeTable/TimeTable';
import Authorization from './pages/Register/Authorization';
import Register from './pages/Register/Register';
import Setting from './pages/Setting/Setting';
import {Routes, Route} from 'react-router-dom'
import './App.css';
import My_page from './pages/Register/My_page';

function App(props) {
  return (
    <div>
      <Routes>
      <Route path='/' element={<TimeTable state={props.state.parametrTable} dispatch={props.dispatch}/>}/>
      <Route path='/authorization' element={<Authorization state={props.state.userData} dispatch={props.dispatch}/>}/>
      <Route path='/register' element={<Register state={props.state.userData} dispatch={props.dispatch}/>}/>
      <Route path='/setting' element={<Setting />}/>
      <Route path='/mypage' element={<My_page dispatch={props.dispatch}/>}/>
      </Routes>
    </div>
  );
}

export default App;
