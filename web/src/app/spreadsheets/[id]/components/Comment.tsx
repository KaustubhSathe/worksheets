import { RootState } from "@/app/lib/redux/store";
import { useState } from "react";
import { useSelector } from "react-redux";

export default function Comment() {
    const spreadSheetMetaData = useSelector((state: RootState) => state.spreadSheetMetaData).value;
    const [comment, setComment] = useState();
    return (
        <div id="comment" className="w-[350px] hidden bg-white shadow-lg shadow-slate-400 rounded-lg p-4">
            <div className="flex gap-3">
                <div className="w-[40px] h-[40px] rounded-full bg-teal-800 flex justify-center">
                    <span className="mt-auto mb-auto text-white font-normal text-2xl">{spreadSheetMetaData?.UserName?.at(0)}</span>
                </div>
                <span className="mt-auto mb-auto text-black font-normal text-2xl">{spreadSheetMetaData?.UserName}</span>
            </div>
            <input type="text" value={comment} className="resize-none outline-none border-[1px] rounded-full pl-4 pr-4 pt-2 pb-2 mb-3 mt-3 border-black w-full" />
            <div className="flex justify-end gap-3">
                <button className="w-[75px] h-[36px] text-blue-700 font-semibold hover:rounded-full hover:bg-blue-100">
                    Cancel
                </button>
                <button className="mr-4 w-[100px] h-[36px] font-semibold text-gray-500 bg-gray-200 rounded-full">
                    Comment
                </button>
            </div>
        </div>
    )
}