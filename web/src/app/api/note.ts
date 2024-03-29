export async function CreateNote(access_token: string, content: string, spreadsheetID: string, sheetNo: number, cellID: string) {
    const response = await fetch(`${process.env.API_DOMAIN}/api/note`, {
        method: "POST",
        headers: {
            'spreadsheet_access_token': access_token,
        },
        body: JSON.stringify({
            "Content": content,
            "SpreadSheetID": spreadsheetID,
            "SheetNo": sheetNo,
            "CellID": cellID,
        })
    })

    return response;
}

export async function GetNotes(access_token: string, spreadsheet_id: string) {
    const response = await fetch(`${process.env.API_DOMAIN}/api/note?spreadsheet_id=${spreadsheet_id}`, {
        method: "GET",
        headers: {
            'spreadsheet_access_token': access_token,
        },
    })

    return response;
}