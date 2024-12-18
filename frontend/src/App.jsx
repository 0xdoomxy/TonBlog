import './App.css';
import {HashRouter as Router, Route, Routes} from 'react-router-dom';
import {AboutPage, ArticlePage, CreatePage, HomePage, HotDetails, NewDetails, SearchPage, TagDetails} from "./pages";
import {Bounce, ToastContainer} from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import Airport from "./pages/airport";
import React, {createContext, useState} from "react";
import {DiscoverWalletProviders} from "./components/WalletProviders";
export const Web3Wallet = createContext();
export default function App() {
  const [selectedWallet, setSelectedWallet] = useState();
  const [userAccount, setUserAccount] = useState('');
  const [searchWalletModal,setSearchWalletModal] = useState(false);

    return (
        <>
          <Web3Wallet.Provider value={{searchWalletModal,setSearchWalletModal,selectedWallet,setSelectedWallet,userAccount,setUserAccount}}>
              <DiscoverWalletProviders  searchWalletModal={searchWalletModal} setSearchWalletModal={setSearchWalletModal} setSelectedWallet={setSelectedWallet} setUserAccount={setUserAccount}/>
              <Router>
                <Routes>
                    <Route path={'/'} Component={HomePage}/>
                    <Route path={'/about'} Component={AboutPage}/>
                    <Route path={'/article/create'} Component={CreatePage}/>
                    <Route path={'/article/:articleId'} Component={ArticlePage}/>
                    <Route path={'/search'} Component={SearchPage}/>
                    <Route path={"/article/hot"} Component={HotDetails}/>
                    <Route path={'/article/newest'} Component={NewDetails}/>
                    <Route path={"/articles/tag"} Component={TagDetails}/>
                    <Route path={"/airport"} Component={Airport}/>
                </Routes>
            </Router>
          </Web3Wallet.Provider>
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
