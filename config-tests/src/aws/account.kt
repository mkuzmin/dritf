package aws

import aws.sdk.kotlin.services.account.AccountClient
import aws.sdk.kotlin.services.account.paginators.listRegionsPaginated
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.runBlocking

fun readRegions(): List<String> = runBlocking {
    AccountClient {
        region = "us-east-1" //TODO AWS_REGION env variable does not work
    }.use { client ->
        client
            .listRegionsPaginated()
            .toList()
            .flatMap { it.regions ?: error("Cannot get regions") }
            .map { it.regionName ?: error("Region name is null") }
    }
}
