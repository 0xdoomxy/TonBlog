import React from 'react'


export const Search = ({onKeyDown}) => {
    return (
        <div className="w-full h-full flex justify-end md:justify-center items-center">
            <div className=" h-4/5 w-4/5 border-2 flex justify-between  items-center rounded-2xl mr-2 ">
                <div className=" absolute px-2">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                         stroke="currentColor" className="size-5">
                        <path strokeLinecap="round" strokeLinejoin="round"
                              d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"/>
                    </svg>
                </div>
                <input type="input" onKeyDown={(target) => onKeyDown(target)}
                       className="pl-8 bg-slate-50 w-full rounded-2xl h-full">
                </input>
            </div>
        </div>
    )

}
