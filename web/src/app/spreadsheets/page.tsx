'use client'

import Link from "next/link";
import Image from 'next/image'
import Sheet from '../../../public/sheets.svg'
import { HiOutlineMagnifyingGlass } from 'react-icons/hi2'
import { CgProfile } from 'react-icons/cg'
import { PiPlusLight } from 'react-icons/pi'
import { useRouter } from "next/navigation";
import { Authenticate } from "../api/auth";
import { CreateSpreadSheet, GetSpreadSheet } from "../api/spreadsheet";
import { useCallback, useEffect, useRef, useState } from "react";
import SpreadSheetTable from "./components/SpreadSheetTable";
import { SpreadSheet } from "../types/SpreadSheet";
import Template from "./components/Template";
import Loading from "./components/Loading";

export default function Dashboard() {
  const router = useRouter()
  const [spreadsheets, setSpreadSheets] = useState<SpreadSheet[]>([]);
  const [loader, setLoader] = useState<boolean>(true);
  const [profileVisible, setProfileVisible] = useState<boolean>(false);
  const [profileName, setProfileName] = useState<string>("");
  const [createNewSpreadSheetLoader, setCreateNewSpreadSheetLoader] = useState<boolean>(false);
  const [searchItems, setSearchItems] = useState<SpreadSheet[]>([]);

  const authenticate = useCallback(Authenticate, []);
  const getspreadsheet = useCallback(GetSpreadSheet, []);

  const ref1 = useRef<HTMLDivElement>(null);

  const click = useCallback((e: MouseEvent) => {
    if (ref1.current && !ref1.current.contains(e.target as Node)) {
      setProfileVisible(false);
    }
  }, []);

  useEffect(() => {
    document.addEventListener("click", click);

    return () => {
      document.removeEventListener("click", click);
    };
  }, [click]);

  useEffect(() => {
    const access_token = ((new URL(window.location.href).searchParams.get("access_token")) || localStorage.getItem("spreadsheet_access_token"))
    if (access_token === null) {
      return router.push("/")
    } else {
      authenticate(access_token)
        .then(async res => {
          if (res.status === 200) {
            localStorage.setItem("spreadsheet_access_token", access_token)
            const userInfo = await res.json();

            setProfileName(userInfo.login);
            getspreadsheet(access_token, "")
              .then(res => {
                if (res.status === 200) {
                  return res.json();
                } else {
                }
              }).then(res => {
                setSpreadSheets(res);
                setLoader(false);
              })
          } else {
            localStorage.removeItem("spreadsheet_access_token");
            return router.push("/");
          }
        })
    }
  }, [router, authenticate, getspreadsheet]);

  const createSpreadSheet = () => {
    const access_token = ((new URL(window.location.href).searchParams.get("access_token")) || localStorage.getItem("spreadsheet_access_token")) || "";
    setCreateNewSpreadSheetLoader(true);
    CreateSpreadSheet(access_token)
      .then(res => {
        if (res.status === 200) {
          return res.json();
        } else {
        }
      }).then((res: SpreadSheet) => {
        router.push(`/spreadsheets/${res.SK.slice(12)}`)
      })
  }

  return (
    <>
      {createNewSpreadSheetLoader && <>
        <div className="absolute w-[100%] h-[100%] bg-black opacity-20 z-100 flex justify-center align-middle" >
        </div >
        <div className="absolute top-[50%] left-[50%] z-1000"><Loading /></div>
      </>}
      <div className="m-0 p-0">
        <div className="relative h-[64px] w-full bg-[#ffffff] flex justify-between">
          <div className="ml-4 flex mr-4">
            <Link href="/spreadsheets" className='mb-auto mt-auto min-w-[60px] pl-[10px] pr-[10px] flex align-middle justify-center hover:cursor-pointer'>
              <Image title='Sheets Home' width={30} height={30} src={Sheet} alt="sheet-icon" />
            </Link>
            <span className="mt-auto mb-auto inline-block font-sans font-semibold text-2xl text-[#5f6368]">Sheets</span>
          </div>
          <div className="mr-4 h-[48px] w-[60%] bg-[#f1f3f4] mt-auto mb-auto flex align-middle justify-start rounded-xl relative focus-within:bg-[#ffffff] focus-within:scale-[1.01] focus-within:shadow-sm focus-within:shadow-black">
            <div className="left-[4px] top-[4px] absolute h-[40px] w-[40px] flex align-middle justify-center hover:bg-slate-200 hover:rounded-full hover:cursor-pointer">
              <HiOutlineMagnifyingGlass className="w-[25px] h-[25px] mt-auto mb-auto" />
            </div>
            <input type="text" onChange={(e) => {
              setSearchItems(spreadsheets.filter(x => e.target.value !== "" && x.SpreadSheetTitle.toLowerCase().startsWith(e.target.value)))
            }} className="ml-[50px] w-full mt-[8px] mb-[8px] bg-inherit mr-[40px] outline-none" placeholder="Search" />
            {searchItems.length !== 0 && <div className="absolute w-full top-[55px] flex flex-col gap-1">
              {
                searchItems.map(x => (
                  <div key={x.PK} className="w-full h-[30px] rounded-lg bg-slate-200 pl-2 hover:cursor-pointer hover:bg-slate-300 flex flex-col" onClick={() => {
                    router.push(`/spreadsheets/${x.SK.slice(12)}`);
                  }}>
                    <span className="mt-auto mb-auto">{x.SpreadSheetTitle}</span>
                  </div>
                ))
              }
            </div>}
          </div>
          <div ref={ref1} onClick={() => setProfileVisible(!profileVisible)} className="mr-4 mt-auto mb-auto min-h-[44px] min-w-[44px] flex align-middle justify-center hover:bg-slate-200 hover:rounded-full hover:cursor-pointer">
            <CgProfile className="w-[25px] h-[25px] mt-auto mb-auto" />
          </div>
          {profileVisible && <div className="shadow-black shadow-md absolute right-[16px] bottom-[-200px] sm:bottom-[-195px] bg-[#E9EEF6] w-[200px] h-[200px] sm:w-[300px] sm:h-[200px] rounded-2xl flex flex-col align-middle justify-center gap-4">
            <div className="ml-auto mr-auto w-[80%] h-[40px] rounded-2xl text-center">
              <span className="m-auto block font-bold">Hi, {profileName}!!</span>
            </div>
            <div className="ml-auto mr-auto bg-slate-400 w-[80%] h-[40px] rounded-2xl text-center hover:bg-slate-500 hover:cursor-pointer flex" onClick={() => {
              localStorage.removeItem("spreadsheet_access_token");
              router.push("/");
            }}>
              <span className="m-auto block font-bold">Log Out</span>
            </div>
          </div>}
        </div>

        <div className="h-[calc(100vh-64px)] w-full">
          <div className="h-[250px] w-full bg-[#f1f3f4] flex flex-col justify-center align-middle">
            <div className="h-[64px] w-[75%] flex justify-start ml-[14%]">
              <span className="mt-auto mb-auto font-medium font-roboto">Start a new spreadsheet from template</span>
            </div>
            <div className="w-[75%] flex justify-start align-middle m-auto mt-0 mb-auto overflow-x-scroll">
              <Template onClick={createSpreadSheet} templateName="Blank spreadsheet" bgColor="#FFFFFF"/>
              <Template onClick={createSpreadSheet} templateName="Red spreadsheet" bgColor="#FF0000"/>
              <Template onClick={createSpreadSheet} templateName="Green spreadsheet" bgColor="#00FF00"/>
              <Template onClick={createSpreadSheet} templateName="Blue spreadsheet" bgColor="#0000FF"/>
            </div>
          </div>
          <div className="w-[75%] ml-auto mr-auto">
            <SpreadSheetTable loader={loader} spreadsheets={spreadsheets} setSpreadSheets={setSpreadSheets} />
          </div>
        </div>
        <PiPlusLight onClick={createSpreadSheet} className="z-10 fixed bottom-[24px] right-[24px] w-[60px] h-[60px] hover:opacity-[50%] hover:cursor-pointer shadow-sm shadow-black rounded-full bg-white" />
      </div>
    </>
  )
}
