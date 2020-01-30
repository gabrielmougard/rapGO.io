import { CONVERTER_SERVER } from '../CONSTANTS.js';
import { fetchRapEnded, downloadEnded } from '../actions/index';
import { call, put, takeLatest, all } from 'redux-saga/effects';
import axios from 'axios';

function* fetchRap(action) {
    const { inputBLOB } = action.payload;
    console.log("[SAGA] blob")
    console.log(inputBLOB)
    const data = new FormData();
    data.append('file', inputBLOB.blob, "recording.mp3");

    try {

        var response = yield call([axios, axios.post], 'http://localhost:3001/upload', data, { // http://converterserver:3001/upload
            headers: {
                'Content-Type': `multipart/form-data; boundary=${data._boundary}`,
            }
        });

        console.log(response)
        const { Status, OutputUUID } = response.data
        if (Status == 200 && OutputUUID) {
            console.log("[SAGA] : outputUUID detected.");
            yield put(fetchRapEnded(true, OutputUUID));
        } else {
            console.log("[SAGA] : error, status code is "+Status)
            yield put(fetchRapEnded(false));
        }
    } catch (e) {
        console.log(e)
        yield put(fetchRapEnded(false));
    }

}

function* fetchOutput(action) {
    const { uuid } = action.payload;
    
    const data = {
        outputUUID: uuid
    }
    //TODO : call axios and download the mp3 file
    try {
        var response = yield call([axios, axios.post], 'http://localhost:3001/download', data, {
            responseType: 'arraybuffer',
            headers: {
                'Content-Type': `json/application`,
            }
        })
        yield put(downloadEnded(true, response)) //response.data is the raw data of the output_<uuid>.mp3
    } catch (e) {
        console.log(e)
        yield put(downloadEnded(false));
    }
}

function* actionWatcher() {
    yield takeLatest('GET_RAP', fetchRap);
    yield takeLatest('DOWNLOAD_OUTPUT', fetchOutput);
}

export default function* rootSaga() {
    yield all([
        actionWatcher(),
    ]);
}