export const getRap = (inputBLOB) => ({
    type: 'GET_RAP',
    payload: {
        inputBLOB: inputBLOB
    }
});

export const fetchRapEnded = (sucess,outputUUID = null) => ({
    type: 'FETCH_RAP_ENDED',
    payload: {
        sucess: sucess,
        outputUUID: outputUUID
    }
})

export const downloadOutput = (uuid) => ({
    type: 'DOWNLOAD_OUTPUT',
    payload: {
        uuid: uuid
    }
})