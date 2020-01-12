import { CONVERTER_SERVER } from '../CONSTANTS.js';
import { fetchRapEnded } from '../actions/index';
import { call, put, takeLatest, all } from 'redux-saga/effects';
import axios from 'axios';

function* fetchRap(action) {
    const { inputBLOB } = action.payload;
    console.log("[SAGA] blob")
    console.log(inputBLOB)
    const data = new FormData();
    data.append('file', inputBLOB.blob, "recording.mp3");

    try {
        var response = yield call([axios, axios.post], 'http://converterserver:3001/test', data, { // http://converterserver:3001/upload
            headers: {
                'Content-Type': `multipart/form-data; boundary=${data._boundary}`,
            }
        });
        const { status, outputBLOB } = response.data
        if (status == 200 && outputBLOB) {
            console.log("[SAGA] : outputBLOB detected.");
            yield put(fetchRapEnded(true, outputBLOB));
        } else {
            console.log("[SAGA] : error, status code is "+status)
            yield put(fetchRapEnded(false));
        }
    } catch (e) {
        console.log(e)
        yield put(fetchRapEnded(false));
    }

}

function* actionWatcher() {
    yield takeLatest('GET_RAP', fetchRap);
}

export default function* rootSaga() {
    yield all([
        actionWatcher(),
    ]);
}