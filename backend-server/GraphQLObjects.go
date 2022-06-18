package main

import (
	"github.com/graphql-go/graphql"
)

var GraphQLObjectPatient = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "Patient",
		Fields: graphql.Fields{
			"mongoId":          &graphql.Field{ Type: graphql.String },
			"practitionerId":   &graphql.Field{ Type: graphql.String },
			"fullName":         &graphql.Field{ Type: graphql.String },
			"emailAddress":     &graphql.Field{ Type: graphql.String },
			"passwordHash":     &graphql.Field{ Type: graphql.String },
			"sessionIdHash":    &graphql.Field{ Type: graphql.String },
			"sessionIdPlain":   &graphql.Field{ Type: graphql.String },
		},
	},
)

var GraphQLObjectPractitioner = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "Practitioner",
		Fields: graphql.Fields{
			"mongoId":         &graphql.Field{ Type: graphql.String },
			"institution":     &graphql.Field{ Type: graphql.String },
			"fullName":        &graphql.Field{ Type: graphql.String },
			"emailAddress":    &graphql.Field{ Type: graphql.String },
			"passwordHash":    &graphql.Field{ Type: graphql.String },
			"sessionIdHash":   &graphql.Field{ Type: graphql.String },
			"sessionIdPlain":  &graphql.Field{ Type: graphql.String },
		},
	},
)

var GraphQLObjectRecordedEntry = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "RecordedEntry",
		Fields: graphql.Fields{
			"mongoId":          &graphql.Field{ Type: graphql.String },
			"patientId":        &graphql.Field{ Type: graphql.String },
			"rangeOfMotion":    &graphql.Field{ Type: graphql.Float },
			"datetimeRecorded": &graphql.Field{ Type: graphql.DateTime },
			"jointId":          &graphql.Field{ Type: JointTypeEnum },
		},
	},
)

var AccountTypeEnum = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "AccountType",
		Values:      graphql.EnumValueConfigMap{
			"ACCOUNT_TYPE_PATIENT":      &graphql.EnumValueConfig{
				Value: "ACCOUNT_TYPE_PATIENT",
			},
			"ACCOUNT_TYPE_PRACTITIONER": &graphql.EnumValueConfig{
				Value: "ACCOUNT_TYPE_PRACTITIONER",
			},
		},
	},
)

var JointTypeEnum = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "JointType",
		Values:      graphql.EnumValueConfigMap{
			"JOINT_TYPE_KNEE": &graphql.EnumValueConfig{
				Value: "JOINT_TYPE_KNEE",
			},
		},
	},
)

type LoginResponse struct {
	SessionIDPlain string `bson:"session_id_plain" json:"sessionId"`
}

var GraphQLObjectLoginResponse = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "LoginResponse",
		Fields: graphql.Fields{
			"sessionId": &graphql.Field{ Type: graphql.String },
		},
	},
)

var GraphQLObjectPatientIDList = graphql.NewList(graphql.String)

/*
 *
 *  Mutation Objects
 *
 */

var RootMutationObject = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"loginToAccount": &graphql.Field{
				Type: GraphQLObjectLoginResponse,
				Args: graphql.FieldConfigArgument{
					"email":       &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
					"password":    &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
					"accountType": &graphql.ArgumentConfig{ Type: graphql.NewNonNull(AccountTypeEnum) },
				},
				Resolve: loginAccountResolver,
			},
			"createNewROMEntry": &graphql.Field{
				Type: GraphQLObjectRecordedEntry,
				Args: graphql.FieldConfigArgument{
					"sessionId":        &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
					"rangeOfMotion":    &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.Float) },
					"datetimeRecorded": &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.DateTime) },
					"jointId":          &graphql.ArgumentConfig{ Type: graphql.NewNonNull(JointTypeEnum) },
				},
				Resolve: createNewROMEntryResolver,
			},
			"createNewPractitionerAccount": &graphql.Field{
				Type: GraphQLObjectPractitioner,
				Args: graphql.FieldConfigArgument{
					"fullName":        &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
					"institution":     &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
					"emailAddress":    &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
					"password":        &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
				},
				Resolve: createNewPractitionerAccountResolver,
			},
			"createNewPatientAccount": &graphql.Field{
				Type: GraphQLObjectPractitioner,
				Args: graphql.FieldConfigArgument{
					"sessionId":       &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
					"fullName":        &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
					"emailAddress":    &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
					"password":        &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
				},
				Resolve: createNewPatientAccountResolver,
			},
		},
	},
)

func pingTestFunction(_ graphql.ResolveParams) (interface{}, error) {
	return "pong", nil
}

var	RootFieldObject = graphql.Fields{
	"getROMEntries": &graphql.Field{
		Type: graphql.NewList(GraphQLObjectRecordedEntry),
		Args: graphql.FieldConfigArgument{
			"sessionId":       &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
			"patientId":       &graphql.ArgumentConfig{ Type: graphql.String },
			"dateRangeBegin":  &graphql.ArgumentConfig{ Type: graphql.DateTime },
			"dateRangeEnd":    &graphql.ArgumentConfig{ Type: graphql.DateTime },
		},
		Resolve: getROMEntriesResolver,
	},
	"getPatientAccountDetails": &graphql.Field{
		Type: GraphQLObjectPatient,
		Args: graphql.FieldConfigArgument{
			"sessionId":       &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
		},
		Resolve: getPatientAccountDetailsResolver,
	},
	"getPractitionerAccountDetails": &graphql.Field{
		Type: GraphQLObjectPractitioner,
		Args: graphql.FieldConfigArgument{
			"sessionId":       &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
		},
		Resolve: getPractitionerAccountDetailsResolver,
	},
	"getPatientsUnderPractitioner": &graphql.Field{
		Type: GraphQLObjectPatientIDList,
		Args: graphql.FieldConfigArgument{
			"sessionId":       &graphql.ArgumentConfig{ Type: graphql.NewNonNull(graphql.String) },
		},
		Resolve: getPractitionerAccountsPatients,
	},
	"ping": &graphql.Field{
		Type: graphql.String,
		Resolve: pingTestFunction,
	},
}