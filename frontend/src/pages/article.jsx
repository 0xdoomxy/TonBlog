import React,{useEffect,useState,useRef} from "react";
import { Search} from "../components";
import { useParams,useNavigate } from "react-router-dom";
import MarkdownContext from "../components/markdown";
import agent from "../agent/agent";



const Article =()=>{
    //标签颜色
    const labelColorList = ["bg-red-300","bg-yellow-200","bg-green-300","bg-pink-300","bg-gray-200"];
    const [article,setArticle] = useState({tags:[],isLike:false});
    const navigate=useNavigate();
    const navItems=[{
        Name:"Home",
        Target:"/"
    },{
        Name:"About",
        Target:"/about"
    },{Name:"Archieve",Target:"/archieve"}]
    //是否已经登陆
    const [isLogin,setIsLogin] = useState(false);
        //是否需要更换header显示
        const [changeHeader,setChangeHeader]=useState(false);
    //文章唯一id
    const{articleId} =useParams();
    // //markdown文章内容显示ref
    // const [contextDom,setContextDom] =useRef([]);
         //小屏幕点击事件，用来显示菜单栏
     const [showSmallNav,setShowSmallNav]=useState(false);
    function setAsLike(){
        agent.Like.Add(articleId,1).then((res)=>{
            if(article==undefined||article==null){
                return;
            }
            if(article.isLike){
                alert("不要重复点赞");
            }
            if (!res.status){
                alert("点赞失败");
            }
            setArticle((old) => ({...old, isLike: true}));
        });
    }
    function cancelLike(){
        if(article==undefined||article==null){
            return;
        }
        if(!article.isLike){
            alert("不要重复取消点赞");
        }
        agent.Like.Remove(articleId,1).then((res)=>{
            if (!res.status){
                alert("取消点赞失败");
            }
            setArticle((old) => ({...old, isLike: false}));
        });
    }
    function existLike(){
        agent.Like.Find(articleId,1).then((res)=>{
            if (!res.status){
                alert("查询失败");
            }
            if (res==undefined||res==null){
                return;
            }
            setArticle((old) => ({...old, isLike: res.data.exist}));
        });
    }
    function findArticle(){
        agent.Article.Find(articleId).then((res)=>{
            if (!res.status){
                alert("查询失败");
            }
            if (res==undefined||res==null){
                return;
            }
            let item  = res.data;
            item.tags = item.tags.split(",");
            item.create_time  =new Date(item.create_time).toLocaleDateString("zh-CN", {timeZone: "Asia/Shanghai", year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'});
            setArticle(item);
        });
    }
      //组件初始化的时候执行的函数
    useEffect(()=>{
        //初始化文章信息
        findArticle();
        //是否已经点赞
        existLike();
        //** 滚动时出现搜索框 */
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
// 监听是否登陆，如果登陆加载评论信息
useEffect(()=>{
    if(isLogin){
        console.log("已经登陆");
    }
},[isLogin])
    return (
        <div className=" w-full h-full">
            {/* header for search */}
            <div className=" fixed z-10 w-full ">
               <div className="bg-slate-50 w-full border-b-2 h-12 flex justify-evenly md:justify-center items-center ">
                 {!changeHeader&&(<><div  className=" w-1/4 flex justify-center   items-center py-2">
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
                </div></>)}
                {/* 小屏幕显示 */}
                <div className=" flex  pl-12 justify-center items-center  w-1/3 md:hidden ">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6 cursor-pointer" onClick={()=>{setShowSmallNav(!showSmallNav)}}>
<path strokeLinecap="round" strokeLinejoin="round" d="M3.75 5.25h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5" />
</svg>
                </div>
                {changeHeader&&<Search/>}   
                {/* TODO 点击事件 */}
            </div>
            {showSmallNav&&<div className="  bg-gray-50 border-b  z-10 relative flex w-full md:hidden flex-col justify-center items-center">
                {navItems.map((item,index)=>(
                            <div onClick={()=>{navigate(item.Target)}} className="w-full border-y hover:decoration-sky-700 hover:underline  text-center text-lg px-8 cursor-pointer " key={"smallnav"+index}>{item.Name}</div>
                        ))}
                </div>} 
                </div>
                {/* body */}
        <div className="w-full h-full pt-20 flex items-start ">
            <div className=" w-1/6"></div>
            <div className=" w-2/3 h-full">
                {/* 简介 */}
                <div className=" flex justify-between w-full h-40">
                    <div className="w-3/4 flex items-start flex-col">
                        <div className=" text-6xl font-normal text-ellipsis">{article.title}</div>
                        <div className=" flex justify-start items-center py-4 ">{article.tags.map((item,index)=>{
                            return (<div key={"tag"+index}  className={`mx-2 md:w-20   border flex justify-center items-center ${labelColorList[index%labelColorList.length]}`} >item</div>)
                        })}</div>
                        <div className=" text-xl font-serif py-1">{article.creator}</div>
                        <div className=" text-base font-sans ">{article.create_time}</div>
                    </div>
                    <div className="w-1/4 h-full flex justify-center flex-col">
                        <div className="h-1/2 border-x-2 border-t-2 text-sm  md:text-lg  w-full  font-serif flex items-center justify-center">
                            浏览量:{article.access_num}
                        </div>
                        <div className=" h-1/2 border-2 w-full text-sm  md:text-lg   font-serif flex items-center justify-center">
                            点赞量:{article.like_num}
                        </div>
                    </div>
                </div>
                <div className=" w-full h-full pt-20">
              <MarkdownContext context={article.content}/>
        </div>
        <div className=" w-full pt-32 pb-4 flex justify-end  ">
        {!isLogin?<div className="w-1/3 flex justify-end items-center  ">
                        <p className=" cursor-pointer flex w-1/2 text-xl md:text-2xl border-2 rounded-xl justify-center bg-gray-100 " onClick={()=>{setIsLogin(true)}} >登录</p>
                        </div> :<div  className=" w-1/3 flex flex-row justify-end items-center">
                            
                            <div className=" px-2 cursor-pointer ">
                            {!article.isLike?<svg onClick={()=>{setAsLike()}} xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
  <path strokeLinecap="round" strokeLinejoin="round" d="M6.633 10.25c.806 0 1.533-.446 2.031-1.08a9.041 9.041 0 0 1 2.861-2.4c.723-.384 1.35-.956 1.653-1.715a4.498 4.498 0 0 0 .322-1.672V2.75a.75.75 0 0 1 .75-.75 2.25 2.25 0 0 1 2.25 2.25c0 1.152-.26 2.243-.723 3.218-.266.558.107 1.282.725 1.282m0 0h3.126c1.026 0 1.945.694 2.054 1.715.045.422.068.85.068 1.285a11.95 11.95 0 0 1-2.649 7.521c-.388.482-.987.729-1.605.729H13.48c-.483 0-.964-.078-1.423-.23l-3.114-1.04a4.501 4.501 0 0 0-1.423-.23H5.904m10.598-9.75H14.25M5.904 18.5c.083.205.173.405.27.602.197.4-.078.898-.523.898h-.908c-.889 0-1.713-.518-1.972-1.368a12 12 0 0 1-.521-3.507c0-1.553.295-3.036.831-4.398C3.387 9.953 4.167 9.5 5 9.5h1.053c.472 0 .745.556.5.96a8.958 8.958 0 0 0-1.302 4.665c0 1.194.232 2.333.654 3.375Z" />
</svg>:<svg onClick={()=>{cancelLike()}} xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-6">
  <path d="M7.493 18.5c-.425 0-.82-.236-.975-.632A7.48 7.48 0 0 1 6 15.125c0-1.75.599-3.358 1.602-4.634.151-.192.373-.309.6-.397.473-.183.89-.514 1.212-.924a9.042 9.042 0 0 1 2.861-2.4c.723-.384 1.35-.956 1.653-1.715a4.498 4.498 0 0 0 .322-1.672V2.75A.75.75 0 0 1 15 2a2.25 2.25 0 0 1 2.25 2.25c0 1.152-.26 2.243-.723 3.218-.266.558.107 1.282.725 1.282h3.126c1.026 0 1.945.694 2.054 1.715.045.422.068.85.068 1.285a11.95 11.95 0 0 1-2.649 7.521c-.388.482-.987.729-1.605.729H14.23c-.483 0-.964-.078-1.423-.23l-3.114-1.04a4.501 4.501 0 0 0-1.423-.23h-.777ZM2.331 10.727a11.969 11.969 0 0 0-.831 4.398 12 12 0 0 0 .52 3.507C2.28 19.482 3.105 20 3.994 20H4.9c.445 0 .72-.498.523-.898a8.963 8.963 0 0 1-.924-3.977c0-1.708.476-3.305 1.302-4.666.245-.403-.028-.959-.5-.959H4.25c-.832 0-1.612.453-1.918 1.227Z" />
</svg>
}
                            </div>
                            <div className=" px-2 cursor-pointer" onClick={()=>{console.log("触发打赏事件")}} >
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
  <path strokeLinecap="round" strokeLinejoin="round" d="M21 11.25v8.25a1.5 1.5 0 0 1-1.5 1.5H5.25a1.5 1.5 0 0 1-1.5-1.5v-8.25M12 4.875A2.625 2.625 0 1 0 9.375 7.5H12m0-2.625V7.5m0-2.625A2.625 2.625 0 1 1 14.625 7.5H12m0 0V21m-8.625-9.75h18c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125h-18c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125Z" />
</svg>

                            </div>
                            <div className=" px-2 cursor-pointer" onClick={()=>{console.log("触发收藏事件")}}>
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
  <path strokeLinecap="round" strokeLinejoin="round" d="M17.593 3.322c1.1.128 1.907 1.077 1.907 2.185V21L12 17.25 4.5 21V5.507c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0 1 11.186 0Z" />
</svg>

                            </div>
                            </div>}
               
        </div>
        {isLogin&&<div className="w-full flex flex-row">
            <input  type="text" className="w-full  h-32 border-2 rounded-xl" placeholder="评论"></input>
            </div>}
        </div>
        <div className=" w-1/6"></div>
        </div>
        </div>
    )
}
export default Article;