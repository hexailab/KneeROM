package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


/*
 *
 * AccountDescriptiveTitle Struct and Modifiers
 *
 */

type AccountDescriptiveTitleType string
var UnknownAccountDescriptiveTitleError = errors.New("the given string is not a valid descriptive title type")
const (
	AccountDescriptiveTitleTypePatient      AccountDescriptiveTitleType = "ACCOUNT_TYPE_PATIENT"
	AccountDescriptiveTitleTypePractitioner AccountDescriptiveTitleType = "ACCOUNT_TYPE_PRACTITIONER"
	AccountDescriptiveTitleTypeInvalid      AccountDescriptiveTitleType = "ACCOUNT_TYPE_INVALID"
)

// GetUniformStringIdentifier
// Returns a 4-letter string which is unique for each account type.
// Used for various Object identifiers. The suffixing underscore is also returned.
func (acctTitle AccountDescriptiveTitleType) GetUniformStringIdentifier() (string, error) {
	switch acctTitle {
	case AccountDescriptiveTitleTypePatient:
		return "ptt_", nil
	case AccountDescriptiveTitleTypePractitioner:
		return "prc_", nil
	default:
		return "err_", UnknownAccountDescriptiveTitleError
	}
}


/*
 *
 * SessionIDHash Struct and Modifiers
 *
 */

type SessionIDType struct {
	text string
}
const (
	SessionIDTypeLengthWithoutPrefix   int32  = 116 // Must be a multiple of 4.
	SessionIDTypePrefixString          string = "session_"
)

func (s SessionIDType) GetPlaintext() string {
	return s.text
}
func (s SessionIDType) GetHash() string {
	sessionIDHash := sha256.Sum256([]byte(s.text))
	return base64.StdEncoding.EncodeToString(sessionIDHash[:])
}
func (s SessionIDType) GetAccountType() (AccountDescriptiveTitleType, error) {
	if len(s.GetPlaintext()) < 11 {
		return "", UnknownAccountDescriptiveTitleError
	}
	switch s.GetPlaintext()[8:11] {
	case "ptt": return AccountDescriptiveTitleTypePatient, nil
	case "prc": return AccountDescriptiveTitleTypePractitioner, nil
	default: return AccountDescriptiveTitleTypeInvalid, UnknownAccountDescriptiveTitleError
	}
}
func (s SessionIDType) GetMongoCollectionTitle() (string, error) {
	if len(s.GetPlaintext()) < 11 {
		return "", UnknownAccountDescriptiveTitleError
	}
	switch s.GetPlaintext()[8:11] {
	case "ptt": return MongoObjectPatientCollectionTitle, nil
	case "prc": return MongoObjectPractitionerCollectionTitle, nil
	default: return "", UnknownAccountDescriptiveTitleError
	}
}
func (s SessionIDType) GetAccountMongoEntry() (MongoObjectAccountBaseInterface, error) {
	var resp MongoObjectAccountBaseInterface

	collectionTitle, err := s.GetMongoCollectionTitle()
	if err != nil {
		return nil, err
	}

	accountType, err :=  s.GetAccountType()
	if err != nil {
		return nil, err
	}

	var t = s.GetHash()
	println(t)

	accountCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(collectionTitle)
	result := accountCollection.FindOne(
		appVariables.MongoContext,
		&bson.M{
			"session_id_hash": t,
		},
	)
	if result.Err() != nil {
		return nil, err
	}

	var dummyPatient MongoObjectPatient
	var dummyPractitioner MongoObjectPractitioner
	if accountType == AccountDescriptiveTitleTypePatient {
		result.Decode(&dummyPatient)
		resp = MongoObjectAccountBaseInterface(dummyPatient)
	} else if accountType == AccountDescriptiveTitleTypePractitioner {
		result.Decode(&dummyPractitioner)
		resp = MongoObjectAccountBaseInterface(dummyPractitioner)
	} else {
		return nil, UnknownAccountDescriptiveTitleError
	}

	return resp, nil
}

// CreateNewSessionID
// Generates a new plaintext SessionID which is currently unused by any
// pre-existing account.
func CreateNewSessionID(accountType AccountDescriptiveTitleType) (SessionIDType, error) {
	var stringIdentifier, stringIdErr = accountType.GetUniformStringIdentifier()
	if stringIdErr != nil {
		return SessionIDType{""}, stringIdErr
	}

	var resp SessionIDType
	var byteLength = SessionIDTypeLengthWithoutPrefix * 3 / 4
	b := make([]byte, byteLength)
	_, err := rand.Read(b)
	if err != nil {
		return SessionIDType{""}, err
	}

	var unvalidatedSessionID = SessionIDTypePrefixString + stringIdentifier + base64.URLEncoding.EncodeToString(b)
	var coll *mongo.Collection

	switch accountType {
	case AccountDescriptiveTitleTypePatient:
		coll = appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectPatientCollectionTitle)
		break
	case AccountDescriptiveTitleTypePractitioner:
		coll = appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectPractitionerCollectionTitle)
		break
	default:
		return SessionIDType{""}, UnknownAccountDescriptiveTitleError
	}

	var dummyInterface bson.M
	err = coll.FindOne(
		appVariables.MongoContext,
		&bson.D{{"session_id_hash", SessionIDType{unvalidatedSessionID}.GetHash()}},
	).Decode(&dummyInterface)

	if err == nil {
		return CreateNewSessionID(accountType)
	}
	if err != mongo.ErrNoDocuments {
		return SessionIDType{""}, err
	}

	resp.text = unvalidatedSessionID

	return resp, nil
}


/*
 *
 * JointDescriptiveTitle Struct and Modifiers
 *
 */

type JointDescriptiveTitleType string
const (
	JointDescriptiveTitleTypeKnee JointDescriptiveTitleType = "JOINT_TYPE_KNEE"
)