import   '../css/markdown.css';
import gfm from "@bytemd/plugin-gfm";
import gemoji from "@bytemd/plugin-gemoji";
import highlight from "@bytemd/plugin-highlight-ssr";
import mediumZoom from "@bytemd/plugin-medium-zoom";
import { Editor, Viewer } from "@bytemd/react";
import "github-markdown-css";
    //markdown文章显示插件
    const plugins = [gfm(), gemoji(), highlight(), mediumZoom()];
const MarkdownContext = ({context})=>{
    return (
        <Viewer   value={context} plugins={plugins} />
    )
    }

    export default MarkdownContext;