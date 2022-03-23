import React from 'react';
import ReactDOM from 'react-dom';
import './index.scss';
import { App } from './App';
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';
import store from './store';
import { connect } from './services/websocket';

connect('ws://localhost:8080/ws', store.dispatch);

ReactDOM.render(
  <Provider store={store}>
    <React.StrictMode>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </React.StrictMode>
  </Provider>,
  document.getElementById('root'),
);
