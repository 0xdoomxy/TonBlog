import React, { useEffect,useState } from "react";
import {useNavigate} from "react-router-dom";
import * as echarts from 'echarts';
import { useTonWallet,useTonConnectUI,toUserFriendlyAddress  } from '@tonconnect/ui-react';
import { Header, TonAvatar } from "../components";

const ArchivesPage =()=>{
    const navigate=useNavigate();
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
    //访问量表格数据
    const [accessnum,setAccessNum] = useState([{value:1048,date: "1123"}]);
    const accessTableId =  "accessTable";

    useEffect(()=>{
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
               <Header/>
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

export default ArchivesPage;