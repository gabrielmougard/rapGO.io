import React from 'react';
import ReactDOM from 'react-dom';
import createSagaMiddleware from 'redux-saga';
import { createStore, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';
import { logger } from 'redux-logger';

import reducer from './reducers';
import rootSaga from './sagas';
import App from './App';
import * as serviceWorker from './serviceWorker';

import {BaseProvider, LightTheme} from 'baseui';
import { Provider as StyletronProvider } from "styletron-react";
import { Client as Styletron } from "styletron-engine-atomic";

const engine = new Styletron();

const sagaMiddleware = createSagaMiddleware();

const store = createStore(
    reducer,
    applyMiddleware(sagaMiddleware, logger),
);

sagaMiddleware.run(rootSaga);

// for redux saga, see :
// https://medium.com/@lavitr01051977/make-your-first-call-to-api-using-redux-saga-15aa995df5b6
ReactDOM.render(
    <StyletronProvider value={engine}>
        <BaseProvider theme={LightTheme}>
            <Provider store={store}>
                <App />
            </Provider>
        </BaseProvider>
    </StyletronProvider>
    , document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
