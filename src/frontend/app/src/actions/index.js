export const getRap = (inputBLOB) => ({
    type: 'GET_RAP',
    payload: {
        inputBLOB: inputBLOB
    }
});

export const fetchRapEnded = (sucess,outputBLOB = null) => ({
    type: 'FETCH_RAP_ENDED',
    payload: {
        sucess: sucess,
        outputBLOB: outputBLOB
    }
})