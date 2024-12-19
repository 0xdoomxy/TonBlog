import {Header, RunningAirport,FinishAirport} from "../components";
import {AnimatePresence, motion} from "framer-motion";
import React, {useEffect, useState} from 'react';
import "../css/airport.css";


const AirPort = () => {
    const [selectedTab, setSelectedTab] = useState(null);
    const [isAdmin, setIsAdmin] = useState(false);
    const tabs = [
        {icon: "🥬", label: "正在进行的空投",content:<RunningAirport isAdmin={isAdmin}/>},
        {icon: "🍅", label: "已经结束的空投",content:<FinishAirport isAdmin={isAdmin}/>},
    ]    
    useEffect(()=>{
        setSelectedTab(tabs[0]);
    },[])
    return (
        <div className={"w-full h-full flex justify-center items-start"}>
            <Header/>
            <div className={"w-full h-full flex justify-center pt-32 items-center flex-row "}>
                <div className="airpointwindow flex justify-start items-center">
                    <nav className={" justify-center w-full flex items-center"}>
                        <ul className={"w-full flex justify-center items-center"}>
                            {tabs.map((item) => (
                                <li
                                    key={item.label}
                                    className={item === selectedTab ? "selected flex justify-center items-center lg:text-3xl text-xl" : "items-center justify-center lg:text-3xl text-xl "}
                                    onClick={() => setSelectedTab(item)}
                                >
                                    {`${item.icon} ${item.label}`}
                                    {item === selectedTab ? (
                                        <motion.div className="underline" layoutId="underline"/>
                                    ) : null}
                                </li>
                            ))}
                        </ul>
                    </nav>
                    <main className={"w-full pt-20 "}>
                        <AnimatePresence mode="wait">
                            <motion.div
                                key={selectedTab ? selectedTab.label : "empty"}
                                initial={{y: 10, opacity: 0}}
                                animate={{y: 0, opacity: 1}}
                                exit={{y: -10, opacity: 0}}
                                transition={{duration: 0.2}}
                            >
                                {selectedTab ?selectedTab.content: "😋"}
                            </motion.div>
                        </AnimatePresence>
                    </main>
                </div>
            </div>
        </div>
    )
};


export default AirPort;