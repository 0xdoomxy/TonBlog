import React, {useContext, useEffect, useState} from 'react';
import {useNavigate} from "react-router-dom";
import {Search} from './search';
import {AutoComplete, Divider, List, Modal, Skeleton, Tag} from 'antd';
import InfiniteScroll from 'react-infinite-scroll-component';
import Constants from '../util/constants';
import {ArticleClient} from '../agent/agent';
import {toast} from 'react-toastify';
import {MenuToggle} from './menutoggle';
import {useMenuAnimation} from '../hooks/menuanimation';
import "../css/header.css"
import {Web3Wallet} from "../App";
import {motion} from "framer-motion";
import {formatAddress} from "../util/web3.js";
import Avatar from "./avatar.jsx";

const Header = () => {
    const {
        searchWalletModal,
        setSearchWalletModal,
        selectedWallet,
        setSelectedWallet,
        userAccount,
        setUserAccount
    } = useContext(Web3Wallet);
    const labelColorList = ["blue", "purple", "cyan", "green", "magenta", "pink", "red", "orange", "yellow", "volcano", "geekblue", "lime", "gold"];
    //是否需要更换header显示
    const [changeHeader, setChangeHeader] = useState(false);
    const navigate = useNavigate();
    //搜索框显示
    const [openSearch, setOpenSearch] = useState(false);
    //搜索的文本
    const [keyword, setKeyword] = useState("");
    //小屏幕点击事件，用来显示菜单栏
    const [showSmallNav, setShowSmallNav] = useState(false);
    const scope = useMenuAnimation(showSmallNav);
    const [total, setTotal] = useState(0);
    useEffect(() => {
        if (keyword !== null && keyword !== undefined && keyword !== "") {
            searchArticle(1, false);
        }
    }, [keyword])
    //搜索到的文章
    const [searchArticles, setSearchArticles] = useState([]);
    const navItems = [
        {
            Name: "主页",
            Target: "/"
        },
        {
            Name: '空投',
            Target: '/airport'
        },
        {
            Name: "作者简介",
            Target: "/about"
        }]

    //搜索文章(isContinue:是否是跟进页数)
    function searchArticle(page, isContinue) {
        if (keyword === null || keyword === undefined) {
            return;
        }

        ArticleClient.Search(keyword, page, Constants.PageSize).then((data) => {
            if (data === undefined || data === null) {
                toast.error("系统故障啦");
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
            if (data.data === null || data.data === undefined) {
                return;
            }
            if (!isContinue) {
                setSearchArticles(data.data.articles.map((item) => {
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
                    });
                    return item;
                }));
            } else {
                setSearchArticles((origin) => [origin, ...data.data.articles.map((item) => {
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
                    });
                    return item;
                })]);
            }
            setTotal(data.data.total);
        })
    }

    useEffect(() => {
        //监听鼠标滚动事件来改变header
        const checkScroll = () => {
            if (window.scrollY > 200) {
                setChangeHeader(true);
            } else {
                setChangeHeader(false);
            }
        };
        window.addEventListener("scroll", checkScroll);
        return () => window.removeEventListener("scroll", checkScroll);
    }, [])
    return (
        <div className="w-full fixed z-10 ">
            <Modal width="75%" onCancel={() => {
                setOpenSearch(false)
            }} closable={false} keyboard={true} footer={null} open={openSearch}>
                <div className='w-full flex flex-col justify-center items-start'>
                    <div className=' w-full flex flex-row items-center '>
                        <div className=' pr-3'>
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                                 stroke="currentColor" className="size-8">
                                <path strokeLinecap="round" strokeLinejoin="round"
                                      d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"/>
                            </svg>
                        </div>
                        <AutoComplete value={keyword} defaultOpen open={true} onChange={(event) => {
                            setKeyword(event)
                        }} backfill className=' h-10 w-full'/>
                    </div>
                    <div className=' flex justify-center w-full mt-4'>
                        <div
                            id="scrollableDiv"
                            style={{
                                height: 400,
                                overflow: 'auto',
                                padding: '0 16px',
                            }}
                            className=' w-full'
                        >
                            <InfiniteScroll className='w-full h-full overflow-y-auto '
                                            dataLength={searchArticles.length}
                                            next={() => {
                                                searchArticle(Number.parseInt(searchArticles.length / Constants.PageSize) + 1, true)
                                            }}
                                            hasMore={total > searchArticles.length && searchArticles.length < Constants.PageSize * 5}
                                            loader={
                                                <Skeleton
                                                    avatar
                                                    paragraph={{
                                                        rows: 1,
                                                    }}
                                                    active
                                                />
                                            }
                                            endMessage={<Divider plain>It is all, nothing more 🤐</Divider>}
                                            scrollableTarget="scrollableDiv"
                            >
                                <List
                                    dataSource={searchArticles}
                                    renderItem={(item) => (
                                        <List.Item className=' w-full cursor-pointer' key={item.id}
                                                   onClick={() => navigate("/article/" + item.id)}>
                                            <List.Item.Meta
                                                avatar={<div className=' flex items-center'>
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none"
                                                         viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor"
                                                         className="size-8">
                                                        <path strokeLinecap="round" strokeLinejoin="round"
                                                              d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"/>
                                                    </svg>
                                                </div>
                                                }
                                                title={item.tags != null && item.tags instanceof Array && item.tags.length > 0 && item.tags.map((tag, index) => (
                                                    <Tag
                                                        color={labelColorList[index % labelColorList.length]}>{tag}</Tag>))}
                                                description={item.title}
                                            />
                                            <div
                                                className=" font-serif text-ellipsis text-sm">浏览量:{item.access_num}</div>
                                        </List.Item>
                                    )}
                                />
                            </InfiniteScroll>
                        </div>
                    </div>
                </div>
            </Modal>
            <div
                className={`bg-slate-50  w-full border-b-2  flex justify-evenly md:justify-center items-center ${changeHeader ? "h-12" : "h-20"}`}>
                {!changeHeader && (<div className=" w-1/6 md:w-full h-full flex items-center justify-center">
                    <div className=" w-1/4 flex justify-center   items-center py-2">
                        <div className=" flex align-middle font-serif text-3xl lg:text-4xl h-full cursor-pointer lg:pl-2 pl-24 "
                             onClick={() => {
                                 window.location.href = "https://github.com/0xdoomxy"
                             }}>0xdoomxy
                        </div>
                    </div>
                    <div className=" w-1/2   hidden  md:flex justify-start items-center">
                        <div className=' w-2/3 flex flex-row justify-evenly'>
                            {navItems.map((item, index) => (
                                <div onClick={() => {
                                    navigate(item.Target)
                                }}
                                     style={{

                                         color: "#222222",
                                         fontFamily: "Basel,sans-serif"
                                     }}   className=" hover:-translate-y-1 duration-500  text-center text-xl  px-4 lg:px-8 cursor-pointer "
                                     key={"nav" + index}>{item.Name}</div>
                            ))}
                        </div>
                        <div className=" lg:pl-24 pl-6  flex justify-start ">
                            <div className=" cursor-pointer " onClick={() => {
                                setOpenSearch(true)
                            }}>
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                     strokeWidth={1.5} stroke="currentColor" className="size-6">
                                    <path strokeLinecap="round" strokeLinejoin="round"
                                          d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"/>
                                </svg>
                            </div>
                        </div>
                    </div>
                    <div
                        className=" hidden md:flex w-1/8 justify-evenly ">{userAccount != null && userAccount.length > 0 ?
                        <Avatar/> :
                        <motion.button style={{width: "150px", height: " 60px"}} whileHover={{scale: 1.1}}
                                       whileTap={{scale: 0.95}} className={"motion-button "}
                                       onClick={() => {
                                           setSearchWalletModal(!searchWalletModal)
                                       }}>Connect Wallet</motion.button>}</div>
                </div>)}
                {changeHeader &&
                    <div className="w-full h-full hidden md:flex">
                        <Search onKeyDown={(event) => {
                            if (event.keyCode !== 13) {
                                return;
                            }
                            if (event.target.value === undefined) {
                                return
                            }
                            navigate(`/search?keyword=${event.target.value}`,)
                        }}/>
                    </div>
                }
                {/* 小屏幕显示 */}
                <div className="w-full h-2/3 md:hidden">
                    <Search onKeyDown={(event) => {
                        if (event.keyCode !== 13) {
                            return;
                        }
                        if (event.target.value === undefined) {
                            return
                        }
                        navigate(`/search?keyword=${event.target.value}`,)
                    }}/>
                </div>
                <div ref={scope} className=" flex  md:pr-4 lg:pl-12  justify-center items-center  w-1/3 md:hidden ">
                    <nav className="menu">
                        <ul className={"flex items-center justify-center flex-col"}>
                            {navItems.map((item, index) => (
                                <li onClick={() => {
                                    navigate(item.Target)
                                }}
                                    key={"smallnav" + index}>{item.Name}</li>
                            ))}
                        </ul>
                    </nav>
                    <MenuToggle toggle={() => setShowSmallNav(!showSmallNav)}/>
                </div>
            </div>
        </div>
    )
}

export default Header;