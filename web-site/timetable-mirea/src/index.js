import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.js';
import Header from './components/header/header.js'
import {BrowserRouter} from 'react-router-dom'
import store from './redux/store.js'


const root = ReactDOM.createRoot(document.getElementById('root'));
let renderAll = (state) =>{
  root.render(
    <React.StrictMode>
      <BrowserRouter>
      <Header />
      <App state={state} dispatch={store.dispatch.bind(store)}/>
      </BrowserRouter>
    </React.StrictMode>
  );
}

renderAll(store.getState());

store.subscribe(()=>{
  let state = store.getState();
  renderAll(state);
});
