import superagentPromise from 'superagent-promise';
import _superagent from 'superagent';


const superagent = superagentPromise(_superagent, Promise);


const API_ROOT ="http://localhost:8080";

const responseBody = (res) =>{
    return res.body;
  }

const encode =encodeURIComponent;

const requests = {
    del: url =>
      superagent.del(`${API_ROOT}${url}`).then(responseBody),
    get: (url) =>
      superagent.get(`${API_ROOT}${url}`).then(responseBody),
    put: (url, body) =>
      superagent.put(`${API_ROOT}${url}`, body).then(responseBody),
    post: (url, body) =>
      superagent.post(`${API_ROOT}${url}`, body).then(responseBody),
  };

const Tag = {
  GetAllTags:()=>requests.get(`/tag/findall`),
}

const Article = {
    ImageDownload:(file)=>requests.get(`/article/image/download?filename=${encode(file)}`),
    ImageUpload:(file)=>requests.post(`/article/image/upload`,file),
    ImageDownloadUrl:(filename)=>`${API_ROOT}/article/image/download?filename=${encode(filename)}`,
    Publish:(article)=>requests.post(`/article/publish`,article),
    Find:(articleId)=>requests.get(`/article/find?id=${encode(articleId)}`),
    FindMaxAccess:(page,pagesize)=>requests.get(`/article/findbymaxaccess?page=${encode(page)}&pagesize=${encode(pagesize)}`),
    FindNewest:(page,pagesize)=>requests.get(`/article/findbycreatetime?page=${encode(page)}&pagesize=${encode(pagesize)}`),
    Search:(keyword,page,pagesize)=>requests.get(`/article/search?page=${encode(page)}&pagesize=${encode(pagesize)}&keyword=${encode(keyword)}`)
}


export default{
    Article,
    Tag,
    API_ROOT
}