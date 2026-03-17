import { configureStore } from '@reduxjs/toolkit'
import authReducer from './slices/auth.slice'
import taxonomySlice from './slices/taxonomy.slice'

const store = configureStore({
  reducer: {
    auth: authReducer,
    taxonomy: taxonomySlice,
  },
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch

export default store
