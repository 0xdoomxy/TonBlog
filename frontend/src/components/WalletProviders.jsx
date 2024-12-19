import {useSyncProviders} from "../hooks/useSyncProviders";
import {toast} from "react-toastify";
import {motion} from "framer-motion";
import {Modal} from "antd";

import {SetAuthorizetion, UserClient} from "../agent/agent";
import React, {useEffect} from "react";
import {MetaMaskSDK} from "@metamask/sdk"

export const DiscoverWalletProviders = (props) => {
    const MMSDK = new MetaMaskSDK({
        dappMetadata: {
            name: "0xdoomxy blog",
        },
        infuraAPIKey: "656103da5ad94eea9d35c65f78079af8",
    })
    const {searchWalletModal, setSearchWalletModal, selectedWallet, setSelectedWallet, setUserAccount} = props;
    const providers = useSyncProviders()
    const handleConnectAndSign = async (providerWithInfo, account) => {
        if (providerWithInfo) {
            const message = JSON.stringify({
                types: {
                    EIP712Domain: [
                        {
                            name: "name",
                            type: "string"
                        },
                        {
                            name: "version",
                            type: "string"
                        },
                        {
                            name: "chainId",
                            type: "uint256"
                        }
                    ],
                    Verify: [
                        {name: "content", type: "string"},
                        {name: "date", type: "uint256"}
                    ]
                },
                domain: {
                    chainId: providerWithInfo.provider.chainId,
                    name: "0xdoomxy blog",
                    version: providerWithInfo.provider.networkVersion,
                },
                primaryType: "Verify",
                message: {
                    content: "Welcome to 0xdoomxy blog",
                    date: Date.now(),
                },
            })

            var params = [account, message]
            var method = "eth_signTypedData_v4"

            var sign = await providerWithInfo.provider.request(
                {
                    method,
                    params,
                    from: account,
                });
            var loginResp = await UserClient.Login({
                "message": message,
                "sign": sign,
            })
            if (!loginResp || !loginResp.status) {
                return false;
            }
            SetAuthorizetion(loginResp.data);
            return true;
        }
        return false;
    }
    const handleConnect = async (providerWithInfo) => {
        const accounts = await (
            providerWithInfo.provider.request({method: 'eth_requestAccounts'})
                .catch((err) => toast.error("获取账号失败", err))
        )
        if (accounts?.[0]) {
            let res = await handleConnectAndSign(providerWithInfo, accounts?.[0]);
            if (!res) {
                toast.error("登录出错");
                return;
            }
            setSelectedWallet(providerWithInfo);
            setUserAccount(accounts?.[0])
            setSearchWalletModal(!searchWalletModal);
            toast.success("登录成功");
        }
    }
    useEffect(() => {
        if (searchWalletModal && providers.length <= 0) {
            toast.error("未找到符合EIP6963的钱包")
            setSearchWalletModal(!searchWalletModal);
        }
    }, [searchWalletModal]);
    return (
        providers.length > 0 &&
        <Modal
            className={"rounded-xl"}
            width={"40%"} closable={false} keyboard footer={null} open={searchWalletModal}
            onCancel={() => setSearchWalletModal(!searchWalletModal)}>
            <div className={"w-full h-24 flex justify-center items-center flex-col"}>
                <motion.div animate={{x: 80, transition: {duration: 1}}}
                            className={" w-full flex justify-start font-serif   items-start flex-col"}>
                    <div className={"w-full lg:text-3xl align-middle font-serif text-wrap text-2xl"}>
                        欢迎来到
                    </div>
                    <div
                        className={"w-full align-middle font-serif pl-20 lg:text-3xl text-2xl"}>0xdoomxy的小世界
                    </div>
                </motion.div>
            </div>
            <div className={"flex justify-center items-center flex-col"}>{providers?.map((provider) => (
                <motion.button style={{width:"100%",height:"100%"}} className={"w-full h-full  my-1 motion-button"} onClick={() => handleConnect(provider)}
                               whileHover={{scale: 1.1}}
                               whileTap={{scale: 0.95}}>
                    <div className={" flex pl-2 justify-start items-center"}>
                        <img src={provider.info.icon} style={{width: "40px", height: "40px"}}
                             alt={provider.info.name}/>
                        <div className={"pl-4"} style={{
                            fontSize: "16px",
                            color: "#222222",
                            fontFamily: "Basel,sans-serif"
                        }}>{provider.info.name}</div>
                    </div>
                </motion.button>
            ))}</div>
        </Modal>

    )
}