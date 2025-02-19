package cloudformation

import kotlinx.serialization.json.Json
import java.io.File

data class Schema(
    val regions: List<String>,
    val services: List<Service>,
) {
    companion object {
        operator fun invoke(directory: String): Schema {
            val regions = readRegions()
            return Schema(
                regions = regions,
                services = readServices(regions, directory),
            )
        }

        private fun readRegions(): List<String> {
            val filename = Companion::class.java.getResource("/regions.txt")
                ?: error("Could not find the file")

            val regions = filename
                .readText()
                .lines()
                .filter { it.isNotBlank() }

            return regions
        }

        private data class ResourceDef(
            val service: String,
            val resource: String,
            val region: String,
            val isListable: Boolean,
            val required: List<String>,
        )

        private fun readServices(regions: List<String>, directory: String): List<Service> {
            val json = Json { ignoreUnknownKeys = true }
            return regions.flatMap { region ->
                val files = File("$directory/$region").listFiles()
                    ?: error("No files found for region $region")
                files.map { it.readText() }
                    .map { json.decodeFromString<AwsType>(it) }
                    .map {
                        val typeName = it.typeName
                        val parts = typeName.split("::")
                        if (parts.size != 3)
                            error(typeName)
                        if (parts[0] != "AWS")
                            return@map null
                        ResourceDef(
                            service = parts[1],
                            resource = parts[2],
                            region = region,
                            isListable = it.handlers?.list != null,
                            required = it.handlers?.list?.handlerSchema?.required ?: emptyList(),
                        )
                    }.filterNotNull()
            }.let(::groupByService)
        }

        private fun groupByService(defs: List<ResourceDef>): List<Service> =
            defs.groupBy { it.service }
                .map { (service, resources) ->
                    resources.groupBy { it.resource }
                    Service(
                        name = service,
                        resourceTypes = resources
                            .groupBy { it.resource }
                            .map { (name, defs) ->
                                ResourceType(
                                    name = name,
                                    regions = defs.map { it.region },
                                    isListable = defs.first().isListable,
                                    parents = defs.first().required,
                                )
                            },
                    )
                }
    }
}

data class Service(
    val name: String,
    val resourceTypes: List<ResourceType>,
) {
    companion object {
        val List<Service>.names
            get() = this.map { it.name }
    }
}

data class ResourceType(
    val name: String,
    val regions: List<String>,
    val isListable: Boolean,
    val parents: List<String>,
) {
    companion object {
        val List<ResourceType>.names
            get() = this.map { it.name }
    }
}
