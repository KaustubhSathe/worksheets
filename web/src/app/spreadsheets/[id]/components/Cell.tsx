import { setValue as setValueFormulaBar } from '../../../lib/redux/formulaBarSlice'
import { setValue as setSelectStart } from '../../../lib/redux/selectStartSlice';
import { STATUS, setValue as setSaved } from "../../../lib/redux/savedSlice"
import { useDispatch } from "react-redux";
import globals from '@/app/lib/globals/globals';
import { useRef } from 'react';
import { DetailedCellError, ExportedCellChange, ExportedChange, FunctionPluginDefinition, SimpleCellAddress } from 'hyperformula';


export default function Cell({ i, j }: { i: number, j: number }) {
    const dispatch = useDispatch()
    const id = String.fromCharCode(65 + j) + (i + 1).toString();
    const oldText = useRef<string>("");

    return (
        <div data-testid="cell-wrapper" className={`relative m-0 p-0 w-full rowbar-${(i + 1).toString()} h-[30px] hover:cursor-cell focus:cursor-text`}>
            <textarea data-testid="cell-textbox" className="overflow-hidden text-sm peer hover:cursor-cell focus:cursor-text overflow-x-clip overflow-y-clip p-[4px] break-words break-all h-full w-full border-b-[1px] border-r-[1px] border-solid border-[#E1E1E1] outline-none m-0 resize-none"
                spellCheck={false}
                id={id}
                onFocus={(e) => {
                    oldText.current = e.currentTarget.value
                    globals.selectStart = e.currentTarget.id
                    dispatch(setSelectStart(e.currentTarget.id))
                    const fontSelector = document.getElementById("fontSelector") as HTMLSelectElement
                    const fontSizeSelector = document.getElementById("fontSizeSelector") as HTMLInputElement
                    const boldSelector = document.getElementById("boldSelector") as HTMLButtonElement
                    const itaclicSelector = document.getElementById("italicSelector") as HTMLButtonElement
                    const strikethroughSelector = document.getElementById("strikethroughSelector") as HTMLButtonElement
                    if (globals.spreadsheet.Versions[0].Sheets[globals.selectedSheet].State[id]) {
                        fontSelector.value = globals.spreadsheet.Versions[0].Sheets[globals.selectedSheet].State[id].FontFamily
                        boldSelector.style.backgroundColor = globals.spreadsheet.Versions[0].Sheets[globals.selectedSheet].State[id].FontWeight === "bold" ? "#d3e3fd" : "inherit"
                        fontSizeSelector.value = globals.spreadsheet.Versions[0].Sheets[globals.selectedSheet].State[id].FontSize ? globals.spreadsheet.Versions[0].Sheets[globals.selectedSheet].State[id].FontSize.toString() : "16"
                        itaclicSelector.style.backgroundColor = globals.spreadsheet.Versions[0].Sheets[globals.selectedSheet].State[id].FontStyle === "italic" ? "#d3e3fd" : "inherit"
                        strikethroughSelector.style.backgroundColor = globals.spreadsheet.Versions[0].Sheets[globals.selectedSheet].State[id].TextDecoration === "line-through" ? "#d3e3fd" : "inherit"
                    }
                }}
                onInput={(e) => {
                    if (globals.saved) {
                        globals.saved = false
                        dispatch(setSaved(STATUS.UNSAVED))
                    }
                    const currentText = e.currentTarget.value
                    const oldTextVal = oldText.current
                    const currentTarget = e.currentTarget
                    globals.undoStack.push({
                        Action: () => {
                            currentTarget.value = currentText
                        },
                        Inverse: () => {
                            currentTarget.value = oldTextVal
                        }
                    })
                    const address = globals.hfInstance.simpleCellAddressFromString(id, globals.selectedSheet)
                    const changes: ExportedChange[] = globals.hfInstance.setCellContents(address as SimpleCellAddress, currentText)
                    for (let i = 1; i < changes.length; i++) {
                        const key = String.fromCharCode(65 + (changes[i] as ExportedCellChange).col) + ((changes[i] as ExportedCellChange).row + 1).toString()
                        const elem = document.getElementById(key) as HTMLTextAreaElement
                        if (elem) {
                            elem.value = changes[i].newValue?.toString() as string
                        }
                    }
                    const calculated = globals.hfInstance.getCellValue(address as SimpleCellAddress)
                    currentTarget.value = calculated instanceof DetailedCellError ? currentText : calculated?.toString() as string
                    oldText.current = calculated instanceof DetailedCellError ? currentText : calculated?.toString() as string
                    dispatch(setValueFormulaBar(currentText))
                }}
                key={String.fromCharCode(65 + j) + (i + 1).toString()}
            />
            <div id={id + "comment"}
                className="w-0 h-0 border-r-[10px] border-solid border-b-[10px] border-b-transparent border-r-transparent absolute top-0 right-0 z-50"></div>
            <div className="absolute bottom-[-3px] right-[-3px] w-[10px] h-[10px] rounded-full peer-focus:bg-[#1a73e8] peer-focus:hover:cursor-crosshair peer-focus:z-10">
            </div>
        </div>
    );
}