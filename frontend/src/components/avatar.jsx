import { useState } from "react";
import {toUserFriendlyAddress  } from '@tonconnect/ui-react';
import { useNavigate } from "react-router-dom";

const TonAvatar = ({wallet,disconnect}) => {
    const [dropdown,setDropDown] =useState(false);
    const navigate = useNavigate();
    return (
        <> 
       {(wallet!==undefined && wallet !== null)?    
     <>
<button onClick={()=>{setDropDown(origin=>setDropDown(!origin))}} type="button" className="relative inline-flex items-center justify-center w-10 h-10 overflow-hidden bg-gray-100 rounded-full dark:bg-gray-600">
    <span className="font-medium text-gray-600 dark:text-gray-300">Ton</span>
</button>
{dropdown?<div id="dropdownInformation" className=" bg-white divide-y divide-gray-100 rounded-lg shadow w-44 dark:bg-gray-700 dark:divide-gray-600 absolute top-12 ">
    <div className="px-4 py-3 text-sm text-gray-900 dark:text-white">
      <div>Welcome to</div>
      <div className="font-medium truncate">{wallet!==undefined && wallet !== null&&toUserFriendlyAddress(wallet.account.address)}</div>
    </div>
    <ul className="py-2 text-sm text-gray-700 dark:text-gray-200" aria-labelledby="dropdownInformationButton">
    <li>
        <a onClick={()=>{navigate("/article/create")}} className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">Create</a>
      </li>
      <li>
        <a href="#" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">Dashboard</a>
      </li>
      <li>
        <a href="#" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">Collection</a>
      </li>
    </ul>
    <div className="py-2">
      <a href="#" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 dark:text-gray-200 dark:hover:text-white" onClick={()=>{disconnect()}}>Sign out</a>
    </div>
</div>:<> </>}
</>
:<></>}
 </>
    )
}

export default TonAvatar;

