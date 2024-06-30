


/**
 * 从一段 markdown 文本中匹配出图片的 url
 * @param {*} markdown 
 * @returns 
 */
function MatchImageUrlFromMarkdown(markdown) {
    var imageRegex = /!\[.*\]\((.*)\)/g;
    var matches = [];
    var match;

    while ((match = imageRegex.exec(markdown)) !== null) {
        matches.push(match[1]);
    }

    return matches;
}

export { MatchImageUrlFromMarkdown };