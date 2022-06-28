package main

import (
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var InvalidCredentialsError = errors.New("invalid credentials for this function")

var loginAccountResolver = func(params graphql.ResolveParams) (interface{}, error) {
	resp, err := LoginToAccount(LoginInput{
		Email:       params.Args["email"].(string),
		Password:    params.Args["password"].(string),
		AccountType: AccountDescriptiveTitleType(params.Args["accountType"].(string)),
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.GetPlaintext())

	return LoginResponse { resp.GetPlaintext() }, err
}

var createNewROMEntryResolver = func(params graphql.ResolveParams) (interface{}, error) {
	session := SessionIDType{params.Args["sessionId"].(string)}

	obj, err := session.GetAccountMongoEntry()
	if err != nil {
		return nil, err
	}

	resp, err := CreateAndSaveNewRecordedEntryObject(CreateAndSaveNewRecordedEntryObjectInput{
		Patient:       obj.GetMongoID(),
		RangeOfMotion: params.Args["rangeOfMotion"].(float64),
		DateTime:      params.Args["datetimeRecorded"].(time.Time),
		JointID:       JointDescriptiveTitleType(params.Args["jointId"].(string)),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

var createNewPractitionerAccountResolver = func(params graphql.ResolveParams) (interface{}, error) {
	practitioner, err := CreateAndSaveNewPractitionerObject(CreateAndSaveNewPractitionerObjectInput{
		FullName:    params.Args["fullName"].(string),
		Email:       params.Args["emailAddress"].(string),
		Password:    params.Args["password"].(string),
		Institution: params.Args["institution"].(string),
	})
	if err != nil {
		return nil, err
	}

	return practitioner, nil
}

var createNewPatientAccountResolver = func(params graphql.ResolveParams) (interface{}, error) {
	session := SessionIDType{params.Args["sessionId"].(string)}

	if act, err := session.GetAccountType(); err != nil || act != AccountDescriptiveTitleTypePractitioner {
		return nil, InvalidCredentialsError
	}

	obj, err := session.GetAccountMongoEntry()
	if err != nil {
		return nil, err
	}

	patient, err := CreateAndSaveNewPatientObject(CreateAndSaveNewPatientObjectInput{
		FullName:     params.Args["fullName"].(string),
		Email:        params.Args["emailAddress"].(string),
		Password:     params.Args["password"].(string),
		Practitioner: obj.GetMongoID(),
	})
	return patient, nil
}

var getROMEntriesResolver = func(params graphql.ResolveParams) (interface{}, error) {
	var returnVal []MongoObjectRecordedEntry
	var find bson.M

	session := SessionIDType{params.Args["sessionId"].(string)}
	act, err := session.GetAccountType()
	if err != nil {
		return nil, err
	}

	obj, err := session.GetAccountMongoEntry()
	if err != nil {
		return nil, err
	}

	if act == AccountDescriptiveTitleTypePatient {
		find = bson.M{
			"patient_id": obj.GetMongoID(),
		}

	} else {
		patientID, err := primitive.ObjectIDFromHex(params.Args["patientId"].(string))
		if err != nil || !IsPatientUnderPractitioner(IsPatientUnderPractitionerInput{
			Practitioner: obj.GetMongoID(),
			Patient:      patientID,
		}) {
			return nil, InvalidCredentialsError
		}

		find = bson.M{
			"patient_id": patientID,
		}
	}


	dateTime := bson.M{}
	doAppend := false

	if params.Args["dateRangeBegin"] != nil {
		dateTime["$gte"] = params.Args["dateRangeBegin"].(time.Time)
		doAppend = true
	}
	if params.Args["dateRangeEnd"] != nil {
		dateTime["$lt"] = params.Args["dateRangeEnd"].(time.Time)
		doAppend = true
	}
	if doAppend {
		find["datetime_recorded"] = dateTime
	}

	entryCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectRecordedEntryCollectionTitle)
	cursor, err := entryCollection.Find(appVariables.MongoContext, find)
	if err != nil {
		return nil, err
	}

	err = cursor.All(appVariables.MongoContext, &returnVal)
	if err != nil {
		return nil, err
	}

	return returnVal, nil
}

var getPatientAccountDetailsResolver = func(params graphql.ResolveParams) (interface{}, error) {
	var patient MongoObjectPatient

	session := SessionIDType{params.Args["sessionId"].(string)}
	act, err := session.GetAccountType()
	if err != nil {
		return nil, err
	}
	if act != AccountDescriptiveTitleTypePatient {
		return nil, InvalidCredentialsError
	}

	obj, err := session.GetAccountMongoEntry()
	if err != nil {
		return nil, err
	}

	patientCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectPatientCollectionTitle)

	err = patientCollection.FindOne(
		appVariables.MongoContext,
		&bson.M{"_id": obj.GetMongoID()},
	).Decode(&patient)
	if err != nil {
		return nil, err
	}

	return patient, nil
}

var getPractitionerAccountDetailsResolver = func(params graphql.ResolveParams) (interface{}, error) {
	var practitioner MongoObjectPractitioner

	session := SessionIDType{params.Args["sessionId"].(string)}
	println(params.Args["sessionId"].(string))
	act, err := session.GetAccountType()
	if err != nil {
		return nil, err
	}
	if act != AccountDescriptiveTitleTypePractitioner {
		return nil, InvalidCredentialsError
	}

	obj, err := session.GetAccountMongoEntry()
	if err != nil {
		return nil, err
	}

	practitionerCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectPractitionerCollectionTitle)

	err = practitionerCollection.FindOne(
		appVariables.MongoContext,
		&bson.M{"_id": obj.GetMongoID()},
	).Decode(&practitioner)
	if err != nil {
		return nil, err
	}

	return practitioner, nil
}

var getPractitionerAccountsPatients = func(params graphql.ResolveParams) (interface{}, error) {
	// For now, we are hard-coding this. This is obviously not how this will work in production.
	return []string{"625d6b075b86fa89f3749cb9"}, nil
}

var getAllStudyEntries = func(params graphql.ResolveParams) (interface{}, error) {
	var resp []MongoObjectPaperRecordedEntry

	paperCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection("paperCollections")
	cursor, err := paperCollection.Find(appVariables.MongoContext, &bson.M{})
	if err != nil {
		return nil, err
	}

	cursor.All(appVariables.MongoContext, &resp)
	return resp, nil
}

var createNewStudyEntry = func(params graphql.ResolveParams) (interface{}, error) {
	var leftKnee = MongoObjectKnee{
		GoniometerTruth: params.Args["leftKneeGoniometerTruth"].(float64),
		AppFrontMeasure: params.Args["leftKneeAppFrontMeasure"].(float64),
		AppSideMeasure:  params.Args["leftKneeAppSideMeasure"].(float64),
	}

	var rightKnee = MongoObjectKnee{
		GoniometerTruth: params.Args["rightKneeGoniometerTruth"].(float64),
		AppFrontMeasure: params.Args["rightKneeAppFrontMeasure"].(float64),
		AppSideMeasure:  params.Args["rightKneeAppSideMeasure"].(float64),
	}

	var newStudyEntry = MongoObjectPaperRecordedEntry{
		MongoID:       primitive.NewObjectID(),
		PatientHeight: params.Args["patientHeight"].(float64),
		PatientWeight: params.Args["patientWeight"].(float64),
		LeftKnee:      leftKnee,
		RightKnee:     rightKnee,
	}

	paperCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection("paperCollections")

	if _, err := paperCollection.InsertOne(appVariables.MongoContext, newStudyEntry); err != nil {
		return nil, err
	}
	return newStudyEntry, nil
}