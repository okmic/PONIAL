import { createSlice } from '@reduxjs/toolkit'

interface TaxonomyState {

}

const initialState: TaxonomyState = {

}

const staticDataSlice = createSlice({
    name: 'staticData',
    initialState,
    reducers: {

    }
})

export const {  } = staticDataSlice.actions
export default staticDataSlice.reducer
