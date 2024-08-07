import { useState, useEffect, useRef } from 'react';
import { toast } from 'react-toastify';
import { useIsConnectionRestored, useTonConnectUI, useTonWallet } from '@tonconnect/ui-react';
import { UserClient,SetAuthorizetion, Authorization } from '../agent/agent';
const localStorageKey = 'blog-auth-token';
const payloadTTLMS = 1000 * 60 * 20;
 const   useBackendAuth = ()=> {
    const isConnectionRestored = useIsConnectionRestored();
    const wallet = useTonWallet();
    const [tonConnectUI] = useTonConnectUI();
    const interval = useRef();
    const [monitorLogin, setMonitorLogin] = useState(null);
    useEffect(() => {
        if (!isConnectionRestored || !SetAuthorizetion) {
            return;
        }

        clearInterval(interval.current);

        if (!wallet) {
            localStorage.removeItem(localStorageKey);
            SetAuthorizetion(null);

            const refreshPayload = async () => {
                tonConnectUI.setConnectRequestParameters({ state: 'loading' });

                const value = await  UserClient.Generate();
                if (!value) {
                    tonConnectUI.setConnectRequestParameters(null);
                } else {
        
                    tonConnectUI.setConnectRequestParameters({state: 'ready',value:{
                        tonProof: value.payload
                    }});
                }
            }

            refreshPayload();
            setInterval(refreshPayload, payloadTTLMS);
            return;
        }

        const token = localStorage.getItem(localStorageKey);
        if (token) {
            SetAuthorizetion(token);
            return;
        }

        if (wallet.connectItems?.tonProof && !('error' in wallet.connectItems.tonProof)) {
          UserClient.Proof({publickey:wallet.account.publicKey,address:wallet.account.address,network:wallet.account.chain,proof:wallet.connectItems.tonProof.proof}).then(result => {    
            try{
            if (result.status) {
                SetAuthorizetion(result.data);
                localStorage.setItem(localStorageKey, result.data);
                if(monitorLogin===undefined &&monitorLogin===null){
                    const inspectLogin = function(){
                        if(Authorization === undefined ||Authorization === null){
                            if(tonConnectUI.connected){
                                tonConnectUI.disconnect();
                            }
                        }
                    }
                    inspectLogin();
                  setMonitorLogin(setInterval(inspectLogin,payloadTTLMS));
                } 
            }else {
                    toast.error('登陆失败');
                    tonConnectUI.disconnect();
                }
            }catch(e){
                toast.error('登陆失败');
                if (tonConnectUI.connected){
                tonConnectUI.disconnect();
                }
            }
            })
        } else {
          toast.error('请尝试换一个钱包');
            tonConnectUI.disconnect();
        }

    }, [wallet, isConnectionRestored, SetAuthorizetion])
}


export { useBackendAuth };