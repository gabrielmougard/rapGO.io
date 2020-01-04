import { GENERATOR_SERVER } from '../CONSTANTS.js';
import { call, put, takeLatest, all } from 'redux-saga/effects';
import axios from 'axios';

function* fetchRap(action) {
    //API call here...
    const { inputBLOB } = action.payload;
    try {
        var response = yield call([axios, axios.post], 'http://'+ GENERATOR_SERVER + '/generate');
    } catch (e) {
        console.log(e)
        yield put(fetchRapEnded(false));
    }

}

function* actionWatcher() {
    yield takeLatest('GET_RAP', fetchRap)
}

export default function* rootSaga() {
    yield all([
        actionWatcher(),
    ]);
}