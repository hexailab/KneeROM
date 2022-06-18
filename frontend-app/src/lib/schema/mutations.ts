import {gql} from "@urql/svelte";

export const GET_SESSION_ACCOUNT_DETAILS = gql`
    query GetAccountDetails($sessionId:String!) {
        getPractitionerAccountDetails(sessionId:$sessionId){
            fullName
        }
    }
`;

export const LOGIN_TO_PRACTITIONER_ACCOUNT = gql`
    mutation LoginToPractitioner($email:String!, $password:String!) {
        loginToAccount(
            email:$email
            password:$password
            accountType:ACCOUNT_TYPE_PRACTITIONER
        ) {
            sessionId
        }
    }
`;