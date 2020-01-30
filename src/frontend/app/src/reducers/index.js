const reducer = (state = {}, action) => {
    switch (action.type) {
        case 'GET_RAP':
            return { ...state, loading: true};
        case 'RAP_RECEIVED':
            return { ...state, rap: action.rapBLOB, loading: false }
        case 'FETCH_RAP_ENDED':
            return { ...state, heartbeatUUID: action.payload.outputUUID}
        case 'DOWNLOAD_ENDED':
            if (action.payload.success) {
                return { ...state, outputResponse: action.payload.response}
            } else {
                return state;
            }
        default:
            return state;
    }
};

export default reducer;