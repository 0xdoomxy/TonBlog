import React, { useEffect,useState } from 'react';
import {TagClient} from "../agent/agent";
import Constants from "../util/constants";
import { useNavigate, useSearchParams  } from "react-router-dom";
import { Header, Spin } from '../components';
import { ToastContainer, toast } from 'react-toastify';
import { Pagination,Empty,Tag,BackTop } from 'antd';
const TagDetails = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const tag = searchParams.get("tag");
    const labelColorList = ["blue", "purple", "cyan", "green", "magenta", "pink", "red", "orange", "yellow", "volcano", "geekblue", "lime", "gold"];
    const [articlesByTag,setArticlesByTag] = useState(undefined);
    const [pageView,setPageView] = useState({
        total:0,
        current:1
    });
    const [isEmpty,setIsEmpty] = useState(false);
    //正在搜索
    const [isLoad,setIsLoading] = useState(true);
 //搜索文章
  function FindArticlesByTag(){
    TagClient.GetArticleByTag(tag,pageView.current,Constants.PageSize).then((data)=>{
        if(!data.status){
            let msg =data.message;
            if(msg === undefined || msg === null){
                msg = "系统出错啦";
            }
            toast.error(msg);
            return;
        }
        setArticlesByTag(data.data.articles.map((item)=>{
            if (item.tags !== "") {
                item.tags = item.tags.split(",");
            } else {
                item.tags = [];
            }
            item.create_time = new Date(item.create_time).toLocaleDateString("zh-CN", {timeZone: "Asia/Shanghai", year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'});
            return item;
        }));
        setPageView((origin)=>{return {...origin,total:data.data.total}});
        setIsLoading(false);
    })
}
useEffect(()=>{
    if(articlesByTag===undefined){
        setIsEmpty(true);
        return;
    }
    if(articlesByTag.length<=0){
        setIsEmpty(true);
    }else{
        setIsEmpty(false);
    }
},[articlesByTag])
    useEffect(()=>{
       //初始化要查找的热点文章
       FindArticlesByTag();
    },[pageView.current,pageView.total])
    //初始化函数
    useEffect(()=>{
        //初始化要查找的热点文章
        FindArticlesByTag();
    },[])
  return (
    <div className=" w-full h-full">
   <BackTop />
        <ToastContainer  />
    {/* header信息 */}
   <Header/>
    {/**热点文章内容主体 */}
    {isLoad?<div className='w-full h-full flex justify-center items-center'><Spin isSpin={isLoad} className=" w-20 h-20"/></div>:<div className='flex justify-center items-center'>
        <div className=' w-1/5 h-full'></div>
        {isEmpty?<div className=' h-screen flex justify-center items-center'><Empty className='pt-12' image={Empty.PRESENTED_IMAGE_SIMPLE} /></div>: <div className='w-3/5 h-full pt-12'>
        <div className=" w-full mt-8">
    {articlesByTag.map((item,index)=>(<div className={`px-2 hover:shadow-lg  transition duration-500 ease-in-out hover:-translate-y-1 hover:scale-105  my-3 min-h-32  border-2 w-full flex  justify-between rounded-md`} key={"newArticle"+index}>
        <div className="flex w-2/3 flex-col justify-center">
        <p className=" font-serif md:text-2xl py-1">{item.title}</p>
        <div className=" flex py-1">
        {item.tags!=null && item.tags instanceof Array &&item.tags.length>0&&item.tags.map((tag,index)=>(<Tag key={"tag"+index} color={labelColorList[index%labelColorList.length]}>{tag}</Tag>) )}
        </div>
        <div className=" font-normal text-md truncate">{item.creator}</div>
        <div  className=" font-normal text-sm">{item.create_time}</div>
        </div>
        <div className=" flex justify-center w-1/3 items-center flex-col">
            <button className=" w-20 h-12 border-2 rounded-xl hover:bg-blue-100" onClick={()=>navigate("/article/"+item.id)}>阅读</button>
            <div className=" font-serif text-ellipsis text-sm">浏览量:{item.access_num}</div>
        </div>
    </div>))}
</div>
<Pagination align="end" onChange={(page)=>{setPageView((origin)=>({total:origin.total,current:page}))}} current={pageView.current} pageSize={Constants.PageSize}  defaultCurrent={pageView.current} total={pageView.total} />
        </div>}
        <div className=' w-1/5 h-full'></div>
    </div>}
    </div>
  );
};


export default TagDetails;