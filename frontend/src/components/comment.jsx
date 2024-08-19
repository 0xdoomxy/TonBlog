import {useCallback, useState} from "react";

import {CommentClient} from "../agent/agent";
import "../css/comment.css";
import {toast} from "react-toastify";
import {Button} from "antd";

const CustomerComment = ({topID, articleId, callBack, setIsOpen}) => {
    const parentInfo = {
        articleid: Number(articleId),
        topid: topID,
    };
    const [content, setContent] = useState("");
    const createComment = useCallback(() => {
            CommentClient.CreateComment({...parentInfo, content: content}).then((res) => {
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
                // TODO need to callback to refresh the comment
                callBack({
                    data: parentInfo,
                    children: []
                });
                setContent("");
                setIsOpen(-1);
            })
        }
        , [parentInfo, callBack]);

    return (<div className=" flex flex-col w-full justify-center items-center">
        <div className="comment">
            <div className=' flex flex-col '>
                <textarea value={content} onChange={(value) => {
                    setContent(value.target.value)
                }} id={"message" + articleId + "_" + topID} rows="4"
                          className="block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                          placeholder="Write your thoughts here..."></textarea>
            </div>
            <div className="w-full mt-1 flex justify-end">
                <Button style={{backgroundColor: 'rgb(0, 152, 234)'}}
                        className="hover:shadow-lg transition duration-500 ease-in-out  hover:-translate-y-1 hover:scale-105 text-white rounded-full"
                        onClick={() => {
                            createComment();
                        }}>提交</Button>
                <Button onClick={() => setIsOpen(-1)}>取消</Button></div>
        </div>
    </div>)
}


export default CustomerComment;