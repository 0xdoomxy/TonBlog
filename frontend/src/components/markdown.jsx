import React, {useMemo} from 'react';
import '../css/markdown.css';
import gfm from '@bytemd/plugin-gfm'
import highlightssr from '@bytemd/plugin-highlight-ssr'
import highlight from '@bytemd/plugin-highlight'
import breaks from '@bytemd/plugin-breaks'
import footnotes from '@bytemd/plugin-footnotes'
import frontmatter from '@bytemd/plugin-frontmatter'
import gemoji from '@bytemd/plugin-gemoji'
import mediumZoom from '@bytemd/plugin-medium-zoom'
import "highlight.js/styles/vs.css";
import { Viewer} from "@bytemd/react";
import "github-markdown-css";

const MarkdownContext = ({context}) => {
    //markdown文章显示插件
    const plugins = useMemo(() => [gfm(), highlightssr(), highlight(), breaks(), footnotes(), frontmatter(), gemoji(), mediumZoom()], []);
    return (
        <Viewer value={context} plugins={plugins}/>
    )
}

export default MarkdownContext;