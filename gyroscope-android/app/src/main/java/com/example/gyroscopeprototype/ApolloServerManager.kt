package com.example.gyroscopeprototype

import com.apollographql.apollo3.ApolloClient
import com.example.PingQuery
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking

class ApolloServerManager {
    private lateinit var apolloClient: ApolloClient

    private fun testConnection() = runBlocking {
        launch {
            val response = apolloClient.query(PingQuery()).execute()

            if (response.data?.ping == "pong") {
                println("SUCCESS - AWS server is on, and we are ready to rumble.")
            } else {
                // TODO: Replace Exception with something a bit prettier, i.e. a popup window.
                throw Exception("ERROR - Either there is an issue with your internet access, or our server is down. Apologies.")
            }
        }
    }

    fun connectToMainServer() {
        apolloClient = ApolloClient.Builder().serverUrl(SERVER_URL).build()
        testConnection()
    }

    companion object {
        private const val SERVER_URL = "http://54.146.48.48:8080/graphql"
    }
}