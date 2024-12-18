


import {useSyncExternalStore} from "react";
import {store} from './store.js'
export const useSyncProviders = ()=> useSyncExternalStore(store.subscribe, store.value, store.value)