

import { About,Main,Article } from "../page";


const routes = [
{
    path:"about",
    element: <About/>,
},
{
    path:"home",
    element:<Main/>,
},
{
    path:"/article/:articleId",
    element:<Article />,
}
];


export default routes;