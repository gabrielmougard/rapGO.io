const reducer = (state = {}, action) => {
    switch (action.type) {
        case 'GET_RAP':
            return { ...state, loading: true};
        case 'RAP_RECEIVED':
            return { ...state, rap: action.rapBLOB, loading: false }
        default:
            return state;
    }
};

export default reducer;