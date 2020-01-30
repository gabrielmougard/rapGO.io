const reducer = (state = {}, action) => {
    switch (action.type) {
        case 'GET_RAP':
            return { ...state, loading: true};
        case 'RAP_RECEIVED':
            return { ...state, rap: action.rapBLOB, loading: false }
        case 'FETCH_RAP_ENDED':
            return { ...state, heartbeatUUID: action.payload.outputUUID}
        default:
            return state;
    }
};

export default reducer;