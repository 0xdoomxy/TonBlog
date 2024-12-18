import React, {useEffect, useState} from "react";

import Constants from "../util/constants";
import {useNavigate} from "react-router-dom";
import {TagClient, ArticleClient} from "../agent/agent";
import {toast} from 'react-toastify';
import {Tag, BackTop} from "antd";
import {Header} from "../components";

const HomePage = () => {
    const defaultTagsViewNum = 12;
    const navigate = useNavigate();
    const labelColorList = ["blue", "purple", "cyan", "green", "magenta", "pink", "red", "orange", "yellow", "volcano", "geekblue", "lime", "gold"];
    //所有的标签列表
    const [allTags, setAllTags] = useState([]);
    //当前可见的标签数量
    const [curTagViewNum, setCurTagViewNum] = useState(defaultTagsViewNum);
    //所有可见的标签列表    
    const [openAllTags, setOpenAllTags] = useState(false);
    //可见标签列表
    // const [tags,setTags]=useState([]);
    //最新文章列表
    const [newArticles, setNewArticles] = useState([]);
    //热门文章列表
    const [hotAriticles, setHotArticles] = useState([]);

    //获取所有标签
    function getAllTags() {
        TagClient.GetAllTags().then((data) => {
            if (data === undefined || data === null) {
                return;
            }
            if (!data.status) {
                let msg = data.message;
                if (msg === undefined || msg === null) {
                    msg = "系统出错啦";
                }
                toast.error(msg);
                return;
            }
            setAllTags(data.data);
        })
    }

    //获取最新文章
    function findTheNewestArticle() {
        ArticleClient.FindNewest(1, Constants.PageSize).then((data) => {
            if (data === undefined || data === null) {
                return;
            }
            if (!data.status) {
                let msg = data.message;
                if (msg === undefined || msg === null) {
                    msg = "系统出错啦";
                }
                toast.error(msg);
                return;
            }
            if (data.data === undefined || data.data === null || data.data.articles === undefined || data.data.articles === null) {
                return;
            }
            setNewArticles(data.data.articles.map((item) => {

                if (item.tags !== "") {
                    item.tags = item.tags.split(",");
                } else {
                    item.tags = [];
                }
                item.create_time = new Date(item.create_time).toLocaleDateString("zh-CN", {
                    timeZone: "Asia/Shanghai",
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                    second: '2-digit'
                });
                return item;
            }));
        })
    }

    /**
     * 获取热度最高文章
     */
    function findTheHotestAritcle() {
        ArticleClient.FindMaxAccess(1, Constants.PageSize).then((data) => {
            if (data === undefined || data === null) {
                return;
            }
            if (!data.status) {
                let msg = data.message;
                if (msg === undefined || msg === null) {
                    msg = "系统出错啦";
                }
                toast.error(msg);
                return;
            }
            if (data.data === undefined || data.data === null || data.data.articles === undefined || data.data.articles === null) {
                return;
            }
            setHotArticles(data.data.articles.map((item) => {
                if (item.tags !== "") {
                    item.tags = item.tags.split(",");
                } else {
                    item.tags = [];
                }
                item.create_time = new Date(item.create_time).toLocaleDateString("zh-CN", {
                    timeZone: "Asia/Shanghai", year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                    second: '2-digit'
                })
                return item;
            }));
        })
    }

    //组件初始化的时候执行的函数
    useEffect(() => {
        //初始化热门文章
        findTheHotestAritcle();
        //初始化最新文章
        findTheNewestArticle();
        //初始化所有标签
        getAllTags();
    }, [])
    const showAllTagsOnclick = () => {
        if (openAllTags) {
            setCurTagViewNum(defaultTagsViewNum);
        } else {
            setCurTagViewNum(allTags.length);
        }
        setOpenAllTags(!openAllTags);
    }
    return (
        <div className=" w-full h-full">
            <BackTop/>
            <Header/>
            <div  className=" w-full  flex justify-center items-center  pt-24">
                <div  className="w-4/5 pt-8   ">
                    <div className="grid grid-flow-row grid-cols-3  md:grid-cols-6 gap-10">
                        {allTags.map((item, index) => {
                            if (index >= curTagViewNum) {
                                return;
                            }
                            return (
                                <div
                                    className="  my-3 min-h-8 h-8  min-w-24 max-w-28 justify-center rounded-xl flex  text-center text-lg cursor-pointer  "
                                    onClick={() => {
                                        navigate(`/articles/tag?tag=${item.Name}`)
                                    }} key={"tag" + item.Name}>
                                    <p style={{
                                        color: "#222222",
                                        fontFamily: "Basel,sans-serif"
                                    }} className=" lg:text-md hover:text-blue-500 text-xl ">{item.Name}</p>
                                    <div
                                        className="min-w-6 min-h-3 h-4 bg-stone-200  text-xs rounded-lg ">{item.ArticleNum}</div>
                                </div>
                            )
                        })}
                    </div>
                    {defaultTagsViewNum<allTags.length&&<div className="w-full flex justify-end cursor-pointer" onClick={showAllTagsOnclick}>
                        { openAllTags ?
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                                 stroke="currentColor" className="size-6">
                                <path strokeLinecap="round" strokeLinejoin="round" d="m4.5 15.75 7.5-7.5 7.5 7.5"/>
                            </svg>
                            : <div className="w-full flex justify-end cursor-pointer" onClick={showAllTagsOnclick}>
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                     strokeWidth={1.5} stroke="currentColor" className="size-6">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5"/>
                                </svg>
                            </div>}
                    </div>}
                    <div className=" items-center justify-evenly flex flex-col ">
                        <div className=" w-full mt-8">
                            <div className="w-full flex justify-between items-center py-6 ">
                                <p className=" font-serif font-semibold text-3xl  text-center">最新文章</p>
                                <div className=" cursor-pointer hover:translate-x-2 duration-500 transition-transform"
                                     onClick={() => {
                                         navigate("/article/newest")
                                     }}>
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                         strokeWidth={2} stroke="currentColor" className="size-7">
                                        <path strokeLinecap="round" strokeLinejoin="round"
                                              d="m8.25 4.5 7.5 7.5-7.5 7.5"/>
                                    </svg>

                                </div>
                            </div>
                            {newArticles.map((item, index) => (<div
                                className={`px-2 hover:shadow-lg  transition duration-500 ease-in-out hover:-translate-y-1 hover:scale-105  my-3 min-h-32  border-2 w-full flex  justify-between rounded-md`}
                                key={"newArticle" + index}>
                                <div className="flex w-2/3 flex-col justify-center">
                                    <p className=" font-serif lg:text-3xl text-2xl py-1">{item.title}</p>
                                    <div className=" flex py-1">
                                        {item.tags !== null && item.tags instanceof Array && item.tags.length > 0 && item.tags.map((tag, index) => (
                                            <Tag className=" text-center text-xl lg:text-2xl"
                                                 color={labelColorList[index % labelColorList.length]}>{tag}</Tag>))}
                                    </div>
                                    <div className=" font-normal text-md lg:text-xl truncate">{item.creator_name}</div>
                                    <div className=" font-normal text-md lg:text-xl">{item.create_time}</div>
                                </div>
                                <div className=" flex justify-center w-1/3 items-center flex-col">
                                    <button style={{
                                        color: "#222222",
                                        fontFamily: "Basel,sans-serif"
                                    }}  className=" w-20 h-12 border-2 rounded-xl text-md lg:text-xl hover:bg-blue-100"
                                            onClick={() => navigate("/article/" + item.id)}>阅读
                                    </button>
                                    <div className=" font-serif text-ellipsis text-md lg:text-xl">浏览量:{item.access_num}</div>
                                </div>
                            </div>))}

                        </div>
                        <div className=" w-full mt-8">
                            <div className="w-full flex justify-between items-center py-6 ">
                                <p className=" font-serif font-semibold text-3xl  text-center">热门文章</p>
                                <div className=" cursor-pointer hover:translate-x-2 duration-500 transition-transform"
                                     onClick={() => {
                                         navigate("/article/hot")
                                     }}>
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                         strokeWidth={2} stroke="currentColor" className="size-7">
                                        <path strokeLinecap="round" strokeLinejoin="round"
                                              d="m8.25 4.5 7.5 7.5-7.5 7.5"/>
                                    </svg>

                                </div>
                            </div>
                            {hotAriticles.map((item, index) => (<div
                                className="  px-2 hover:shadow-lg transition duration-500 ease-in-out  hover:-translate-y-1 hover:scale-105  my-3 min-h-32  border-2 flex  justify-between rounded-md"
                                key={"newArticle" + index}>
                                <div className="flex flex-col w-2/3 justify-center ">
                                    <p className=" font-serif lg:text-3xl text-2xl py-1">{item.title}</p>
                                    <div className=" flex py-1">
                                        {item.tags != null && item.tags instanceof Array && item.tags.length > 0 && item.tags.map((tag, index) => (
                                            <Tag className={"text-xl lg:text-2xl"} color={labelColorList[index % labelColorList.length]}>{tag}</Tag>))}
                                    </div>
                                    <div className=" font-normal text-md lg:text-xl truncate">{item.creator_name}</div>
                                    <div className=" font-normal text-md lg:text-xl">{item.create_time}</div>
                                </div>
                                <div className=" flex justify-center w-1/3 items-center flex-col">
                                    <button style={{
                                        color: "#222222",
                                        fontFamily: "Basel,sans-serif"
                                    }} className=" text-md lg:text-xl w-20 h-12 border-2 rounded-xl hover:bg-blue-100"
                                            onClick={() => navigate("/article/" + item.id)}>阅读
                                    </button>
                                    <div className=" font-serif text-ellipsis text-md lg:text-xl">浏览量:{item.access_num}</div>
                                </div>
                            </div>))}

                        </div>
                    </div>
                </div>
            </div>
            <div className=" flex w-full border-t-2  bg-slate-50 justify-center items-center ">
                <div className=" w-1/5"></div>
                <div className=" h-20  w-full flex justify-around items-center">
                    <div className=" w-1/2 text-md md:pl-10">© 0xdoomxy 保留所有权利</div>
                    <div className=" w-1/2 flex justify-end items-center md:pr-32">
                        <div className="px-2 cursor-pointer" onClick={() => {
                            window.location.href = "https://github.com/0xdoomxy"
                        }}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="1.5em" height="1.5em" viewBox="0 0 16 16">
                                <path fill="currentColor"
                                      d="M8 0c4.42 0 8 3.58 8 8a8.01 8.01 0 0 1-5.45 7.59c-.4.08-.55-.17-.55-.38c0-.27.01-1.13.01-2.2c0-.75-.25-1.23-.54-1.48c1.78-.2 3.65-.88 3.65-3.95c0-.88-.31-1.59-.82-2.15c.08-.2.36-1.02-.08-2.12c0 0-.67-.22-2.2.82c-.64-.18-1.32-.27-2-.27s-1.36.09-2 .27c-1.53-1.03-2.2-.82-2.2-.82c-.44 1.1-.16 1.92-.08 2.12c-.51.56-.82 1.28-.82 2.15c0 3.06 1.86 3.75 3.64 3.95c-.23.2-.44.55-.51 1.07c-.46.21-1.61.55-2.33-.66c-.15-.24-.6-.83-1.23-.82c-.67.01-.27.38.01.53c.34.19.73.9.82 1.13c.16.45.68 1.31 2.69.94c0 .67.01 1.3.01 1.49c0 .21-.15.45-.55.38A7.995 7.995 0 0 1 0 8c0-4.42 3.58-8 8-8"></path>
                            </svg>
                        </div>
                        <div className=" px-2 cursor-pointer ">
                            <svg xmlns="http://www.w3.org/2000/svg" width="1.5em" height="1.5em" viewBox="0 0 32 32">
                                <path fill="currentColor"
                                      d="M7.845 9.983L9.88 27.336c0 .977 2.74 1.77 6.12 1.77s6.12-.793 6.12-1.77L24.5 9.85c-2.455 1.024-6.812 1.134-8.498 1.134c-1.61 0-5.655-.1-8.155-1zm16.285-4.23l-.376-1.68c0-.65-3.472-1.178-7.754-1.178s-7.754.53-7.754 1.18L7.87 5.752c-.714.284-1.12.608-1.12.953V7.99c0 1.1 4.142 1.994 9.25 1.994s9.25-.894 9.25-1.995V6.704c0-.345-.406-.67-1.12-.953z"></path>
                            </svg>
                        </div>
                    </div>
                </div>
                <div className=" w-1/5"></div>
            </div>
        </div>
    )

}

export default HomePage;