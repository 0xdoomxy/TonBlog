import React,{useEffect,useState,useRef} from "react";
import { Search} from "../components";
import { useParams,useNavigate } from "react-router-dom";
import MarkdownContext from "../components/markdown";

const mdDemo = `



## Markdown Basic Syntax

I just love **bold text**. Italicized text is the _cat's meow_. At the command prompt, type \`nano\`.

My favorite markdown editor is [ByteMD](https://github.com/bytedance/bytemd).

1. First item
2. Second item
3. Third item

> Dorothy followed her through many of the beautiful rooms in her castle.

\`\`\`js
import { Editor, Viewer } from 'bytemd';
import gfm from '@bytemd/plugin-gfm';

const plugins = [
  gfm(),
  // Add more plugins here
];

const editor = new Editor({
  target: document.body, // DOM to render
  props: {
    value: '',
    plugins,
  },
});

editor.on('change', (e) => {
  editor.$set({ value: e.detail.value });
});
\`\`\`
## GFM Extended Syntax

Automatic URL Linking: https://github.com/bytedance/bytemd

~~The world is flat.~~ We now know that the world is round.

- [x] Write the press release
- [ ] Update the website
- [ ] Contact the media

| Syntax    | Description |
| --------- | ----------- |
| Header    | Title       |
| Paragraph | Text        |

## Footnotes

Here's a simple footnote,[^1] and here's a longer one.[^bignote]

[^1]: This is the first footnote.
[^bignote]: Here's one with multiple paragraphs and code.

    Indent paragraphs to include them in the footnote.

    \`{ my code }\`

    Add as many paragraphs as you like.
`;

const Article =()=>{
    //标签颜色
    const labelColorList = ["bg-red-300","bg-yellow-200","bg-green-300","bg-pink-300","bg-gray-200"]
    const [context,setContext] = useState(mdDemo);
    const [article,setArticle] = useState({
        title:"title",
        tags:["1","2","3"],
        date:"2021-10-10",
        author:"0xdoomxy",
        looknum:1233333,
        likenum:123
    })
    const navigate=useNavigate();
    const navItems=[{
        Name:"Home",
        Target:"/"
    },{
        Name:"About",
        Target:"/about"
    },{Name:"Archieve",Target:"/archieve"}]

        //是否需要更换header显示
        const [changeHeader,setChangeHeader]=useState(false);
    //文章唯一id
    const{articleId} =useParams();
    // //markdown文章内容显示ref
    // const [contextDom,setContextDom] =useRef([]);
         //小屏幕点击事件，用来显示菜单栏
     const [showSmallNav,setShowSmallNav]=useState(false);


      //组件初始化的时候执行的函数
    useEffect(()=>{
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
        <div className="w-full h-full pt-12 flex items-start ">
            <div className=" w-1/6"></div>
            <div className=" w-2/3 h-full">
                {/* 简介 */}
                <div className=" flex justify-between w-full h-40">
                    <div className="w-3/4 flex items-start flex-col">
                        <div className=" text-6xl font-normal text-ellipsis">{article.title}</div>
                        <div className=" flex justify-start items-center py-4 ">{article.tags.map((item,index)=>{
                            return (<div key={"tag"+index}  className={`mx-2 md:w-20   border flex justify-center items-center ${labelColorList[index%labelColorList.length]}`} >item</div>)
                        })}</div>
                        <div className=" text-xl font-serif py-1">{article.author}</div>
                        <div className=" text-base font-sans ">{article.date}</div>
                    </div>
                    <div className="w-1/4 h-full flex justify-center flex-col">
                        <div className="h-1/2 border-x-2 border-t-2  text-lg  w-full  font-serif flex items-center justify-center">
                            浏览量:{article.looknum}
                        </div>
                        <div className=" h-1/2 border-2 w-full text-lg  font-serif flex items-center justify-center">
                            点赞量:{article.likenum}
                        </div>
                    </div>
                </div>
                <div className=" w-full h-full pt-20">
              <MarkdownContext context={context}/>
        </div>
        </div>
        <div className=" w-1/6"></div>
        </div>
        </div>
    )
}
export default Article;