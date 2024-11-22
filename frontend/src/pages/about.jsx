import React from "react";
import {Timeline} from "antd";
import avatar from "../asserts/avatar.jpg";
import {Header} from "../components";

const AboutPage = () => {
    return (
        <div className="w-full h-full ">
            <Header/>
            <div style={{height: "85%"}} className="w-full flex pt-12  bg-slate-100">
                {/* 左边 */}
                <div
                    className=" h-full md:w-1/3 w-1/2 flex flex-col items-center justify-center border-r-2 border-gray-200">
                    <div
                        className="md:h-[200px] md:w-[200px] h-[100px] w-[100px] rounded-full border-4 border-white dark:border-dark-3">
                        <img
                            src={avatar}
                            alt="avatar"
                            className="h-full w-full rounded-full object-cover object-center"
                        />
                    </div>
                    <div className=" w-full flex justify-center items-center ">
                        <div className=" cursor-pointer px-2" onClick={() => {
                            window.location.href = "https://github.com/0xdoomxy"
                        }}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="1.5em" height="1.5em" viewBox="0 0 16 16">
                                <path fill="currentColor"
                                      d="M8 0c4.42 0 8 3.58 8 8a8.01 8.01 0 0 1-5.45 7.59c-.4.08-.55-.17-.55-.38c0-.27.01-1.13.01-2.2c0-.75-.25-1.23-.54-1.48c1.78-.2 3.65-.88 3.65-3.95c0-.88-.31-1.59-.82-2.15c.08-.2.36-1.02-.08-2.12c0 0-.67-.22-2.2.82c-.64-.18-1.32-.27-2-.27s-1.36.09-2 .27c-1.53-1.03-2.2-.82-2.2-.82c-.44 1.1-.16 1.92-.08 2.12c-.51.56-.82 1.28-.82 2.15c0 3.06 1.86 3.75 3.64 3.95c-.23.2-.44.55-.51 1.07c-.46.21-1.61.55-2.33-.66c-.15-.24-.6-.83-1.23-.82c-.67.01-.27.38.01.53c.34.19.73.9.82 1.13c.16.45.68 1.31 2.69.94c0 .67.01 1.3.01 1.49c0 .21-.15.45-.55.38A7.995 7.995 0 0 1 0 8c0-4.42 3.58-8 8-8"></path>
                            </svg>
                        </div>
                        <div className="px-2 ">
                            <svg xmlns="http://www.w3.org/2000/svg" width="1.5em" height="1.5em" viewBox="0 0 16 16">
                                <path fill="currentColor"
                                      d="M8 .25a.75.75 0 0 1 .673.418l1.882 3.815l4.21.612a.75.75 0 0 1 .416 1.279l-3.046 2.97l.719 4.192a.751.751 0 0 1-1.088.791L8 12.347l-3.766 1.98a.75.75 0 0 1-1.088-.79l.72-4.194L.818 6.374a.75.75 0 0 1 .416-1.28l4.21-.611L7.327.668A.75.75 0 0 1 8 .25"></path>
                            </svg>
                        </div>
                        <div className=" px-2">
                            <svg xmlns="http://www.w3.org/2000/svg" width="1.5em" height="1.5em" viewBox="0 0 512 512">
                                <path fill="currentColor"
                                      d="M408.67 298.53a21 21 0 1 1 20.9-21a20.85 20.85 0 0 1-20.9 21m-102.17 0a21 21 0 1 1 20.9-21a20.84 20.84 0 0 1-20.9 21m152.09 118.86C491.1 394.08 512 359.13 512 319.51c0-71.08-68.5-129.35-154.41-129.35s-154.42 58.27-154.42 129.35s68.5 129.34 154.42 129.34c17.41 0 34.83-2.33 49.92-7c2.49-.86 3.48-1.17 4.64-1.17a16.67 16.67 0 0 1 8.13 2.34L454 462.83a11.6 11.6 0 0 0 3.48 1.17a5 5 0 0 0 4.65-4.66a14.3 14.3 0 0 0-.77-3.86c-.41-1.46-5-16-7.36-25.27a19 19 0 0 1-.33-3.47a11.4 11.4 0 0 1 5-9.35"></path>
                                <path fill="currentColor"
                                      d="M246.13 178.51a24.47 24.47 0 0 1 0-48.94c12.77 0 24.38 11.65 24.38 24.47c1.16 12.82-10.45 24.47-24.38 24.47m-123.06 0A24.47 24.47 0 1 1 147.45 154a24.57 24.57 0 0 1-24.38 24.47M184.6 48C82.43 48 0 116.75 0 203c0 46.61 24.38 88.56 63.85 116.53C67.34 321.84 68 327 68 329a11.4 11.4 0 0 1-.66 4.49C63.85 345.14 59.4 364 59.21 365s-1.16 3.5-1.16 4.66a5.49 5.49 0 0 0 5.8 5.83a7.15 7.15 0 0 0 3.49-1.17L108 351c3.49-2.33 5.81-2.33 9.29-2.33a16.3 16.3 0 0 1 5.81 1.16c18.57 5.83 39.47 8.16 60.37 8.16h10.45a133.2 133.2 0 0 1-5.81-38.45c0-78.08 75.47-141 168.35-141h10.45C354.1 105.1 277.48 48 184.6 48"></path>
                            </svg>
                        </div>
                    </div>
                    <div className=" w-full text-center pt-1 ">
                        <p className=" text-sm indent-4 font-normal">We can only see a short distance ahead, but we can
                            see plenty there that needs to be done.</p>
                    </div>
                </div>
                {/* 右边 */}
                <div  className="w-full h-full pl-6 pt-6  flex flex-col">
                    <div className="w-full h-1/3 flex flex-col">
                        <h1 className=" font-serif font-semibold md:text-xl pb-2">简&nbsp;&nbsp;&nbsp;&nbsp;介</h1>
                        <p className=" font-serif indent-4 md:text-lg text-sm line-clamp-8  md:max-h-full max-h-64  overflow-y-scroll  md:overflow-auto">本科就读于区块链工程专业,主要研究方向有区块链、大数据、云原生、全栈开发,技术栈为Golang、JavaScript、Solidity、java。曾在
                            Bytedance Cloudwego Team 的项目中做过contributor、同时荣获2023年Ethereum Contributor。
                            在学习过程中参与并主导过多个项目的研发工程目前主导的开源项目为<a
                                href="https://github.com/0xdoomxy/blog" className=" underline">blog</a>,致力于成为web3技术知识分享的主流平台之一。
                        </p>
                    </div>
                    <div className="w-full pl-4 h-2/3 flex justify-start items-center">
                        <Timeline
                            pending="keep going..."
                            reverse={false}
                            items={[
                                {
                                    children: '2020年进入成都信息工程大学区块链工程专业学习',
                                },
                                {
                                    children: '2022年开启实习之旅',
                                },
                                {
                                    children: '2023年开始在github开源贡献',
                                },
                                {
                                    children: '2024年于成都信息工程大学毕业',
                                }
                            ]}
                        />
                    </div>
                </div>
            </div>
            <div  style={{height: "15%"}} className=" flex w-full h-20 border-t-2  bg-slate-50 justify-center items-center ">
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

export default AboutPage;