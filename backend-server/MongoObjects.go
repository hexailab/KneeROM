package main

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type MongoObjectAccountBaseInterface interface {
	GetMongoID() primitive.ObjectID
	GetEmail()   string
}

// MongoObjectPatient
// An element of the "patients" Mongo Collection.
// One noticeable aspect of this object is the SessionIDHash element,
// which implies - as is the case - that only one client may be logged in
// as this patient at a time. The patient must be automatically logged out of
// one client when logged into another, perhaps through the use of
// websockets -- this should be communicated effectively and immediately.
//
type MongoObjectPatient struct {
	MongoID          primitive.ObjectID        `bson:"_id"               json:"mongoId"`
	PractitionerID   primitive.ObjectID        `bson:"practitioner_id"   json:"practitionerId"`
	FullName         string                    `bson:"full_name"         json:"fullName"`
	Email            string                    `bson:"email_address"     json:"emailAddress"`
	PasswordHash     string                    `bson:"password_hash"     json:"passwordHash,omitempty"`
	SessionIDHash    string                    `bson:"session_id_hash"   json:"sessionIdHash,omitempty"`

	SessionIDPlain   string			   `bson:"-"                 json:"sessionIdPlain,omitempty"`
}

func (m MongoObjectPatient) GetMongoID() primitive.ObjectID {
	return m.MongoID
}
func (m MongoObjectPatient) GetEmail() string {
	return m.Email
}

const MongoObjectPatientCollectionTitle = "patients"

type CreateAndSaveNewPatientObjectInput struct {
	FullName string
	Email string
	Password string
	Practitioner primitive.ObjectID
}

func CreateAndSaveNewPatientObject(input CreateAndSaveNewPatientObjectInput) (MongoObjectPatient, error) {
	patientCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectPatientCollectionTitle)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return MongoObjectPatient{}, err
	}

	sessionIDPlain, err := CreateNewSessionID(AccountDescriptiveTitleTypePatient)
	if err != nil {
		return MongoObjectPatient{}, err
	}

	patient := MongoObjectPatient{
		MongoID: primitive.NewObjectID(),
		PractitionerID: input.Practitioner,
		FullName: input.FullName,
		Email: input.Email,
		PasswordHash: string(hashedPassword),
		SessionIDHash: sessionIDPlain.GetHash(),
	}

	var dummyInterface bson.M
	err = patientCollection.FindOne(
		appVariables.MongoContext,
		bson.D{{"email_address", input.Email}},
	).Decode(&dummyInterface)

	if err == nil {
		return MongoObjectPatient{}, EmailInUseError
	}
	if err != mongo.ErrNoDocuments {
		return MongoObjectPatient{}, err
	}

	if _, err = patientCollection.InsertOne(appVariables.MongoContext, patient); err != nil {
		return MongoObjectPatient{}, err
	}
	patient.SessionIDPlain = sessionIDPlain.GetPlaintext()

	return patient, err
}

var InvalidLoginInputMessage = errors.New("invalid_login_input")

type LoginInput struct {
	Email       string			`bson:"email_address"`
	Password    string			`bson:"password_hash"`
	AccountType AccountDescriptiveTitleType `bson:"account_type,omitempty"`
}
func LoginToAccount(input LoginInput) (SessionIDType, error) {
	var accountCollection *mongo.Collection
	var acct LoginInput
	if input.AccountType == AccountDescriptiveTitleTypePatient {
		accountCollection = appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectPatientCollectionTitle)
	} else if input.AccountType == AccountDescriptiveTitleTypePractitioner {
		accountCollection = appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectPractitionerCollectionTitle)
	}

	result := accountCollection.FindOne(
		appVariables.MongoContext,
		&bson.M{
			"email_address": input.Email,
		},
	)
	if result.Err() == mongo.ErrNoDocuments {
		return SessionIDType{""}, InvalidLoginInputMessage
	}
	result.Decode(&acct)

	sessionIDPlain, err := CreateNewSessionID(input.AccountType)
	if err != nil {
		return SessionIDType{""}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(acct.Password), []byte(input.Password))
	if err != nil {
		return SessionIDType{""}, InvalidLoginInputMessage
	}

	result = accountCollection.FindOneAndUpdate(
		appVariables.MongoContext,
		&bson.M{
			"email_address": input.Email,
		},
		&bson.D{{"$set",
			&bson.M{
				"session_id_hash": sessionIDPlain.GetHash(),
			}},
		},
	)
	if result.Err() != nil {
		return SessionIDType{""}, err
	}

	return sessionIDPlain, nil
}

// MongoObjectPractitioner
// An element of the "practitioners" Mongo Collection.
// One noticeable aspect of this object is the SessionIDHash element,
// which implies - as is the case - that only one client may be logged in
// as this patient at a time. The practitioner must be automatically logged out of
// one client when logged into another, perhaps through the use of
// websockets -- this should be communicated effectively and immediately.
//
type MongoObjectPractitioner struct {
	MongoID          primitive.ObjectID        `bson:"_id"                json:"mongoId"`
	FullName         string                    `bson:"full_name"          json:"fullName"`
	Institution	     string		               `bson:"institution"        json:"institution"`
	Email            string                    `bson:"email_address"      json:"emailAddress"`
	PasswordHash     string                    `bson:"password_hash"      json:"passwordHash,omitempty"`
	SessionIDHash    string                    `bson:"session_id_hash"    json:"sessionIdHash,omitempty"`

	SessionIDPlain   string                    `bson:"-"                  json:"sessionIdPlain,omitempty"`
}

func (m MongoObjectPractitioner) GetMongoID() primitive.ObjectID {
	return m.MongoID
}
func (m MongoObjectPractitioner) GetEmail() string {
	return m.Email
}

const MongoObjectPractitionerCollectionTitle = "practitioners"

var EmailInUseError = errors.New("email is already in use")

type CreateAndSaveNewPractitionerObjectInput struct {
	FullName    string
	Email       string
	Password    string
	Institution string
}
func CreateAndSaveNewPractitionerObject(input CreateAndSaveNewPractitionerObjectInput) (MongoObjectPractitioner, error) {
	practitionerCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectPractitionerCollectionTitle)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return MongoObjectPractitioner{}, err
	}

	sessionIDHash, err := CreateNewSessionID(AccountDescriptiveTitleTypePractitioner)
	if err != nil {
		return MongoObjectPractitioner{}, err
	}

	practitioner := MongoObjectPractitioner{
		MongoID: primitive.NewObjectID(),
		FullName: input.FullName,
		Email: input.Email,
		Institution: input.Institution,
		PasswordHash: string(hashedPassword),
		SessionIDHash: sessionIDHash.GetHash(),
	}

	var dummyInterface bson.M
	err = practitionerCollection.FindOne(
		appVariables.MongoContext,
		bson.D{{"email_address", input.Email}},
	).Decode(&dummyInterface)

	if err == nil {
		return MongoObjectPractitioner{}, EmailInUseError
	}
	if err != mongo.ErrNoDocuments {
		return MongoObjectPractitioner{}, err
	}

	if _, err = practitionerCollection.InsertOne(appVariables.MongoContext, practitioner); err != nil {
		return MongoObjectPractitioner{}, err
	}

	practitioner.SessionIDPlain = sessionIDHash.GetPlaintext()

	return practitioner, err
}

// MongoObjectRecordedEntry
// An element of the "recording_entries" Mongo Collection.
//
type MongoObjectRecordedEntry struct {
	MongoID          primitive.ObjectID        `bson:"_id"                json:"mongoId"`
	PatientID        primitive.ObjectID        `bson:"patient_id"         json:"patientId"`
	RangeOfMotion    float64		   		   `bson:"range_of_motion"    json:"rangeOfMotion"`
	DateTimeRecorded time.Time		           `bson:"datetime_recorded"  json:"datetimeRecorded"`
	JointID          JointDescriptiveTitleType `bson:"joint_id"           json:"jointId"`
}

const MongoObjectRecordedEntryCollectionTitle = "recording_entries"

type CreateAndSaveNewRecordedEntryObjectInput struct {
	Patient       primitive.ObjectID
	RangeOfMotion float64
	DateTime      time.Time
	JointID       JointDescriptiveTitleType
}

func CreateAndSaveNewRecordedEntryObject(input CreateAndSaveNewRecordedEntryObjectInput) (MongoObjectRecordedEntry, error) {
	entryCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectRecordedEntryCollectionTitle)

	recordedEntry := MongoObjectRecordedEntry{
		MongoID: primitive.NewObjectID(),
		PatientID: input.Patient,
		RangeOfMotion: input.RangeOfMotion,
		DateTimeRecorded: input.DateTime,
		JointID: input.JointID,
	}

	if _, err := entryCollection.InsertOne(appVariables.MongoContext, recordedEntry); err != nil {
		return MongoObjectRecordedEntry{}, err
	}

	return recordedEntry, nil
}

type IsPatientUnderPractitionerInput struct {
	Practitioner primitive.ObjectID
	Patient primitive.ObjectID
}

func IsPatientUnderPractitioner(input IsPatientUnderPractitionerInput) bool {
	var patient MongoObjectPatient

	patientCollection := appVariables.MongoClient.Database(MongoClientDatabaseTitle).Collection(MongoObjectPatientCollectionTitle)
	err := patientCollection.FindOne(appVariables.MongoContext, &bson.M{
		"_id": input.Patient,
		"practitioner_id": input.Practitioner,
	}).Decode(&patient)

	return err == nil
}