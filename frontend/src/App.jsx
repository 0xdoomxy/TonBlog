
import './App.css';
import { Route, HashRouter as Router, Routes } from 'react-router-dom';
import { CreatePage, ArticlePage, SearchPage, HomePage, ArchivesPage, AboutPage,HotDetails, NewDetails, TagDetails} from "./pages";
import { Bounce, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css'; 
import { useBackendAuth } from "./hooks/ton";
function App() {
  return (
  <>
        {useBackendAuth()}
    <Router>
      <Routes>
      <Route path={'/'} Component={HomePage} />
      <Route path={'/about'} Component={AboutPage} />
      <Route path={'/archieve'} Component={ArchivesPage} />
      <Route path={'/article/create'} Component={CreatePage}/> 
      <Route path={'/article/:articleId'} Component={ArticlePage} />
      <Route path={'/search'} Component={SearchPage}/>
      <Route path={"/article/hot"} Component={HotDetails}/>
      <Route path={'/article/newest'} Component={NewDetails}/>
      <Route path={"/articles/tag"} Component={TagDetails} />
      </Routes>
    </Router>
    <ToastContainer
position="top-center"
autoClose={5000}
hideProgressBar={false}
newestOnTop={false}
closeOnClick
rtl={false}
pauseOnFocusLoss
draggable
pauseOnHover
theme="light"
transition={Bounce}
/>
    </>
  );
}

export default App;
