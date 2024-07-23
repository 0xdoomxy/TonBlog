import React, { useEffect,useState } from 'react';
import {ArticleClient} from "../agent/agent";
import Constants from "../util/constants";
import { useNavigate,useSearchParams  } from "react-router-dom";
import { Header, Spin } from '../components';
import { ToastContainer, toast } from 'react-toastify';
import {Tag,Empty } from 'antd';
const SearchPage = () => {
    const [params] = useSearchParams()
    const navigate = useNavigate();
    const labelColorList = ["blue", "purple", "cyan", "green", "magenta", "pink", "red", "orange", "yellow", "volcano", "geekblue", "lime", "gold"];
    const [searchArticles,setSearchArticles] = useState(undefined);
    const searchKeyword=params.get("keyword");
    const [isEmpty,setIsEmpty] = useState(false);
    //正在搜索
    const [isLoad,setIsLoading] = useState(true);
 //搜索文章
  function searchArticle(){
    if (searchKeyword=== null || searchKeyword === undefined || searchKeyword === ""){
        return;
    }
    
    ArticleClient.Search(searchKeyword,1,Constants.PageSize).then((data)=>{
        if(!data.status){
            let msg =data.message;
            if(msg === undefined || msg === null){
                msg = "系统出错啦";
            }
            toast.error(msg);
            return;
        }
        setSearchArticles(data.data.articles.map((item)=>{
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
        setIsLoading(false);
    })

}
useEffect(()=>{
    if(searchArticles===undefined){
        setIsEmpty(true);
        return;
    }
    if(searchArticles.length<=0){
        setIsEmpty(true);
    }else{
        setIsEmpty(false);
    }
},[searchArticles])
    //初始化函数
    useEffect(()=>{
        //及时搜索文章
        searchArticle();
    },[])
  return (
    <div className=" w-full h-full">

        <ToastContainer  />
    {/* header信息 */}
   <Header/>
    {/**搜索内容主体 */}
    {isLoad?<div className='w-full h-full flex justify-center items-center'><Spin isSpin={isLoad} className=" w-20 h-20"/></div>:<div className='flex justify-center items-center'>
    {!isEmpty?<>
        <div className=' w-1/5 h-full'></div>
        <div className='w-3/5 h-full pt-12'>
        <div className=" w-full mt-8">
   {searchArticles.map((item,index)=>(<div className={`px-2 hover:shadow-lg  transition duration-500 ease-in-out hover:-translate-y-1 hover:scale-105  my-3 min-h-32  border-2 w-full flex  justify-between rounded-md`} key={"newArticle"+index}>
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
        </div>
        <div className=' w-1/5 h-full'></div>
       </> :<div className=' h-screen flex justify-center items-center'><Empty/></div>}
    </div>}
    </div>
  );
};


export default SearchPage;