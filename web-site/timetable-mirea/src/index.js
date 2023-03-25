import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.js';
import Header from './components/header/header.js'
import Footer from './components/footer/footer.js'
import {BrowserRouter} from 'react-router-dom'


const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <BrowserRouter>
    <Header />
    <App />
    <Footer info='&copy; Copyright 2023, TimeTable MIREA'/>
    </BrowserRouter>
  </React.StrictMode>
);
