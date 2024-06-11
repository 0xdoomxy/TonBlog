import React, { useEffect,useState,useRef } from "react";
import {useNavigate} from "react-router-dom";
import * as echarts from 'echarts';


const Archives =()=>{
    const navigate=useNavigate();
    const navItems=[{
        Name:"Home",
        Target:"/"
    },{
        Name:"About",
        Target:"/about"
    },{Name:"Archieve",Target:"/archieve"}]
    //访问量表格数据
    const [accessnum,setAccessNum] = useState([{value:1048,date: "1123"}]);
    const accessTableId =  "accessTable";

    useEffect(()=>{
        console.log(document.getElementById(accessTableId))
        echarts.init(document.getElementById(accessTableId)).setOption({
            xAxis:{
                type:'category',
            data:accessnum.map((item)=>{return item.date}),
            },
            yAxis:{
                type:"value"
            },
            series:[{
                data:accessnum.map((item)=>{return item.value}),
                type:'line'
            
            }]
        })
    },[accessnum])
     //小屏幕点击事件，用来显示菜单栏
     const [showSmallNav,setShowSmallNav]=useState(false);
    return (
        <div className="w-full h-full">
               {/* header信息 */}
               <div className="w-full fixed z-10">
               <div className="  bg-slate-50 w-full border-b-2 h-12 flex justify-evenly md:justify-center items-center ">
                <div  className=" w-1/4 flex justify-center   items-center py-2">
                <h1 className=" flex align-middle font-serif text-wrap h-full text-xl md:text-3xl cursor-pointer"  onClick={()=>{window.location.href="https://github.com/0xdoomxy"}}>0xdoomxy</h1>
                </div>
                <div className="w-1/2   hidden md:flex justify-start items-center">
                        {navItems.map((item,index)=>(
                            <div onClick={()=>{navigate(item.Target)}} className=" hover:-translate-y-1 duration-500  text-center text-lg px-8 cursor-pointer " key={"nav"+index}>{item.Name}</div>
                        ))}
                        <div className=" pl-24 ">
                            <div className=" cursor-pointer  ">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
<path strokeLinecap="round" strokeLinejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
</svg>
</div>
                        </div>
                </div>
                {/* 小屏幕显示 */}
                <div className=" flex  pl-12 justify-center items-center  w-1/3 md:hidden ">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6 cursor-pointer" onClick={()=>{setShowSmallNav(!showSmallNav)}}>
<path strokeLinecap="round" strokeLinejoin="round" d="M3.75 5.25h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5" />
</svg>
                </div>
                {/* TODO 点击事件 */}
            </div>
            {showSmallNav&&<div className="  bg-gray-50 border-b  z-10 relative flex w-full md:hidden flex-col justify-center items-center">
                {navItems.map((item,index)=>(
                            <div onClick={()=>{navigate(item.Target)}} className="w-full border-y hover:decoration-sky-700 hover:underline  text-center text-lg px-8 cursor-pointer " key={"smallnav"+index}>{item.Name}</div>
                        ))}
                </div>}
                </div>
            <div className=" flex justify-center w-full h-full pt-12">
                    <div className=" w-1/6"></div>
                    <div className="w-2/3 pt-8">
                        {/* 文本框 */}
                        <div className=" w-full h-32">
                            <p className=" text-ellipsis  indent-4   md:text-xl w-full">博客建立于2024年6月2日,旨在于自我学习过程中的经验分享。</p>   
                        </div>
                        {/* 浏览量访问展示图 */}
                        <div className=" flex justify-between">
                            <div id={accessTableId}></div>
                            <div>there</div>
                        </div>
                        <div></div>
                    </div>
                    <div className="w-1/6"></div>
            </div>
            {/* TODO */}
        </div>
    )

}

export default Archives;