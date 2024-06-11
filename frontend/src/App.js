
import './App.css';
import { Route, HashRouter as Router, Routes } from 'react-router-dom';
import {About, Home,Archives} from "./pages";

function App() {
  return (
  <>
    <Router>
      <Routes>
      <Route path={'/'} Component={Home} />
      <Route path={'/about'} Component={About} />
      <Route path={'/archieve'} Component={Archives} />
      </Routes>
    </Router>
    </>
  );
}

export default App;
