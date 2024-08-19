import React, {useEffect, useState} from "react";
import {CustomerComment, Header} from "../components";
import {useParams} from "react-router-dom";
import MarkdownContext from "../components/markdown";
import {LikeClient, ArticleClient, CommentClient, Authorization} from "../agent/agent";
import {useTonWallet, useTonConnectUI} from "@tonconnect/ui-react";
import {Tag, Modal, Input, InputNumber, Segmented, Button, BackTop, Avatar, Tooltip} from "antd";
import {Comment} from "@ant-design/compatible";
import {toast} from "react-toastify";
import {FormatDateDistance} from "../util/time";
import {UserOutlined, MoneyCollectOutlined} from "@ant-design/icons";

const ArticlePage = () => {
    //文章唯一id
    const {articleId} = useParams();
    //标签颜色
    const labelColorList = ["blue", "purple", "cyan", "green", "magenta", "pink", "red", "orange", "yellow", "volcano", "geekblue", "lime", "gold"];
    const [article, setArticle] = useState({tags: [], isLike: false});
    const [comments, setComments] = useState(new Map());
    //正在打开的评论栏
    const [isOpen, setIsOpen] = useState(-1);
    const wallet = useTonWallet();
    const [tonConnectUI] = useTonConnectUI();
    //是否正在打赏中
    const [rewardModal, setRewardModal] = useState(false);
    //打赏价格
    const [rewardInfo, setRewardInfo] = useState({
        address: "0:9cc2ceadf8282782c3bfe6b7ad0933e59b6f7257025f3fad607106738d91dea0",
        prices: 0
    });
    const [createCommentInfo, setCreateCommentInfo] = useState({
        articleid: Number(articleId),
        topid: 0,
    })

    //reward
    function reward() {
        if (wallet === undefined || wallet === null) {
            toast.error("请先登陆");
            return;
        }
        tonConnectUI.sendTransaction({
            messages: [{
                address: rewardInfo.address,
                amount: rewardInfo.prices * 1e9,
                validUntil: Math.floor(Date.now() / 1000) + 600
            }]
        }).then((res) => {
            if (res.status) {
                toast.success("打赏成功");
            } else {
                toast.error("打赏失败");
            }
        });
    }

    //TODO创建评论
    function createComment() {
        CommentClient.CreateComment(createCommentInfo).then((res) => {
            if (res === undefined || res === null) {
                return;
            }
            if (!res.status) {
                let msg = res.message;
                if (msg === undefined || msg === null) {
                    msg = "系统出错啦";
                }
                toast.error(msg);
                return;
            }
            toast.success("评论成功");
            comments.get(createCommentInfo.topid).push({
                data: {
                    TopID: createCommentInfo.topid,
                    Content: createCommentInfo.content,
                },
                children: []
            });
            createCommentInfo.content = "";
        })

    }

    //评论显示组件生成器
    function commentView(comment) {
        return (<>
                {comment !== undefined && comment != null && comment instanceof Array && comment.map((item, index) => {
                    return (
                        <Comment
                            actions={[(
                                <div className=" w-full flex flex-col justify-center"><span
                                    key="comment-nested-reply-to" className="w-full cursor-pointer mb-2"
                                    onClick={() => setIsOpen(item.data.ID)}>回复</span>{isOpen === item.data.ID &&
                                    <CustomerComment topID={item.data.ID} setIsOpen={setIsOpen} articleId={articleId}
                                                     key={"sub_comment_" + index} callBack={(data) => {
                                        let children = comments.get(item.data.ID);
                                        if (children === undefined) {
                                            children = [data];
                                        } else {
                                            children.push(data);
                                        }
                                    }}/>}</div>
                            )]}
                            author={<a>{item.data.Creator}</a>}
                            avatar={
                                <Avatar
                                    src="https://zos.alipayobjects.com/rmsportal/ODTLcjxAfvqbxHnVXCYX.png"
                                    alt="Han Solo"
                                />
                            }
                            datetime={
                                <Tooltip title={FormatDateDistance(item.data.CreateAt)}>
                                <span>
                                    {FormatDateDistance(item.data.CreateAt)}
                                </span>
                                </Tooltip>
                            }
                            content={
                                <p>
                                    {item.data.Content}
                                </p>
                            }
                        >
                            {commentView(item.children)}
                        </Comment>
                    )
                })}
            </>
        )

    }

    function searchCommentByArticle() {
        CommentClient.SearchByArticle(articleId).then((res) => {
            if (res === undefined || res === null) {
                return;
            }
            if (!res.status) {
                let msg = res.message;
                if (msg === undefined || msg === null) {
                    msg = "系统出错啦";
                }
                toast.error(msg);
                return;
            }
            if (res.data === undefined || res.data === null) {
                return;
            }
            //实现解析comments,递归解析
            var allComments = res.data;
            allComments.sort((a, b) => a.ID - b.ID);
            var top = new Map();
            for (let i = 0; i < allComments.length; i++) {
                let item = {data: allComments[i]};
                if (top.has(item.data.ID)) {
                    let subItem;
                    subItem = top.get(item.ID);
                    item.children = subItem;
                } else {
                    let childs = [];
                    item.children = childs;
                    top.set(item.data.ID, childs);
                }
                if (!top.has(item.data.TopID)) {
                    let childs = [];
                    childs.push(item);
                    top.set(item.data.TopID, childs);
                } else {
                    top.get(item.data.TopID).push(item);
                }
            }
            setComments(top);
        })
    }

    function setAsLike() {
        LikeClient.Add(articleId, 1).then((res) => {
            if (res === undefined || res === null) {
                return;
            }
            if (article === undefined || article === null) {
                return;
            }
            if (article.isLike) {
                toast.error("不要重复点赞");
            }
            if (!res.status) {
                toast.error("点赞失败");
            }
            setArticle((old) => ({...old, isLike: true, like_num: old.like_num + 1}));
        });
    }

    function cancelLike() {
        if (article === undefined || article === null) {
            return;
        }
        if (!article.isLike) {
            toast.error("不要重复取消点赞");
        }
        LikeClient.Remove(articleId, 1).then((res) => {
            if (res === undefined || res === null) {
                return;
            }
            if (!res.status) {
                toast.error("取消点赞失败");
            }
            setArticle((old) => ({...old, isLike: false, like_num: old.like_num - 1}));
        });
    }

    function existLike() {
        LikeClient.Find(articleId, 1).then((res) => {
            if (res === undefined || res === null) {
                return;
            }
            if (!res.status) {
                toast.error("查询失败");
            }
            if (res.data === undefined || res.data === null) {
                return;
            }
            setArticle((old) => ({...old, isLike: res.data.exist}));
        });
    }

    function findArticle() {
        ArticleClient.Find(articleId).then((res) => {
            if (res === undefined || res === null) {
                return;
            }
            if (!res.status) {
                toast.error("查询失败");
            }
            if (res.data === undefined || res.data === null) {
                return;
            }
            let item = res.data;
            item.tags = item.tags.split(",");
            item.create_time = new Date(item.create_time).toLocaleDateString("zh-CN", {
                timeZone: "Asia/Shanghai", year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit'
            });
            setArticle(item);
        });
    }

    //TODO this function has delay to do
    useEffect(() => {
        if (tonConnectUI.connected) {
            //TODO dangerous: this function will be invalid when verify user to login spend long time
            setTimeout(() => {
                if (Authorization !== undefined) {
                    //是否已经点赞
                    existLike();
                    //完成登录时初始化评论信息
                    searchCommentByArticle();
                }
            }, 3000);

        }
    }, [tonConnectUI.connected, Authorization])
    //组件初始化的时候执行的函数
    useEffect(() => {
        //初始化文章信息
        findArticle();
    }, [])

    return (
        <div className=" w-full h-full">
            <BackTop/>
            {/* header for search */}
            <Header/>
            {/* body */}
            <div className="w-full h-full pt-20 flex items-start ">
                <Modal title="Ton" className=" z-0" onOk={() => {
                    reward()
                }} open={rewardModal} okText="赞助" cancelText="离开" onCancel={() => {
                    setRewardModal(false)
                }}>
                    <Input className=" mt-4" prefix={<UserOutlined/>} value={rewardInfo.address} disabled></Input>
                    <InputNumber className=" my-4" value={rewardInfo.prices} onChange={(value) => {
                        setRewardInfo((old) => ({...old, prices: value}))
                    }}
                                 suffix="Ton"
                                 style={{
                                     width: '100%',
                                 }}
                                 addonBefore={<svg width="20" height="20" viewBox="0 0 56 56" fill="none"
                                                   xmlns="http://www.w3.org/2000/svg">
                                     <path
                                         d="M28 56C43.464 56 56 43.464 56 28C56 12.536 43.464 0 28 0C12.536 0 0 12.536 0 28C0 43.464 12.536 56 28 56Z"
                                         fill="#0098EA"/>
                                     <path
                                         d="M37.5603 15.6277H18.4386C14.9228 15.6277 12.6944 19.4202 14.4632 22.4861L26.2644 42.9409C27.0345 44.2765 28.9644 44.2765 29.7345 42.9409L41.5381 22.4861C43.3045 19.4251 41.0761 15.6277 37.5627 15.6277H37.5603ZM26.2548 36.8068L23.6847 31.8327L17.4833 20.7414C17.0742 20.0315 17.5795 19.1218 18.4362 19.1218H26.2524V36.8092L26.2548 36.8068ZM38.5108 20.739L32.3118 31.8351L29.7417 36.8068V19.1194H37.5579C38.4146 19.1194 38.9199 20.0291 38.5108 20.739Z"
                                         fill="white"/>
                                 </svg>
                                 }
                    />
                    <div className=" mt-4 w-full flex justify-center items-center">
                        <div className=" w-full flex justify-start">
                            <Segmented
                                options={[
                                    {label: '赞助记录', value: 'all', icon: <MoneyCollectOutlined/>},
                                    {label: '我的赞助', value: 'my', icon: <UserOutlined/>},
                                ]}
                            />
                        </div>
                    </div>
                </Modal>
                <div className=" w-1/6"></div>
                <div className=" w-2/3 h-full">
                    {/* 简介 */}
                    <div className=" flex justify-between w-full ">
                        <div className="w-3/4 flex items-start flex-col">
                            <h1 className="w-full text-6xl font-normal  max-h-32 line-clamp-2">{article.title}</h1>
                            <div className=" flex justify-start items-center py-4 ">{article.tags.map((item, index) => {
                                return (<Tag key={"tag" + index}
                                             color={labelColorList[index % labelColorList.length]}>{item}</Tag>)
                            })}</div>
                            <div className="w-full text-xl font-serif py-4  truncate "
                                 id={article.creator}>{article.creator_name}</div>
                            <div className=" text-base font-sans ">{article.create_time}</div>
                        </div>
                        <div className="w-1/4 h-48 flex justify-center flex-col">
                            <div
                                className="h-1/2 border-x-2 border-t-2 text-sm  md:text-lg  w-full  font-serif flex items-center justify-center">
                                浏览量:{article.access_num}
                            </div>
                            <div
                                className=" h-1/2 border-2 w-full text-sm  md:text-lg   font-serif flex items-center justify-center">
                                点赞量:{article.like_num}
                            </div>
                        </div>
                    </div>
                    <div className=" pt-20">
                        <MarkdownContext context={article.content}/>
                    </div>
                    <div className=" w-full pt-24 pb-4 flex justify-end  ">
                        {!tonConnectUI.connected ? <div className="w-1/3 flex justify-end items-center  ">
                            <Button style={{backgroundColor: 'rgb(0, 152, 234)'}}
                                    className=" hover:shadow-lg transition duration-500 ease-in-out  hover:-translate-y-1 hover:scale-105  rounded-full md:w-32 h-10 text-white "
                                    onClick={() => tonConnectUI.openModal()}>Conntect Wallet</Button>
                        </div> : <div className=" w-1/3 flex flex-row justify-end items-center">
                            <div className=" px-2 cursor-pointer ">
                                {!article.isLike ? <svg onClick={() => {
                                    setAsLike()
                                }} xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                                                        stroke="currentColor" className="size-6">
                                    <path strokeLinecap="round" strokeLinejoin="round"
                                          d="M6.633 10.25c.806 0 1.533-.446 2.031-1.08a9.041 9.041 0 0 1 2.861-2.4c.723-.384 1.35-.956 1.653-1.715a4.498 4.498 0 0 0 .322-1.672V2.75a.75.75 0 0 1 .75-.75 2.25 2.25 0 0 1 2.25 2.25c0 1.152-.26 2.243-.723 3.218-.266.558.107 1.282.725 1.282m0 0h3.126c1.026 0 1.945.694 2.054 1.715.045.422.068.85.068 1.285a11.95 11.95 0 0 1-2.649 7.521c-.388.482-.987.729-1.605.729H13.48c-.483 0-.964-.078-1.423-.23l-3.114-1.04a4.501 4.501 0 0 0-1.423-.23H5.904m10.598-9.75H14.25M5.904 18.5c.083.205.173.405.27.602.197.4-.078.898-.523.898h-.908c-.889 0-1.713-.518-1.972-1.368a12 12 0 0 1-.521-3.507c0-1.553.295-3.036.831-4.398C3.387 9.953 4.167 9.5 5 9.5h1.053c.472 0 .745.556.5.96a8.958 8.958 0 0 0-1.302 4.665c0 1.194.232 2.333.654 3.375Z"/>
                                </svg> : <svg onClick={() => {
                                    cancelLike()
                                }} xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"
                                              className="size-6">
                                    <path
                                        d="M7.493 18.5c-.425 0-.82-.236-.975-.632A7.48 7.48 0 0 1 6 15.125c0-1.75.599-3.358 1.602-4.634.151-.192.373-.309.6-.397.473-.183.89-.514 1.212-.924a9.042 9.042 0 0 1 2.861-2.4c.723-.384 1.35-.956 1.653-1.715a4.498 4.498 0 0 0 .322-1.672V2.75A.75.75 0 0 1 15 2a2.25 2.25 0 0 1 2.25 2.25c0 1.152-.26 2.243-.723 3.218-.266.558.107 1.282.725 1.282h3.126c1.026 0 1.945.694 2.054 1.715.045.422.068.85.068 1.285a11.95 11.95 0 0 1-2.649 7.521c-.388.482-.987.729-1.605.729H14.23c-.483 0-.964-.078-1.423-.23l-3.114-1.04a4.501 4.501 0 0 0-1.423-.23h-.777ZM2.331 10.727a11.969 11.969 0 0 0-.831 4.398 12 12 0 0 0 .52 3.507C2.28 19.482 3.105 20 3.994 20H4.9c.445 0 .72-.498.523-.898a8.963 8.963 0 0 1-.924-3.977c0-1.708.476-3.305 1.302-4.666.245-.403-.028-.959-.5-.959H4.25c-.832 0-1.612.453-1.918 1.227Z"/>
                                </svg>
                                }
                            </div>
                            <div className=" px-2 cursor-pointer" onClick={() => {
                                setRewardModal(true)
                            }}>
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                     strokeWidth={1.5} stroke="currentColor" className="size-6">
                                    <path strokeLinecap="round" strokeLinejoin="round"
                                          d="M21 11.25v8.25a1.5 1.5 0 0 1-1.5 1.5H5.25a1.5 1.5 0 0 1-1.5-1.5v-8.25M12 4.875A2.625 2.625 0 1 0 9.375 7.5H12m0-2.625V7.5m0-2.625A2.625 2.625 0 1 1 14.625 7.5H12m0 0V21m-8.625-9.75h18c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125h-18c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125Z"/>
                                </svg>

                            </div>
                        </div>}

                    </div>
                    {tonConnectUI.connected && <>
                        <div className="w-full flex flex-col">
                            <div className=" w-full flex flex-col">
                                <label htmlFor="message"
                                       className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">评论</label>
                                <textarea value={createCommentInfo.content} onChange={(value) => {
                                    setCreateCommentInfo((origin) => ({...origin, content: value.target.value}))
                                }} id="message" rows="4"
                                          className="block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                                          placeholder="Write your thoughts here..."></textarea>
                            </div>
                            <div className=" mt-1 flex justify-end">
                                <Button style={{backgroundColor: 'rgb(0, 152, 234)'}}
                                        className="hover:shadow-lg transition duration-500 ease-in-out  hover:-translate-y-1 hover:scale-105 text-white rounded-full"
                                        onClick={() => {
                                            createComment();
                                        }}>提交</Button>
                            </div>

                        </div>
                        <div className=" mt-4">
                            {commentView(comments.get(0))}
                        </div>
                    </>}
                </div>
                <div className=" w-1/6"></div>
            </div>
        </div>
    )
}
export default ArticlePage;