
import React, { useEffect, useState } from 'react';

import TonAvatar from './avatar';
import { useTonWallet,useTonConnectUI,toUserFriendlyAddress  } from '@tonconnect/ui-react';
import { useNavigate } from "react-router-dom";
import { Search } from './search';
import {Modal,AutoComplete} from 'antd';
const Header = () => {
        //是否需要更换header显示
        const [changeHeader,setChangeHeader]=useState(false);
        const navigate = useNavigate();
        //搜索框显示
        const [openSearch,setOpenSearch]=useState(false);
            //小屏幕点击事件，用来显示菜单栏
    const [showSmallNav,setShowSmallNav]=useState(false);
        const navItems=[{
            Name:"Home",
            Target:"/"
        },{
            Name:"About",
            Target:"/about"
        },{Name:"Archieve",Target:"/archieve"}]
        //tron 钱包
        const wallet = useTonWallet();
        //tron 连接
        const [tonConnectUI] = useTonConnectUI();
    useEffect(()=>{
          //监听鼠标滚动事件来改变header
          const checkScroll =()=>{
            if(window.scrollY >200){
                setChangeHeader(true);
            }else{  
                setChangeHeader(false);
            }
        };
        window.addEventListener("scroll",checkScroll);
        return ()=>window.removeEventListener("scroll",checkScroll);
    },[])
    return (
        <div className="w-full fixed z-10 ">
            <Modal  transitionName="move-up">
            <AutoComplete/>
            </Modal>
        <div className="   bg-slate-50 w-full border-b-2 h-12 flex justify-evenly md:justify-center items-center ">
            {!changeHeader&&(<div className="w-full h-full flex items-center justify-center"><div  className=" w-1/4 flex justify-center   items-center py-2">
            <h1 className=" flex align-middle font-serif text-wrap h-full text-xl md:text-3xl cursor-pointer pl-2 "  onClick={()=>{window.location.href="https://github.com/0xdoomxy"}}>0xdoomxy</h1>
            </div>
            <div className="w-1/2   hidden  md:flex justify-start items-center">
                    {navItems.map((item,index)=>(
                        <div onClick={()=>{navigate(item.Target)}} className=" hover:-translate-y-1 duration-500  text-center text-lg px-4 lg:px-8 cursor-pointer " key={"nav"+index}>{item.Name}</div>
                    ))}
                    <div className=" lg:pl-24 pl-6   ">
                        <div className=" cursor-pointer  ">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
<path strokeLinecap="round" strokeLinejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
</svg>
</div>
                    </div>
            </div>
            <div className=" hidden md:flex w-1/8 justify-evenly "><TonAvatar wallet={wallet} disconnect={()=>{tonConnectUI.disconnect()}}/></div>
            </div>)}
            {changeHeader&&
            <div className="w-full h-full">
                <Search  onKeyDown={(event)=>{if(event.keyCode!==13){return;}if(event.target.value == undefined || event.target.value == null ){return }navigate(`/search?keyword=${event.target.value}`,)}}/>
               </div> 
                }
            {/* 小屏幕显示 */}
            <div className=" flex  pl-12 justify-center items-center  w-1/3 md:hidden ">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6 cursor-pointer" onClick={()=>{setShowSmallNav(!showSmallNav)}}>
<path strokeLinecap="round" strokeLinejoin="round" d="M3.75 5.25h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5" />
</svg>
            </div>
        </div>
        {showSmallNav&&<div className="w-full   h-full backdrop-blur absolute top-12    flex   flex-col justify-start items-center">
                {navItems.map((item,index)=>(
                            <div onClick={()=>{navigate(item.Target)}} className="w-full border-y hover:decoration-sky-700 hover:underline  text-center text-lg px-8 cursor-pointer " key={"smallnav"+index}>{item.Name}</div>
                        ))}
                </div>}
        </div>
    )
}

export default Header;