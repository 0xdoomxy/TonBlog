




let providers = []




export const store = {
    value:()=>providers,
    subscribe:(callback)=>{
        function onAnnouncement(event){
            if(providers.map(p=>p.info.uuid).includes(event.detail.info.uuid)) return
            providers = [...providers,event.detail]
            callback()
        }
        window.addEventListener("eip6963:announceProvider",onAnnouncement)
        window.dispatchEvent(new Event("eip6963:requestProvider"));
        return ()=>window.removeEventListener("eip6963:announceProvider", onAnnouncement)
    }
}