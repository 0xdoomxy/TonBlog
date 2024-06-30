import {React, useMemo, useState}  from 'react';
import 'bytemd/dist/index.css';
import "highlight.js/styles/vs.css";
import { Editor } from '@bytemd/react';
import gfm from '@bytemd/plugin-gfm'
import highlightssr from '@bytemd/plugin-highlight-ssr'
import highlight from '@bytemd/plugin-highlight'
import breaks from '@bytemd/plugin-breaks'
import footnotes from '@bytemd/plugin-footnotes'
import frontmatter from '@bytemd/plugin-frontmatter'
import gemoji from '@bytemd/plugin-gemoji'
import mediumZoom from '@bytemd/plugin-medium-zoom'
import zhHans from 'bytemd/lib/locales/zh_Hans.json'
import agent from '../agent/agent';
import { MatchImageUrlFromMarkdown } from '../util/image';
const Create = () => {
    const plugins = useMemo(() => [gfm(), highlightssr(), highlight(), breaks(), footnotes(), frontmatter(), gemoji(), mediumZoom()], []);
    const [article,setArticle] = useState({
        title:'请输入文章标题',
        content:'',
        tags:[],
    })
    //用户上传的所有图片(包括为用到的和已用到的图片)
    const [allPictures,setAllPictures] = useState([]);
    function publicArticle(){
        console.log(article)
        console.log(allPictures)
        var realP = [];
        const allUrl =  MatchImageUrlFromMarkdown(article.content);
        for (let i =0;i<allUrl.length;i++){
            for(let j =0;j<allPictures.length;j++){
                if (allUrl[i].includes(allPictures[j])){
                    realP.push(allPictures[j]);
            }
        }
        let pushParams = new FormData();
        pushParams.append('title',article.title);
        pushParams.append('content',article.content);
        pushParams.append('tags',article.tags.join(','));
        pushParams.append('images',realP.join(','));
        agent.Article.Publish(pushParams).then((res)=>{
            console.log(res);
            if (res.status){
                alert('发布成功');
            }
        });
    }
    }
    return (
        <div className='w-full h-full flex flex-col'>
            <div className='w-full h-20 flex flex-row justify-start py-2 '>
                <div className='md:w-24 flex justify-center items-center '>
                <svg  xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className=" w-12 cursor-pointer ">
  <path strokeLinecap="round" strokeLinejoin="round" d="M9 15 3 9m0 0 6-6M3 9h12a6 6 0 0 1 0 12h-3" />
</svg>

                </div>
                <input  className='  w-4/5 indent-8 text-3xl font-semibold font-serif outline-none ' onChange={(e)=>{setArticle({...article,title:e.target.value})}} placeholder={article.title} ></input>
                <div className=' flex justify-end items-center  h-full w-1/10'>
                    <button  className=' h-full  md:w-28 border-2 rounded-xl bg-blue-400 ' onClick={publicArticle}>发布文章</button>
                </div>
            </div>
        <Editor className=" w-full"
       locale={zhHans}
       value={article.content}  //markdown内容
       plugins={plugins}  //markdown中用到的插件，如表格、数学公式、流程图
       onChange={(v) => {
         setArticle({ ...article, content: v });
       }}
       uploadImages={async (files) => {   
            //上传图片
            let form = new FormData();
            form.append('file', files[0]);
            let filename = await agent.Article.ImageUpload(form).then((res) => {
                if (res.status) {
                    return res.data;
                }
            });
            setAllPictures((origin)=>[...origin, filename]);
         return [
                {
                   title: filename,
                   url: agent.Article.ImageDownloadUrl(filename),
                },
                ];
        }}
 />
        </div> 
    );
    }

export default Create;