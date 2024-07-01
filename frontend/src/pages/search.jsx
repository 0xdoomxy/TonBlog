import React, { useEffect,useState } from 'react';
import agent from "../agent/agent";
import Constants from "../util/constants";
import { useNavigate,useSearchParams  } from "react-router-dom";
import { Spin } from '../components';
const Search = () => {
    const [params] = useSearchParams()
    const navigate = useNavigate();
    const labelColorList = ["bg-red-300","bg-yellow-200","bg-green-300","bg-pink-300","bg-gray-200"]
     const navItems=[{
            Name:"Home",
            Target:"/"
        },{
            Name:"About",
            Target:"/about"
        },{Name:"Archieve",Target:"/archieve"}]
    //是否需要更换header显示
    const [changeHeader,setChangeHeader]=useState(false);
    //小屏幕点击事件，用来显示菜单栏
    const [showSmallNav,setShowSmallNav]=useState(false);
    const [searchArticles,setSearchArticles] = useState(undefined);
    const [searchKeyword,setSearchKeyword] =useState(params.get("keyword"));
    //正在搜索
    const [isLoad,setIsLoading] = useState(true);
 //搜索文章
  function searchArticle(){
    if (searchKeyword== null || searchKeyword == undefined || searchKeyword === ""){
        return;
    }
    agent.Article.Search(searchKeyword,1,Constants.PageSize).then((data)=>{
        if(!data.status){
            alert("failed",data.message);
            return;
        }
        setSearchArticles(data.data.articles.map((item)=>{
            item.tags = item.tags.split(",");
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
    //初始化函数
    useEffect(()=>{
        //及时搜索文章
        searchArticle();
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
    <div className=" w-full h-full">
    {/* header信息 */}
    <div className="w-full fixed z-10 ">
    <div className="   bg-slate-50 w-full border-b-2 h-12 flex justify-evenly md:justify-center items-center ">
        {!changeHeader&&(<div className="w-full h-full flex items-center justify-center"><div  className=" w-1/4 flex justify-center   items-center py-2">
        <h1 className=" flex align-middle font-serif text-wrap h-full text-xl md:text-3xl cursor-pointer pl-2 "  onClick={()=>{window.location.href="https://github.com/0xdoomxy"}}>0xdoomxy</h1>
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
        </div></div>)}
        {changeHeader&&<Search  onKeyDown={searchArticle}/>}
        {/* 小屏幕显示 */}
        <div className=" flex  pl-12 justify-center items-center  w-1/3 md:hidden ">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6 cursor-pointer" onClick={()=>{setShowSmallNav(!showSmallNav)}}>
<path strokeLinecap="round" strokeLinejoin="round" d="M3.75 5.25h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5" />
</svg>
        </div>
    </div>
    </div>
    {/**搜索内容主体 */}
    {isLoad?<div className='w-full h-full flex justify-center items-center'><Spin isSpin={isLoad} className=" w-20 h-20"/></div>:<div className='flex justify-center items-center'>
        <div className=' w-1/5 h-full'></div>
        <div className='w-3/5 h-full pt-12'>

        <div className=" w-full mt-8">
    {searchArticles.map((item,index)=>(<div className={`px-2 hover:shadow-lg  transition duration-500 ease-in-out hover:-translate-y-1 hover:scale-105  my-3 min-h-32  border-2 w-full flex  justify-between rounded-md`} key={"newArticle"+index}>
        <div className="flex w-2/3 flex-col justify-center">
        <p className=" font-serif md:text-2xl py-1">{item.title}</p>
        <div className=" flex py-1">
        {item.tags!=null &&item.tags.length>0&&item.tags.map((tag,index)=>(<div key={"tag"+index} className={"md:min-w-16 w-16 min-h-5  font-semibold items-center flex justify-center mx-1 "+labelColorList[index%labelColorList.length]+" text-xs rounded-lg"}>{tag}</div>) )}
        </div>
        <div className=" font-normal text-md">{item.creator}</div>
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
    </div>}
    </div>
  );
};


export default Search;