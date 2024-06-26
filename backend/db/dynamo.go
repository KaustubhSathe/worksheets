package db

import (
	"backend-go/config"
	"backend-go/db/model"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Dynamo struct {
	Client *dynamodb.DynamoDB
}

func NewDynamo() *Dynamo {
	return &Dynamo{
		Client: initializeDynamo(),
	}
}

func initializeDynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	return svc
}

func (db *Dynamo) CreateUser(userInfo *model.UserInfo) error {
	fmt.Println(userInfo)
	user, err := dynamodbattribute.MarshalMap(model.User{
		Base: model.Base{
			PK:        db.UserPK(int64(userInfo.User["id"].(float64))),
			SK:        db.UserSK(int64(userInfo.User["id"].(float64))),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserName: userInfo.User["login"].(string),
		ID:       int64(userInfo.User["id"].(float64)),
	})
	if err != nil {
		log.Fatalf("Got error marshalling User: %s", err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      user,
		TableName: aws.String(config.SPREADSHEETTABLE),
	}

	_, err = db.Client.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return err
	}

	return nil
}

func (db *Dynamo) CreateSpreadSheet(spreadsheetID string, user *model.User) (*model.SpreadSheet, error) {
	ss := &model.SpreadSheet{
		Base: model.Base{
			PK:        db.UserPK(user.ID),
			SK:        db.SpreadSheetSK(spreadsheetID),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserName:         user.UserName,
		UserID:           user.ID,
		SpreadSheetTitle: "Untitled spreadsheet",
		Favorited:        false,
		Versions: []model.Version{
			{
				VersionName: "Version1",
				VersionID:   uuid.NewString(),
				CreatedAt:   time.Now(),
				Sheets:      []model.Sheet{},
			},
		},
		LastOpened: time.Now(),
	}
	spreadsheet, err := dynamodbattribute.MarshalMap(ss)
	if err != nil {
		return nil, err
	}
	ss1 := &model.SpreadSheet{
		Base: model.Base{
			PK:        db.SpreadSheetPK(spreadsheetID),
			SK:        db.SpreadSheetSK(spreadsheetID),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserName:         user.UserName,
		UserID:           user.ID,
		SpreadSheetTitle: "Untitled spreadsheet",
		Favorited:        false,
		Versions: []model.Version{
			{
				VersionName: "Version1",
				VersionID:   uuid.NewString(),
				CreatedAt:   time.Now(),
				Sheets:      []model.Sheet{},
			},
		},
		LastOpened: time.Now(),
	}
	spreadsheet1, err := dynamodbattribute.MarshalMap(ss1)
	if err != nil {
		return nil, err
	}

	_, err = db.Client.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			config.SPREADSHEETTABLE: {
				{
					PutRequest: &dynamodb.PutRequest{
						Item: spreadsheet,
					},
				},
				{
					PutRequest: &dynamodb.PutRequest{
						Item: spreadsheet1,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func (db *Dynamo) CopySpreadSheet(spreadsheetCopy *model.SpreadSheet, user *model.User) (*model.SpreadSheet, error) {
	ss := &model.SpreadSheet{
		Base: model.Base{
			PK:        db.UserPK(user.ID),
			SK:        spreadsheetCopy.SK,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserName:         user.UserName,
		UserID:           user.ID,
		SpreadSheetTitle: spreadsheetCopy.SpreadSheetTitle,
		Favorited:        spreadsheetCopy.Favorited,
		Versions:         spreadsheetCopy.Versions,
		LastOpened:       time.Now(),
	}
	spreadsheet, err := dynamodbattribute.MarshalMap(ss)
	if err != nil {
		return nil, err
	}

	_, err = db.Client.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			config.SPREADSHEETTABLE: {
				{
					PutRequest: &dynamodb.PutRequest{
						Item: spreadsheet,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func (db *Dynamo) DeleteSpreadSheet(spreadsheetID string, user *model.User) (*model.SpreadSheet, error) {
	_, err := db.Client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(config.SPREADSHEETTABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(db.UserPK(user.ID)),
			},
			"SK": {
				S: aws.String(db.SpreadSheetSK(spreadsheetID)),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (db *Dynamo) UpdateSpreadSheetTitle(spreadsheetID string, user *model.User, newTitle string) error {
	_, err := db.Client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(config.SPREADSHEETTABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(db.UserPK(user.ID)),
			},
			"SK": {
				S: aws.String(db.SpreadSheetSK(spreadsheetID)),
			},
		},
		UpdateExpression: aws.String("set SpreadSheetTitle = :spreadSheetTitle"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":spreadSheetTitle": {
				S: aws.String(newTitle),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Dynamo) UpdateSpreadSheets(spreadsheetID string, user *model.User, newTitle string) error {
	_, err := db.Client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(config.SPREADSHEETTABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(db.UserPK(user.ID)),
			},
			"SK": {
				S: aws.String(db.SpreadSheetSK(spreadsheetID)),
			},
		},
		UpdateExpression: aws.String("set SpreadSheetTitle = :spreadSheetTitle"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":spreadSheetTitle": {
				S: aws.String(newTitle),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Dynamo) GetSpreadSheet(spreadsheetID string) (*model.SpreadSheet, error) {
	res, err := db.Client.Query(&dynamodb.QueryInput{
		TableName:              aws.String(config.SPREADSHEETTABLE),
		KeyConditionExpression: aws.String("#PK = :pk AND #SK = :sk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(db.SpreadSheetPK(spreadsheetID)),
			},
			":sk": {
				S: aws.String(db.SpreadSheetSK(spreadsheetID)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#PK": aws.String("PK"),
			"#SK": aws.String("SK"),
		},
	})
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return nil, nil
	}

	spreadsheet := model.SpreadSheet{}
	err = dynamodbattribute.UnmarshalMap(res.Items[0], &spreadsheet)
	if err != nil {
		return nil, err
	}
	return &spreadsheet, nil
}

func (db *Dynamo) GetSpreadSheets(spreadsheetID string, userID int64) ([]*model.SpreadSheet, error) {
	res, err := db.Client.Query(&dynamodb.QueryInput{
		TableName:              aws.String(config.SPREADSHEETTABLE),
		KeyConditionExpression: aws.String("#PK = :pk AND begins_with(#SK, :sk)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(db.UserPK(userID)),
			},
			":sk": {
				S: aws.String(db.SpreadSheetSK(spreadsheetID)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#PK": aws.String("PK"),
			"#SK": aws.String("SK"),
		},
	})
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return nil, nil
	}
	spreadsheets := []*model.SpreadSheet{}

	for i := 0; i < len(res.Items); i++ {
		spreadsheet := model.SpreadSheet{}
		err = dynamodbattribute.UnmarshalMap(res.Items[i], &spreadsheet)
		if err != nil {
			return nil, err
		}
		spreadsheets = append(spreadsheets, &spreadsheet)
	}

	return spreadsheets, nil
}

func (db *Dynamo) CreateComment(spreadsheetID string, userID int64, userName string, sheetNo int64, cellID string, content string) (*model.Comment, error) {
	cc := &model.Comment{
		Base: model.Base{
			PK:        db.SpreadSheetPK(spreadsheetID),
			SK:        db.CommentSK(uuid.NewString()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserName:      userName,
		UserID:        userID,
		SpreadSheetID: spreadsheetID,
		SheetNo:       sheetNo,
		CellID:        cellID,
		Content:       content,
	}
	comment, err := dynamodbattribute.MarshalMap(cc)
	if err != nil {
		return nil, err
	}

	_, err = db.Client.PutItem(&dynamodb.PutItemInput{
		Item:      comment,
		TableName: aws.String(config.SPREADSHEETTABLE),
	})
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func (db *Dynamo) GetComments(spreadsheetID string) ([]*model.Comment, error) {
	res, err := db.Client.Query(&dynamodb.QueryInput{
		TableName:              aws.String(config.SPREADSHEETTABLE),
		KeyConditionExpression: aws.String("#PK = :pk AND begins_with(#SK, :sk)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(db.SpreadSheetPK(spreadsheetID)),
			},
			":sk": {
				S: aws.String(db.CommentSK("")),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#PK": aws.String("PK"),
			"#SK": aws.String("SK"),
		},
	})
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return nil, nil
	}
	comments := []*model.Comment{}

	for i := 0; i < len(res.Items); i++ {
		comment := model.Comment{}
		err = dynamodbattribute.UnmarshalMap(res.Items[i], &comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (db *Dynamo) DeleteComment(spreadsheetID string, commentID string) error {
	_, err := db.Client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(config.SPREADSHEETTABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(db.SpreadSheetPK(spreadsheetID)),
			},
			"SK": {
				S: aws.String(db.CommentSK(commentID)),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Dynamo) UpdateComment(spreadsheetID, commentID, content string) error {
	_, err := db.Client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(config.SPREADSHEETTABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(db.SpreadSheetPK(spreadsheetID)),
			},
			"SK": {
				S: aws.String(db.CommentSK(commentID)),
			},
		},
		UpdateExpression: aws.String("set Content = :content"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":content": {
				S: aws.String(content),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Dynamo) CreateNote(spreadsheetID string, userID int64, userName string, sheetNo int64, cellID string, content string) (*model.Note, error) {
	nn := &model.Note{
		Base: model.Base{
			PK:        db.SpreadSheetPK(spreadsheetID),
			SK:        db.NoteSK(cellID),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserName:      userName,
		UserID:        userID,
		SpreadSheetID: spreadsheetID,
		SheetNo:       sheetNo,
		CellID:        cellID,
		Content:       content,
	}
	note, err := dynamodbattribute.MarshalMap(nn)
	if err != nil {
		return nil, err
	}

	_, err = db.Client.PutItem(&dynamodb.PutItemInput{
		Item:      note,
		TableName: aws.String(config.SPREADSHEETTABLE),
	})
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func (db *Dynamo) GetNotes(spreadsheetID string) ([]*model.Note, error) {
	res, err := db.Client.Query(&dynamodb.QueryInput{
		TableName:              aws.String(config.SPREADSHEETTABLE),
		KeyConditionExpression: aws.String("#PK = :pk AND begins_with(#SK, :sk)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(db.SpreadSheetPK(spreadsheetID)),
			},
			":sk": {
				S: aws.String(db.NoteSK("")),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#PK": aws.String("PK"),
			"#SK": aws.String("SK"),
		},
	})
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return nil, nil
	}
	notes := []*model.Note{}

	for i := 0; i < len(res.Items); i++ {
		note := model.Note{}
		err = dynamodbattribute.UnmarshalMap(res.Items[i], &note)
		if err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

func (db *Dynamo) DeleteNote(spreadsheetID string, noteID string) error {
	_, err := db.Client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(config.SPREADSHEETTABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(db.SpreadSheetPK(spreadsheetID)),
			},
			"SK": {
				S: aws.String(db.NoteSK(noteID)),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Dynamo) UpdateNote(spreadsheetID, noteID, content string) error {
	_, err := db.Client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(config.SPREADSHEETTABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(db.SpreadSheetPK(spreadsheetID)),
			},
			"SK": {
				S: aws.String(db.NoteSK(noteID)),
			},
		},
		UpdateExpression: aws.String("set Content = :content"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":content": {
				S: aws.String(content),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Dynamo) UserPK(userID int64) string {
	return fmt.Sprintf("USER#%d", userID)
}

func (db *Dynamo) UserSK(userID int64) string {
	return fmt.Sprintf("USER#%d", userID)
}

func (db *Dynamo) SpreadSheetPK(spreadsheetID string) string {
	return fmt.Sprintf("SPREADSHEET#%s", spreadsheetID)
}

func (db *Dynamo) SpreadSheetSK(spreadsheetID string) string {
	return fmt.Sprintf("SPREADSHEET#%s", spreadsheetID)
}

func (db *Dynamo) CommentSK(commentID string) string {
	return fmt.Sprintf("COMMENT#%s", commentID)
}

func (db *Dynamo) NoteSK(noteID string) string {
	return fmt.Sprintf("NOTE#%s", noteID)
}
